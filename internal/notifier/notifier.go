// Package notifier provides alerting functionality for cronwatch.
package notifier

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

// Alert represents a notification to be sent when a job fails or times out.
type Alert struct {
	JobName   string
	Command   string
	Error     error
	Output    string
	Duration  time.Duration
	Timestamp time.Time
}

// Notifier sends alerts via a configured backend.
type Notifier struct {
	cfg Config
}

// Config holds notifier configuration.
type Config struct {
	SMTPHost string
	SMTPPort int
	From     string
	To       []string
	Username string
	Password string
}

// New creates a new Notifier with the given config.
func New(cfg Config) *Notifier {
	return &Notifier{cfg: cfg}
}

// Send dispatches an alert using the configured backend.
func (n *Notifier) Send(alert Alert) error {
	if n.cfg.SMTPHost == "" {
		return fmt.Errorf("notifier: smtp host is not configured")
	}
	subject := fmt.Sprintf("[cronwatch] Job '%s' failed", alert.JobName)
	body := formatBody(alert)
	return n.sendEmail(subject, body)
}

func (n *Notifier) sendEmail(subject, body string) error {
	addr := fmt.Sprintf("%s:%d", n.cfg.SMTPHost, n.cfg.Port())
	auth := smtp.PlainAuth("", n.cfg.Username, n.cfg.Password, n.cfg.SMTPHost)
	msg := buildMessage(n.cfg.From, n.cfg.To, subject, body)
	return smtp.SendMail(addr, auth, n.cfg.From, n.cfg.To, msg)
}

func (n *Config) Port() int {
	if n.SMTPPort == 0 {
		return 587
	}
	return n.SMTPPort
}

func buildMessage(from string, to []string, subject, body string) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "From: %s\r\n", from)
	fmt.Fprintf(&buf, "To: %s\r\n", strings.Join(to, ", "))
	fmt.Fprintf(&buf, "Subject: %s\r\n", subject)
	fmt.Fprintf(&buf, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&buf, "Content-Type: text/plain; charset=utf-8\r\n")
	fmt.Fprintf(&buf, "\r\n")
	fmt.Fprintf(&buf, "%s", body)
	return buf.Bytes()
}

func formatBody(a Alert) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Job:       %s\n", a.JobName))
	sb.WriteString(fmt.Sprintf("Command:   %s\n", a.Command))
	sb.WriteString(fmt.Sprintf("Time:      %s\n", a.Timestamp.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Duration:  %s\n", a.Duration))
	sb.WriteString(fmt.Sprintf("Error:     %v\n", a.Error))
	if a.Output != "" {
		sb.WriteString(fmt.Sprintf("\nOutput:\n%s\n", a.Output))
	}
	return sb.String()
}
