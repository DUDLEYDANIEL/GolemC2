package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DUDLEYDANIEL/GolemC2/internal/agent/tasks"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var store *InMemoryStore

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
			tasks := store.GetTasks(agentID)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(tasks); err != nil {
				http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
			}
		} else if r.Method == http.MethodPost {
			var task tasks.Task
			if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			task.ID = uuid.New().String()
			task.CreatedAt = time.Now()
			store.AddTask(agentID, task)
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

		var res result
		if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		res.Completed = time.Now()
		store.Addresult(agentID, res)
		w.WriteHeader(http.StatusNoContent)
	})
}

func generateUUID() string {
	return uuid.New().String()
}
