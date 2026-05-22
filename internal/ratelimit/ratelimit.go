// Package ratelimit provides a token-bucket style rate limiter
// for controlling how frequently alerts are emitted per job.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks per-key rate limits using a sliding window approach.
type Limiter struct {
	mu       sync.Mutex
	window   time.Duration
	maxCalls int
	buckets  map[string][]time.Time
}

// New creates a Limiter that allows at most maxCalls events per key
// within the given window duration.
func New(window time.Duration, maxCalls int) *Limiter {
	return &Limiter{
		window:   window,
		maxCalls: maxCalls,
		buckets:  make(map[string][]time.Time),
	}
}

// Allow reports whether an event for the given key is permitted under
// the current rate limit. It records the event if allowed.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)

	times := l.buckets[key]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= l.maxCalls {
		l.buckets[key] = filtered
		return false
	}

	l.buckets[key] = append(filtered, now)
	return true
}

// Reset clears the event history for a specific key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}

// Remaining returns how many more calls are allowed for the key
// within the current window.
func (l *Limiter) Remaining(key string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)

	count := 0
	for _, t := range l.buckets[key] {
		if t.After(cutoff) {
			count++
		}
	}

	remaining := l.maxCalls - count
	if remaining < 0 {
		return 0
	}
	return remaining
}
