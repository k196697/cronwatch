package notifier

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func makeAlert() Alert {
	return Alert{
		JobName:   "backup",
		Command:   "/usr/bin/backup.sh",
		Error:     errors.New("exit status 1"),
		Output:    "disk full",
		Duration:  2 * time.Second,
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
}

func TestSend_NoSMTPHost(t *testing.T) {
	n := New(Config{})
	err := n.Send(makeAlert())
	if err == nil {
		t.Fatal("expected error when smtp host is empty")
	}
	if !strings.Contains(err.Error(), "smtp host") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestFormatBody_ContainsFields(t *testing.T) {
	a := makeAlert()
	body := formatBody(a)

	checks := []string{"backup", "/usr/bin/backup.sh", "exit status 1", "disk full", "2s"}
	for _, c := range checks {
		if !strings.Contains(body, c) {
			t.Errorf("body missing %q\nbody: %s", c, body)
		}
	}
}

func TestFormatBody_NoOutputSection(t *testing.T) {
	a := makeAlert()
	a.Output = ""
	body := formatBody(a)
	if strings.Contains(body, "Output:") {
		t.Error("expected no Output section when output is empty")
	}
}

func TestBuildMessage_Headers(t *testing.T) {
	msg := buildMessage("from@example.com", []string{"to@example.com"}, "Test Subject", "body text")
	s := string(msg)
	for _, h := range []string{"From: from@example.com", "To: to@example.com", "Subject: Test Subject"} {
		if !strings.Contains(s, h) {
			t.Errorf("missing header %q in message", h)
		}
	}
}

func TestConfig_Port_Default(t *testing.T) {
	cfg := Config{}
	if cfg.Port() != 587 {
		t.Errorf("expected default port 587, got %d", cfg.Port())
	}
}

func TestConfig_Port_Custom(t *testing.T) {
	cfg := Config{SMTPPort: 465}
	if cfg.Port() != 465 {
		t.Errorf("expected port 465, got %d", cfg.Port())
	}
}
