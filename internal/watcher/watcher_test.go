package watcher_test

import (
	"testing"
	"time"

	"github.com/user/cronwatch/internal/config"
	"github.com/user/cronwatch/internal/watcher"
)

func makeConfig() *config.Config {
	return &config.Config{
		Jobs: []config.Job{
			{
				Name:     "test-job",
				Schedule: "@every 1m",
				Command:  "echo hello",
				Timeout:  30,
			},
		},
		SMTP: config.SMTP{
			Host: "",
		},
	}
}

func TestNew_ReturnsWatcher(t *testing.T) {
	cfg := makeConfig()
	w, err := watcher.New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestRun_StopsOnSignal(t *testing.T) {
	cfg := makeConfig()
	w, err := watcher.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- w.Run()
	}()

	// Give it a moment to start, then stop it.
	time.Sleep(50 * time.Millisecond)
	w.Stop()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("expected nil error on clean stop, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("watcher did not stop within timeout")
	}
}

func TestStop_IdempotentDoesNotPanic(t *testing.T) {
	cfg := makeConfig()
	w, err := watcher.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Stopping before Run should not panic.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Stop panicked: %v", r)
		}
	}()
	w.Stop()
}
