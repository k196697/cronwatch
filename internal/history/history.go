// Package history provides a simple in-memory store for tracking
// recent cron job execution results.
package history

import (
	"sync"
	"time"
)

// Result holds the outcome of a single job execution.
type Result struct {
	JobName   string
	StartedAt time.Time
	Duration  time.Duration
	Success   bool
	Output    string
}

// Store keeps the last N results per job name.
type Store struct {
	mu      sync.RWMutex
	records map[string][]Result
	limit   int
}

// New creates a Store that retains at most limit results per job.
func New(limit int) *Store {
	if limit <= 0 {
		limit = 10
	}
	return &Store{
		records: make(map[string][]Result),
		limit:   limit,
	}
}

// Record appends a result for the given job, evicting the oldest entry
// when the per-job limit is exceeded.
func (s *Store) Record(r Result) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries := s.records[r.JobName]
	entries = append(entries, r)
	if len(entries) > s.limit {
		entries = entries[len(entries)-s.limit:]
	}
	s.records[r.JobName] = entries
}

// Latest returns the most recent result for a job, and false if none exists.
func (s *Store) Latest(jobName string) (Result, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := s.records[jobName]
	if len(entries) == 0 {
		return Result{}, false
	}
	return entries[len(entries)-1], true
}

// All returns a copy of all stored results for a job.
func (s *Store) All(jobName string) []Result {
	s.mu.RLock()
	defer s.mu.RUnlock()

	src := s.records[jobName]
	out := make([]Result, len(src))
	copy(out, src)
	return out
}

// ConsecutiveFailures returns the number of consecutive failures
// at the tail of the history for a job.
func (s *Store) ConsecutiveFailures(jobName string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := s.records[jobName]
	count := 0
	for i := len(entries) - 1; i >= 0; i-- {
		if !entries[i].Success {
			count++
		} else {
			break
		}
	}
	return count
}
