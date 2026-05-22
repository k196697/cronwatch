package notifier

// Sender is the interface implemented by any alert backend.
// Implementations must be safe for concurrent use.
type Sender interface {
	Send(alert Alert) error
}

// Ensure *Notifier satisfies Sender at compile time.
var _ Sender = (*Notifier)(nil)

// NoopSender discards all alerts. Useful when notifications are disabled.
type NoopSender struct{}

// Send does nothing and always returns nil.
func (NoopSender) Send(_ Alert) error { return nil }
