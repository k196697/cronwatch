package scheduler

import (
	"testing"
	"time"

	"cronwatch/internal/config"
	"cronwatch/internal/notifier"
	"cronwatch/internal/runner"
)

// stubNotifier records the last alert it received.
type stubNotifier struct {
	sentAlerts []notifier.Alert
}

func (s *stubNotifier) Send(a notifier.Alert) error {
	s.sentAlerts = append(s.sentAlerts, a)
	return nil
}

func makeConfig(schedule, command string, timeoutSec int) *config.Config {
	return &config.Config{
		Jobs: []config.Job{
			{
				Name:           "test-job",
				Schedule:       schedule,
				Command:        command,
				TimeoutSeconds: timeoutSec,
			},
		},
	}
}

func TestStart_RegistersJobs(t *testing.T) {
	cfg := makeConfig("@every 1h", "echo hello", 5)
	r := runner.New()
	n := &stubNotifier{}
	s := New(cfg, r, n)

	if err := s.Start(); err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer s.Stop()

	s.mu.Lock()
	_, ok := s.entries["test-job"]
	s.mu.Unlock()

	if !ok {
		t.Error("expected test-job to be registered in entries map")
	}
}

func TestStart_InvalidSchedule(t *testing.T) {
	cfg := makeConfig("not-a-cron", "echo hi", 5)
	r := runner.New()
	n := &stubNotifier{}
	s := New(cfg, r, n)

	if err := s.Start(); err == nil {
		t.Error("expected error for invalid schedule, got nil")
		s.Stop()
	}
}

func TestExecute_SendsAlertOnFailure(t *testing.T) {
	cfg := makeConfig("@every 1h", "false", 5)
	r := runner.New()
	n := &stubNotifier{}
	s := New(cfg, r, n)

	s.execute(cfg.Jobs[0])

	if len(n.sentAlerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(n.sentAlerts))
	}
	if n.sentAlerts[0].JobName != "test-job" {
		t.Errorf("unexpected job name: %s", n.sentAlerts[0].JobName)
	}
}

func TestExecute_NoAlertOnSuccess(t *testing.T) {
	cfg := makeConfig("@every 1h", "true", 5)
	r := runner.New()
	n := &stubNotifier{}
	s := New(cfg, r, n)

	s.execute(cfg.Jobs[0])

	if len(n.sentAlerts) != 0 {
		t.Errorf("expected no alerts, got %d", len(n.sentAlerts))
	}
}

func TestExecute_TimeoutSendsAlert(t *testing.T) {
	cfg := makeConfig("@every 1h", "sleep 10", 1)
	r := runner.New()
	n := &stubNotifier{}
	s := New(cfg, r, n)

	start := time.Now()
	s.execute(cfg.Jobs[0])
	elapsed := time.Since(start)

	if elapsed > 3*time.Second {
		t.Errorf("execute took too long: %v", elapsed)
	}
	if len(n.sentAlerts) != 1 || !n.sentAlerts[0].TimedOut {
		t.Error("expected a timed-out alert")
	}
}
