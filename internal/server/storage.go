package server 

import (
  "sync"
  "time"
)

type InMemoryStore struct {
  mu      sync.Mutex
  tasks   map[string][]Task
  results map[string][]result
  agents  map[string]AgentInfo
}

func NewStore() *InMemoryStore {
  return &InMemoryStore{
    tasks:   make(map[string][]Task),
    results: make(map[string][]result),
    agents:  make(map[string]AgentInfo),
  }
}

func (s *InMemoryStore) AddTask(agentID string, t Task ){
  s.mu.Lock()
  defer s.mu.Unlock()
  s.tasks[agentID] = append(s.tasks[agentID], t)
}

func (s *InMemoryStore) GetTasks (agentID string) []Task {
  s.mu.Lock()
  defer s.mu.Unlock()
  ts := s.tasks[agentID]
  s.tasks[agentID] = []Task{}
  return ts 
} 

func (s *InMemoryStore) Addresult(agentID string, r result){
  s.mu.Lock()
  defer s.mu.Unlock()
  s.results[agentID] = append(s.results[agentID], r)
}

func (s *InMemoryStore) Getresults (agentID string) []result{
  s.mu.Lock()
  defer s.mu.Unlock()
  if results, exists := s.results[agentID]; exists{
    return results
  }
  return []result{}
}

func (s *InMemoryStore) UpdateAgentLastSeen(agentID string, lastSeen time.Time){
  s.mu.Lock()
  defer s.mu.Unlock()
  info , exists := s.agents[agentID]
  if !exists {
    info = AgentInfo{ID: agentID}
  }
  info.LastSeen = lastSeen
  s.agents[agentID] = info 
}

func (s *InMemoryStore) GetAgentInfo (agentID string)(AgentInfo, bool){
  s.mu.Lock()
  defer s.mu.Unlock()
  info , exists := s.agents[agentID]
  return info, exists
}

func (s *InMemoryStore) GetAllAgents() []AgentInfo{
  s.mu.Lock()
  defer s.mu.Unlock()
  agents := make([]AgentInfo, 0, len(s.agents))
  for _, info := range s.agents {
    agents = append(agents, info)
  }
  return agents
}


