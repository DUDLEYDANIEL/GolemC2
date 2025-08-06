package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DUDLEYDANIEL/GolemC2/internal/common"
	"github.com/DUDLEYDANIEL/GolemC2/pkg/crypto"
	"github.com/DUDLEYDANIEL/GolemC2/pkg/logging"
	"github.com/spf13/cobra"
)

// Agent represents an active agent in the C2 system
type Agent struct {
	ID       string `json:"id"`
	LastSeen string `json:"last_seen"`
}

// Task represents a task to be assigned to an agent
type Task struct {
	Type   string            `json:"type"`
	Params map[string]string `json:"params"`
}

// Result represents the outcome of a task
type Result struct {
	TaskID   string `json:"task_id"`
	Output   string `json:"output"`
	ExitCode int    `json:"exit_code"`
}

func main() {
	// Parse and validate configuration
	cfg, err := common.ParseFlags()
	if err != nil {
		logging.Log.Fatal("Failed to parse flags: ", err)
	}
	if err := cfg.Validate("cli"); err != nil {
		logging.Log.Fatal("Configuration validation failed: ", err)
	}

	// Initialize logger
	if err := logging.Init("info", "", false); err != nil {
		logging.Log.Fatal("Failed to initialize logger: ", err)
	}

	// Load TLS configuration for secure communication
	tlsConfig, err := crypto.LoadTLSConfig(cfg.TLSCertPath, cfg.TLSKeyPath, cfg.CACertPath)
	if err != nil {
		logging.Log.Fatal("Failed to load TLS config: ", err)
	}

	// Create HTTP client with TLS
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Define the root command
	rootCmd := &cobra.Command{
		Use:   "c2tool",
		Short: "Command-line interface for the C2 system",
	}

	// Define subcommands
	agentsCmd := &cobra.Command{
		Use:   "agents",
		Short: "List active agents",
		Run: func(cmd *cobra.Command, args []string) {
			listAgents(client, cfg.ServerUrl)
		},
	}

	taskCmd := &cobra.Command{
		Use:   "task [agent_id] [task_type] [params...]",
		Short: "Assign a task to an agent",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			agentID := args[0]
			taskType := args[1]
			params := make(map[string]string)
			for i := 2; i < len(args); i += 2 {
				if i+1 < len(args) {
					params[args[i]] = args[i+1]
				}
			}
			assignTask(client, cfg.ServerUrl, agentID, taskType, params)
		},
	}

	resultsCmd := &cobra.Command{
		Use:   "results [agent_id]",
		Short: "View results for an agent",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			agentID := args[0]
			viewResults(client, cfg.ServerUrl, agentID)
		},
	}

	// Add subcommands to root command
	rootCmd.AddCommand(agentsCmd, taskCmd, resultsCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		logging.Log.Fatal(err)
	}
}

// listAgents fetches and displays the list of active agents from the server
func listAgents(client *http.Client, serverURL string) {
	resp, err := client.Get(serverURL + "/agents")
	if err != nil {
		logging.Log.Errorf("Failed to list agents: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logging.Log.Errorf("Server returned non-200 status: %d", resp.StatusCode)
		return
	}

	var agents []Agent
	if err := json.NewDecoder(resp.Body).Decode(&agents); err != nil {
		logging.Log.Errorf("Failed to decode agents: %v", err)
		return
	}

	for _, agent := range agents {
		fmt.Printf("ID: %s, Last Seen: %s\n", agent.ID, agent.LastSeen)
	}
}

// assignTask sends a task assignment to the server for a specific agent
func assignTask(client *http.Client, serverURL, agentID, taskType string, params map[string]string) {
	task := Task{
		Type:   taskType,
		Params: params,
	}
	body, err := json.Marshal(task)
	if err != nil {
		logging.Log.Errorf("Failed to marshal task: %v", err)
		return
	}

	resp, err := client.Post(serverURL+"/tasks?agent_id="+agentID, "application/json", bytes.NewReader(body))
	if err != nil {
		logging.Log.Errorf("Failed to assign task: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		logging.Log.Errorf("Server returned non-201 status: %d", resp.StatusCode)
		return
	}

	var taskResp map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		logging.Log.Errorf("Failed to decode task response: %v", err)
		return
	}

	fmt.Printf("Task assigned with ID: %s\n", taskResp["task_id"])
}

// viewResults fetches and displays the results for a specific agent
func viewResults(client *http.Client, serverURL, agentID string) {
	resp, err := client.Get(serverURL + "/results?agent_id=" + agentID)
	if err != nil {
		logging.Log.Errorf("Failed to view results: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logging.Log.Errorf("Server returned non-200 status: %d", resp.StatusCode)
		return
	}

	var results []Result
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		logging.Log.Errorf("Failed to decode results: %v", err)
		return
	}

	for _, result := range results {
		fmt.Printf("Task ID: %s, Output: %s, Exit Code: %d\n", result.TaskID, result.Output, result.ExitCode)
	}
}