package jobpriority

import (
	"errors"
	"sync"
)

// Level represents a job priority tier.
type Level int

const (
	Low    Level = 1
	Normal Level = 5
	High   Level = 10
)

// ErrInvalidPriority is returned when an out-of-range priority is assigned.
var ErrInvalidPriority = errors.New("jobpriority: level must be between 1 and 10")

// ErrEmptyJobName is returned when an empty job name is provided.
var ErrEmptyJobName = errors.New("jobpriority: job name must not be empty")

// Registry maps job names to their priority levels.
type Registry struct {
	mu       sync.RWMutex
	priority map[string]Level
}

// New returns an initialised Registry.
func New() *Registry {
	return &Registry{
		priority: make(map[string]Level),
	}
}

// Set assigns a priority level to the named job.
func (r *Registry) Set(job string, level Level) error {
	if job == "" {
		return ErrEmptyJobName
	}
	if level < 1 || level > 10 {
		return ErrInvalidPriority
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.priority[job] = level
	return nil
}

// Get returns the priority level for the named job.
// If no level has been set, Normal is returned as the default.
func (r *Registry) Get(job string) Level {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if lvl, ok := r.priority[job]; ok {
		return lvl
	}
	return Normal
}

// Delete removes the priority entry for the named job.
func (r *Registry) Delete(job string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.priority, job)
}

// All returns a snapshot of all registered priorities.
func (r *Registry) All() map[string]Level {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]Level, len(r.priority))
	for k, v := range r.priority {
		out[k] = v
	}
	return out
}
