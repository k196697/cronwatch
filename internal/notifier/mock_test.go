package notifier

import "sync"

// MockNotifier records alerts sent during tests without making real network calls.
type MockNotifier struct {
	mu     sync.Mutex
	alerts []Alert
	err    error
}

// SetError configures the mock to return err on Send.
func (m *MockNotifier) SetError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.err = err
}

// Send records the alert and returns any configured error.
func (m *MockNotifier) Send(a Alert) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.err != nil {
		return m.err
	}
	m.alerts = append(m.alerts, a)
	return nil
}

// Alerts returns a copy of all recorded alerts.
func (m *MockNotifier) Alerts() []Alert {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Alert, len(m.alerts))
	copy(out, m.alerts)
	return out
}

// Count returns the number of alerts received.
func (m *MockNotifier) Count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.alerts)
}
