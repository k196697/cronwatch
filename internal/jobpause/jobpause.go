package jobpause

import (
	"errors"
	"sync"
	"time"
)

// Store tracks paused jobs and their optional resume times.
type Store struct {
	mu      sync.RWMutex
	paused  map[string]entry
}

type entry struct {
	pausedAt time.Time
	resumeAt *time.Time // nil means paused indefinitely
}

// New returns an initialised Store.
func New() *Store {
	return &Store{
		paused: make(map[string]entry),
	}
}

// Pause marks a job as paused. If duration is zero the pause is indefinite.
func (s *Store) Pause(job string, duration time.Duration) error {
	if job == "" {
		return errors.New("jobpause: job name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	e := entry{pausedAt: time.Now()}
	if duration > 0 {
		t := e.pausedAt.Add(duration)
		e.resumeAt = &t
	}
	s.paused[job] = e
	return nil
}

// Resume removes a manual pause for the given job.
func (s *Store) Resume(job string) error {
	if job == "" {
		return errors.New("jobpause: job name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.paused, job)
	return nil
}

// IsPaused reports whether the job is currently paused.
// Timed pauses that have expired are automatically cleared.
func (s *Store) IsPaused(job string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.paused[job]
	if !ok {
		return false
	}
	if e.resumeAt != nil && time.Now().After(*e.resumeAt) {
		delete(s.paused, job)
		return false
	}
	return true
}

// All returns a snapshot of all currently paused jobs and their entries.
func (s *Store) All() map[string]time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make(map[string]time.Time, len(s.paused))
	for k, e := range s.paused {
		out[k] = e.pausedAt
	}
	return out
}
