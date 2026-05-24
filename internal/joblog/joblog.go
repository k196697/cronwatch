package joblog

import (
	"errors"
	"sync"
	"time"
)

// Entry represents a single log line emitted during a job run.
type Entry struct {
	JobName   string
	Timestamp time.Time
	Level     string // "info", "warn", "error"
	Message   string
}

// Store holds recent log entries per job, capped at a configurable limit.
type Store struct {
	mu    sync.RWMutex
	logs  map[string][]Entry
	limit int
}

// New returns a Store that retains at most limit entries per job.
func New(limit int) (*Store, error) {
	if limit <= 0 {
		return nil, errors.New("joblog: limit must be greater than zero")
	}
	return &Store{
		logs:  make(map[string][]Entry),
		limit: limit,
	}, nil
}

// Append adds a log entry for the named job.
func (s *Store) Append(jobName, level, message string) error {
	if jobName == "" {
		return errors.New("joblog: job name must not be empty")
	}
	if level == "" {
		return errors.New("joblog: level must not be empty")
	}
	e := Entry{
		JobName:   jobName,
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logs[jobName] = append(s.logs[jobName], e)
	if len(s.logs[jobName]) > s.limit {
		s.logs[jobName] = s.logs[jobName][len(s.logs[jobName])-s.limit:]
	}
	return nil
}

// Get returns a copy of all retained log entries for the named job.
func (s *Store) Get(jobName string) []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.logs[jobName]
	out := make([]Entry, len(src))
	copy(out, src)
	return out
}

// Clear removes all log entries for the named job.
func (s *Store) Clear(jobName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.logs, jobName)
}
