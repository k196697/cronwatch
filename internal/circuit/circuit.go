// Package circuit implements a simple circuit breaker for job execution.
// It tracks consecutive failures and opens the circuit after a threshold,
// preventing further executions until a reset timeout has elapsed.
package circuit

import (
	"fmt"
	"sync"
	"time"
)

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // failing; requests blocked
	StateHalfOpen              // testing if service recovered
)

// Breaker is a per-job circuit breaker.
type Breaker struct {
	mu          sync.Mutex
	threshold   int
	resetAfter  time.Duration
	failures    int
	state       State
	openedAt    time.Time
}

// New creates a Breaker that opens after threshold consecutive failures
// and attempts a reset after resetAfter duration.
func New(threshold int, resetAfter time.Duration) *Breaker {
	if threshold <= 0 {
		threshold = 3
	}
	if resetAfter <= 0 {
		resetAfter = 30 * time.Second
	}
	return &Breaker{
		threshold:  threshold,
		resetAfter: resetAfter,
		state:      StateClosed,
	}
}

// Allow reports whether the circuit allows execution.
// It transitions an open circuit to half-open once the reset window passes.
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return true
	case StateHalfOpen:
		return true
	case StateOpen:
		if time.Since(b.openedAt) >= b.resetAfter {
			b.state = StateHalfOpen
			return true
		}
		return false
	}
	return false
}

// RecordSuccess resets the breaker to closed state.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure counter and opens the circuit
// if the threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.threshold {
		b.state = StateOpen
		b.openedAt = time.Now()
	}
}

// State returns the current circuit state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// String returns a human-readable state label.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return fmt.Sprintf("unknown(%d)", int(s))
	}
}
