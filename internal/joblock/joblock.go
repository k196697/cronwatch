// Package joblock provides a simple per-job mutex to prevent concurrent
// execution of the same cron job when a previous run is still in progress.
package joblock

import (
	"fmt"
	"sync"
)

// Locker manages per-job execution locks.
type Locker struct {
	mu    sync.Mutex
	locks map[string]*entry
}

type entry struct {
	mu     sync.Mutex
	locked bool
}

// New returns an initialised Locker.
func New() *Locker {
	return &Locker{
		locks: make(map[string]*entry),
	}
}

// TryAcquire attempts to acquire the lock for the given job name.
// It returns true and a release function on success, or false if the job is
// already running.
func (l *Locker) TryAcquire(job string) (bool, func(), error) {
	if job == "" {
		return false, nil, fmt.Errorf("joblock: job name must not be empty")
	}

	l.mu.Lock()
	e, ok := l.locks[job]
	if !ok {
		e = &entry{}
		l.locks[job] = e
	}
	l.mu.Unlock()

	e.mu.Lock()
	if e.locked {
		e.mu.Unlock()
		return false, nil, nil
	}
	e.locked = true
	e.mu.Unlock()

	release := func() {
		e.mu.Lock()
		e.locked = false
		e.mu.Unlock()
	}

	return true, release, nil
}

// IsLocked reports whether the named job currently holds a lock.
func (l *Locker) IsLocked(job string) bool {
	l.mu.Lock()
	e, ok := l.locks[job]
	l.mu.Unlock()
	if !ok {
		return false
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.locked
}
