package circuit

import (
	"testing"
	"time"
)

func newBreaker(threshold int) *Breaker {
	return New(threshold, 50*time.Millisecond)
}

func TestAllow_ClosedByDefault(t *testing.T) {
	b := newBreaker(3)
	if !b.Allow() {
		t.Fatal("expected Allow() == true for a new breaker")
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	b := newBreaker(3)
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if b.State() != StateOpen {
		t.Fatalf("expected StateOpen after %d failures, got %s", 3, b.State())
	}
	if b.Allow() {
		t.Fatal("expected Allow() == false when circuit is open")
	}
}

func TestRecordSuccess_ResetsClosed(t *testing.T) {
	b := newBreaker(2)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatal("expected circuit to be open")
	}
	// wait for reset window
	time.Sleep(60 * time.Millisecond)
	b.Allow() // transitions to half-open
	b.RecordSuccess()
	if b.State() != StateClosed {
		t.Fatalf("expected StateClosed after success, got %s", b.State())
	}
	if !b.Allow() {
		t.Fatal("expected Allow() == true after reset")
	}
}

func TestHalfOpen_AfterResetWindow(t *testing.T) {
	b := newBreaker(1)
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatal("expected StateOpen")
	}
	time.Sleep(60 * time.Millisecond)
	if !b.Allow() {
		t.Fatal("expected Allow() == true after reset window")
	}
	if b.State() != StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %s", b.State())
	}
}

func TestHalfOpen_FailureReopens(t *testing.T) {
	b := newBreaker(1)
	b.RecordFailure()
	time.Sleep(60 * time.Millisecond)
	b.Allow() // half-open
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatalf("expected StateOpen after failure in half-open, got %s", b.State())
	}
}

func TestNew_DefaultsOnInvalidArgs(t *testing.T) {
	b := New(0, 0)
	if b.threshold != 3 {
		t.Errorf("expected default threshold 3, got %d", b.threshold)
	}
	if b.resetAfter != 30*time.Second {
		t.Errorf("expected default resetAfter 30s, got %v", b.resetAfter)
	}
}

func TestState_String(t *testing.T) {
	cases := map[State]string{
		StateClosed:   "closed",
		StateOpen:     "open",
		StateHalfOpen: "half-open",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("State(%d).String() = %q, want %q", int(s), got, want)
		}
	}
}
