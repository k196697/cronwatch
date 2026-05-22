// Package alertmanager coordinates alert dispatching with deduplication
// and rate-limiting to prevent notification storms.
package alertmanager

import (
	"sync"
	"time"

	"github.com/example/cronwatch/internal/notifier"
)

// Alert represents a trigger event for a cron job.
type Alert struct {
	JobName   string
	Command   string
	Error     string
	Output    string
	Duration  time.Duration
	Timestamp time.Time
	Timeout   bool
}

// Manager deduplicates and rate-limits alerts before forwarding them.
type Manager struct {
	notifier  notifier.Notifier
	cooldown  time.Duration
	mu        sync.Mutex
	lastSent  map[string]time.Time
}

// New creates a Manager with the given notifier and per-job cooldown duration.
func New(n notifier.Notifier, cooldown time.Duration) *Manager {
	return &Manager{
		notifier: n,
		cooldown: cooldown,
		lastSent: make(map[string]time.Time),
	}
}

// Send dispatches an alert for the given job unless a cooldown is active.
// Returns true if the alert was forwarded, false if suppressed.
func (m *Manager) Send(a Alert) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if last, ok := m.lastSent[a.JobName]; ok {
		if time.Since(last) < m.cooldown {
			return false, nil
		}
	}

	na := notifier.Alert{
		JobName:  a.JobName,
		Command:  a.Command,
		Error:    a.Error,
		Output:   a.Output,
		Duration: a.Duration,
		Timeout:  a.Timeout,
	}

	if err := m.notifier.Send(na); err != nil {
		return false, err
	}

	m.lastSent[a.JobName] = time.Now()
	return true, nil
}

// Reset clears the cooldown state for a specific job.
func (m *Manager) Reset(jobName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.lastSent, jobName)
}
