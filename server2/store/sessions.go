package store

import (
	"server2/models"
	"sync"

	"github.com/google/uuid"
)

// manages all active sessions
type SessionStore struct {
	sessions  map[string]*models.Session
	mu        sync.RWMutex
	dimension int
}

//  creates a new session store
func NewSessionStore(dimension int) *SessionStore {
	return &SessionStore{
		sessions:  make(map[string]*models.Session),
		dimension: dimension,
	}
}

// creates a new session and returns its ID
func (s *SessionStore) Create() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	session := models.NewSession(id, s.dimension)
	s.sessions[id] = session
	return id
}

// returns a session by ID
func (s *SessionStore) Get(id string) *models.Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[id]
}

