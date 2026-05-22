package ratelimit_test

import (
	"testing"
	"time"

	"github.com/cronwatch/internal/ratelimit"
)

func TestAllow_PermitsUpToMax(t *testing.T) {
	l := ratelimit.New(time.Minute, 3)

	for i := 0; i < 3; i++ {
		if !l.Allow("job-a") {
			t.Fatalf("expected call %d to be allowed", i+1)
		}
	}

	if l.Allow("job-a") {
		t.Fatal("expected 4th call to be denied")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	l := ratelimit.New(time.Minute, 1)

	if !l.Allow("job-a") {
		t.Fatal("expected job-a to be allowed")
	}
	if !l.Allow("job-b") {
		t.Fatal("expected job-b to be allowed independently")
	}
	if l.Allow("job-a") {
		t.Fatal("expected job-a second call to be denied")
	}
}

func TestAllow_WindowExpiry(t *testing.T) {
	l := ratelimit.New(50*time.Millisecond, 1)

	if !l.Allow("job-a") {
		t.Fatal("expected first call to be allowed")
	}
	if l.Allow("job-a") {
		t.Fatal("expected second call to be denied within window")
	}

	time.Sleep(60 * time.Millisecond)

	if !l.Allow("job-a") {
		t.Fatal("expected call after window expiry to be allowed")
	}
}

func TestReset_ClearsHistory(t *testing.T) {
	l := ratelimit.New(time.Minute, 1)

	l.Allow("job-a")
	if l.Allow("job-a") {
		t.Fatal("expected second call to be denied")
	}

	l.Reset("job-a")

	if !l.Allow("job-a") {
		t.Fatal("expected call after reset to be allowed")
	}
}

func TestRemaining_DecrementsOnAllow(t *testing.T) {
	l := ratelimit.New(time.Minute, 3)

	if got := l.Remaining("job-a"); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}

	l.Allow("job-a")
	if got := l.Remaining("job-a"); got != 2 {
		t.Fatalf("expected 2 remaining, got %d", got)
	}

	l.Allow("job-a")
	l.Allow("job-a")
	if got := l.Remaining("job-a"); got != 0 {
		t.Fatalf("expected 0 remaining, got %d", got)
	}
}
