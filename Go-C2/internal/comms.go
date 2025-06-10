package internal

type Task struct{
  ID string `json:"id"`
  Command string `json:"command"`
}

type Result struct {
  TaskID string `json:"task_id"`
  Output string `json:"output"`
}


