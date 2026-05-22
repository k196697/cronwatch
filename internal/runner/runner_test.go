package runner_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/cronwatch/internal/runner"
)

func TestRun_Success(t *testing.T) {
	r := runner.New()
	res := r.Run(context.Background(), "echo-job", "echo hello", 0)

	if res.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", res.ExitCode)
	}
	if !strings.Contains(res.Output, "hello") {
		t.Errorf("expected output to contain 'hello', got %q", res.Output)
	}
	if res.TimedOut {
		t.Error("expected TimedOut to be false")
	}
	if res.Err != nil {
		t.Errorf("unexpected error: %v", res.Err)
	}
}

func TestRun_Failure(t *testing.T) {
	r := runner.New()
	res := r.Run(context.Background(), "fail-job", "exit 2", 0)

	if res.ExitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", res.ExitCode)
	}
	if res.Err == nil {
		t.Error("expected an error for non-zero exit")
	}
}

func TestRun_Timeout(t *testing.T) {
	r := runner.New()
	res := r.Run(context.Background(), "slow-job", "sleep 10", 100*time.Millisecond)

	if !res.TimedOut {
		t.Error("expected TimedOut to be true")
	}
	if res.ExitCode != -1 {
		t.Errorf("expected exit code -1 on timeout, got %d", res.ExitCode)
	}
}

func TestRun_CapturesOutput(t *testing.T) {
	r := runner.New()
	res := r.Run(context.Background(), "output-job", "echo stdout && echo stderr >&2", 0)

	if !strings.Contains(res.Output, "stdout") {
		t.Errorf("expected stdout in output, got %q", res.Output)
	}
	if !strings.Contains(res.Output, "stderr") {
		t.Errorf("expected stderr in combined output, got %q", res.Output)
	}
}

func TestRun_Duration(t *testing.T) {
	r := runner.New()
	res := r.Run(context.Background(), "dur-job", "sleep 0.05", 0)

	if res.Duration < 50*time.Millisecond {
		t.Errorf("expected duration >= 50ms, got %v", res.Duration)
	}
}
