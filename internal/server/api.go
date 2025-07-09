package server 

import (
  "encoding/json"
  "net/http"
  "time"

  "github.com/google/uuid"
)

var store *InMemoryStore

func RegisterHandler(mux *http.ServeMux, s *InMemoryStore){
  store = s

  mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request){
    if r.Method != http.MethodPost {
      http.Error(w, "Method Not Allowed" , http.StatusMethodNotAllowed)
      return 
    }

    agentID := r.URL.QUERY().GET("agent_id")
    if agentID == "" {
      agentID = uuid.New().String()
    }
    
    store.UpdateAgentLastSeen(agentID, time.Now())

    w.Header().Set("Content-Type", "application/json")

    if err:= json.NewEncoder(w).Encode(map[string]string{"agent_id": agentID}); err!= nil{
      http.Error(w, "Failed to encode response", http.StatusInternalError)
    }
  })

  mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request){
    if r.Method  != http.MethodGet {
      http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
      return 
    }

    agentID := r.URL.QUERY().Get("agent_id")
    if agentID == ""{
      http.Error(w, "Missing agent_id",  http.StatusBadRequest)
      return 
    }
    tasks := store.GetTasks(agentID)
    w.Header().Set("Content-Type":"application/json")
    if err := json.NewEncoder(w).Encode(tasks); err != nil {
      http.Error(w, "Failed to encode the respnse", http.StatusInternalError)
      return 
    }
  })

  mux.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request){
    if r.Method != http.MethodPost{
      http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
      return 
    }

    agentID := r.URL.QUERY().Get("agent_id")
    if agentID == ""{
      http.Error(w, "Missing agent_id", http.StatusBadRequest)
      return 
    }

    var res Result 
    if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
      http.Error(w, "Invalid request body", http.StatusBadRequest)
      return 
    }
    res.Completed = time.Now()
    store.AddResult(agentID, res)
    w.WriteHeader(http.StatusNoContent)
  })
}

func generateUUID() string{
  return uuid.New().String()
}

