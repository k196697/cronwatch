// Package jobstore provides an in-memory registry of known cron jobs
// and their current execution state.
package jobstore

import (
	"sync"
	"time"
)

// State represents the current execution state of a cron job.
type State int

const (
	StateUnknown State = iota
	StateRunning
	StateSuccess
	StateFailed
	StateTimedOut
)

// Entry holds the current state and metadata for a single job.
type Entry struct {
	JobName   string
	State     State
	LastStart time.Time
	LastEnd   time.Time
	PID       int
}

// Store is a thread-safe in-memory registry of job states.
type Store struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New returns an initialised Store.
func New() *Store {
	return &Store{
		entries: make(map[string]*Entry),
	}
}

// Set creates or updates the entry for the given job name.
func (s *Store) Set(e Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := e
	s.entries[e.JobName] = &copy
}

// Get returns the entry for the given job name and whether it was found.
func (s *Store) Get(jobName string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[jobName]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot of all entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, *e)
	}
	return out
}

// Delete removes the entry for the given job name.
func (s *Store) Delete(jobName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, jobName)
}
