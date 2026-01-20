package models

import (
	"sync"
)

//  represents a user's food selection session
type Session struct {
	ID           string
	IntentVector []float64
	SeenFoods    map[string]bool
	Completed    bool
	FinalChoice  string
	mu           sync.RWMutex
}

//  creates a new session with a neutral intent vector
func NewSession(id string, vectorSize int) *Session {
	intent := make([]float64, vectorSize)
	return &Session{
		ID:           id,
		IntentVector: intent,
		SeenFoods:    make(map[string]bool),
		Completed:    false,
	}
}

//  marks a food as seen in this session
func (s *Session) MarkSeen(foodID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.SeenFoods[foodID] = true
}

//  checks if a food has been seen
func (s *Session) HasSeen(foodID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.SeenFoods[foodID]
}

//  updates the session's intent vector
func (s *Session) UpdateIntent(newIntent []float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IntentVector = newIntent
}

// returns a copy of the intent vector
func (s *Session) GetIntent() []float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]float64, len(s.IntentVector))
	copy(result, s.IntentVector)
	return result
}

//  marks the session as completed with the final choice
func (s *Session) Complete(foodName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Completed = true
	s.FinalChoice = foodName
}

//  checks if the session is completed
func (s *Session) IsCompleted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Completed
}

