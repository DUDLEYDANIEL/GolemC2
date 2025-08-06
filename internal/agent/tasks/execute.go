package tasks

import (
	"bytes"
	"os/exec"
	"strings"
)

type ExecuteHandler struct{}

func (h *ExecuteHandler) Run(params map[string]string) (string, int){
	cmdStr, ok := params["command"]
	if !ok{
		return "Missing command parameter", 1
	}
	cmdArgs := strings.Fields(cmdStr)
	if len(cmdArgs) == 0{
		return "Empty Command", 1
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return out.String(), 1
	}
	return out.String(), 0
	}