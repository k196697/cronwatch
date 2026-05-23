// Package debounce provides a mechanism to suppress repeated alerts
// for the same job within a configurable quiet period.
package debounce

import (
	"sync"
	"time"
)

// Debouncer tracks the last alert time per job and suppresses
// subsequent alerts that arrive within the quiet window.
type Debouncer struct {
	mu      sync.Mutex
	last    map[string]time.Time
	window  time.Duration
	nowFunc func() time.Time
}

// New creates a Debouncer with the given quiet window duration.
// Alerts for the same job key are suppressed until the window elapses.
func New(window time.Duration) *Debouncer {
	return &Debouncer{
		last:    make(map[string]time.Time),
		window:  window,
		nowFunc: time.Now,
	}
}

// Allow returns true if an alert for the given job key should be
// forwarded, and false if it falls within the quiet window of the
// previous alert. When true is returned the timestamp is updated.
func (d *Debouncer) Allow(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.nowFunc()
	if t, ok := d.last[key]; ok && now.Sub(t) < d.window {
		return false
	}
	d.last[key] = now
	return true
}

// Reset clears the recorded timestamp for the given key, allowing
// the next alert to pass through immediately.
func (d *Debouncer) Reset(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.last, key)
}

// ResetAll clears all recorded timestamps.
func (d *Debouncer) ResetAll() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.last = make(map[string]time.Time)
}
