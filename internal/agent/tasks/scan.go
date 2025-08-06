package tasks

import (
	"fmt"
	"net"
	"time"
)

type ScanHandler struct{}

func (h *ScanHandler) Run(params map[string]string) (string, int){
	host, ok := params["host"]
	if !ok {
		return "Missing host parameter", 1
	}
	port , ok := params["port"]
	if !ok {
		return "Missing port parameter", 1
	}
	address := fmt.Sprintf("%s :%s", host, port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return fmt.Sprintf("Connection to %s failed: %v", address, err), 1
	}
	defer conn.Close()
	return fmt.Sprintf("Port %s is open", port), 0
}