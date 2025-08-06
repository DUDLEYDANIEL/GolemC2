package tasks

import "time"

type Task struct {
    ID        string            `json:"id"`
    Type      string            `json:"type"`
    Params    map[string]string `json:"params"`
    CreatedAt time.Time         `json:"created_at"`
}

type Result struct {
    TaskID    string    `json:"task_id"`
    Output    string    `json:"output"`
    ExitCode  int       `json:"exit_code"`
    Completed time.Time `json:"completed"`
}

type TaskHandler interface {
    Run(params map[string]string) (string, int)
}