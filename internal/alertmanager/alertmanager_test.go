package alertmanager_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/alertmanager"
	"github.com/example/cronwatch/internal/notifier"
)

type mockNotifier struct {
	calls atomic.Int32
	fail  bool
}

func (m *mockNotifier) Send(_ notifier.Alert) error {
	m.calls.Add(1)
	if m.fail {
		return errors.New("smtp error")
	}
	return nil
}

func makeAlert(job string) alertmanager.Alert {
	return alertmanager.Alert{
		JobName:   job,
		Command:   "/usr/bin/backup",
		Error:     "exit status 1",
		Duration:  2 * time.Second,
		Timestamp: time.Now(),
	}
}

func TestSend_ForwardsAlert(t *testing.T) {
	mn := &mockNotifier{}
	mgr := alertmanager.New(mn, time.Minute)

	sent, err := mgr.Send(makeAlert("backup"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sent {
		t.Fatal("expected alert to be sent")
	}
	if mn.calls.Load() != 1 {
		t.Fatalf("expected 1 notifier call, got %d", mn.calls.Load())
	}
}

func TestSend_SuppressesDuringCooldown(t *testing.T) {
	mn := &mockNotifier{}
	mgr := alertmanager.New(mn, time.Hour)

	mgr.Send(makeAlert("backup")) //nolint:errcheck
	sent, err := mgr.Send(makeAlert("backup"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sent {
		t.Fatal("expected alert to be suppressed during cooldown")
	}
	if mn.calls.Load() != 1 {
		t.Fatalf("expected 1 notifier call, got %d", mn.calls.Load())
	}
}

func TestSend_DifferentJobsNotSuppressed(t *testing.T) {
	mn := &mockNotifier{}
	mgr := alertmanager.New(mn, time.Hour)

	mgr.Send(makeAlert("job-a")) //nolint:errcheck
	sent, err := mgr.Send(makeAlert("job-b"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sent {
		t.Fatal("expected alert for different job to be sent")
	}
}

func TestReset_ClearsCooldown(t *testing.T) {
	mn := &mockNotifier{}
	mgr := alertmanager.New(mn, time.Hour)

	mgr.Send(makeAlert("backup")) //nolint:errcheck
	mgr.Reset("backup")
	sent, err := mgr.Send(makeAlert("backup"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sent {
		t.Fatal("expected alert to be sent after reset")
	}
}

func TestSend_PropagatesNotifierError(t *testing.T) {
	mn := &mockNotifier{fail: true}
	mgr := alertmanager.New(mn, time.Minute)

	_, err := mgr.Send(makeAlert("backup"))
	if err == nil {
		t.Fatal("expected error from notifier")
	}
}
