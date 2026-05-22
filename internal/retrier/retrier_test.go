package retrier_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cronwatch/internal/retrier"
)

var errFake = errors.New("fake error")

func TestRun_SuccessOnFirstAttempt(t *testing.T) {
	r := retrier.New(retrier.Config{MaxAttempts: 3, Delay: 0})
	res := r.Run(context.Background(), func(_ context.Context) (string, error) {
		return "ok", nil
	})
	if res.Err != nil {
		t.Fatalf("expected no error, got %v", res.Err)
	}
	if res.Attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", res.Attempts)
	}
	if res.Output != "ok" {
		t.Fatalf("expected output 'ok', got %q", res.Output)
	}
}

func TestRun_RetriesOnFailure(t *testing.T) {
	calls := 0
	r := retrier.New(retrier.Config{MaxAttempts: 3, Delay: 0})
	res := r.Run(context.Background(), func(_ context.Context) (string, error) {
		calls++
		if calls < 3 {
			return "", errFake
		}
		return "done", nil
	})
	if res.Err != nil {
		t.Fatalf("expected success on 3rd attempt, got %v", res.Err)
	}
	if res.Attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", res.Attempts)
	}
}

func TestRun_ExhaustsAttempts(t *testing.T) {
	r := retrier.New(retrier.Config{MaxAttempts: 2, Delay: 0})
	res := r.Run(context.Background(), func(_ context.Context) (string, error) {
		return "", errFake
	})
	if res.Err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
	if res.Attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", res.Attempts)
	}
}

func TestRun_RespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	r := retrier.New(retrier.Config{MaxAttempts: 5, Delay: 10 * time.Millisecond})
	res := r.Run(ctx, func(_ context.Context) (string, error) {
		return "", errFake
	})
	if res.Err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestRun_DefaultsMaxAttemptsToOne(t *testing.T) {
	r := retrier.New(retrier.Config{MaxAttempts: 0})
	calls := 0
	r.Run(context.Background(), func(_ context.Context) (string, error) {
		calls++
		return "", errFake
	})
	if calls != 1 {
		t.Fatalf("expected exactly 1 call, got %d", calls)
	}
}
