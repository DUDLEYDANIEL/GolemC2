package server 

import (
	"time"
)

type Task struct {
	ID string  	`json:"id"`
	Command string `json:"command"`
	Args []string 	`json:"args"`	
}

type result struct {
	TaskID string `json:"task_id"`		
	Output []byte 	`json:"output"`
	ExitCode int 	`json:"exit_code"`
	Completed time.Time `json:"completed"`
}

type AgentInfo struct {
  ID string `json:"id"`
  LastSeen time.Time `json:"last_seen"`
}

