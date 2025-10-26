package server

import (
	"sync"
	"time"
)

type InMemoryStore struct {
	mu      sync.Mutex
	tasks   map[string][]Task
	results map[string][]Result
	agents  map[string]AgentInfo
}

func NewStore() *InMemoryStore {
	return &InMemoryStore{
		tasks:   make(map[string][]Task),
		results: make(map[string][]Result),
		agents:  make(map[string]AgentInfo),
	}
}

func (s *InMemoryStore) AddTask(agentID string, t Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[agentID] = append(s.tasks[agentID], t)
	return nil
}

func (s *InMemoryStore) GetTasks(agentID string) []Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	ts := s.tasks[agentID]
	s.tasks[agentID] = []Task{}
	return ts
}

func (s *InMemoryStore) AddResult(agentID string, r Result) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.results[agentID] == nil {
		s.results[agentID] = []Result{}
	}
	s.results[agentID] = append(s.results[agentID], r)
	return nil
}

func (s *InMemoryStore) GetResults(agentID string) []Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	if results, exists := s.results[agentID]; exists {
		return results
	}
	return []Result{}
}

func (s *InMemoryStore) UpdateAgentLastSeen(agentID string, lastSeen time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	info, exists := s.agents[agentID]
	if !exists {
		info = AgentInfo{ID: agentID}
	}
	info.LastSeen = lastSeen
	s.agents[agentID] = info
}

func (s *InMemoryStore) GetAgentInfo(agentID string) (AgentInfo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	info, exists := s.agents[agentID]
	return info, exists
}

func (s *InMemoryStore) GetAllAgents() []AgentInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	agents := make([]AgentInfo, 0, len(s.agents))
	for _, info := range s.agents {
		agents = append(agents, info)
	}
	return agents
}
