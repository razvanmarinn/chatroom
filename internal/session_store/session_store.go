package session_store

import (
	"sync"
)

var SessionStore = &InMemorySessionStore{
	sessions: make(map[string]string),
	mu:       &sync.RWMutex{},
}

type InMemorySessionStore struct {
	sessions map[string]string
	mu       *sync.RWMutex
}

func (s *InMemorySessionStore) Set(token, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[token] = userID
}

func (s *InMemorySessionStore) Get(token string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userID, found := s.sessions[token]
	return userID, found
}

func (s *InMemorySessionStore) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}
