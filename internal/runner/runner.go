package runner

import (
	"context"
	"os/exec"
	"time"
)

// Result holds the outcome of a single cron job execution.
type Result struct {
	JobName  string
	Command  string
	ExitCode int
	Output   string
	Duration time.Duration
	TimedOut bool
	Err      error
}

// Runner executes shell commands with an optional timeout.
type Runner struct{}

// New returns a new Runner instance.
func New() *Runner {
	return &Runner{}
}

// Run executes the given command string under a shell with the provided timeout.
// A timeout of zero means no timeout is enforced.
func (r *Runner) Run(ctx context.Context, jobName, command string, timeout time.Duration) Result {
	result := Result{
		JobName: jobName,
		Command: command,
	}

	runCtx := ctx
	var cancel context.CancelFunc
	if timeout > 0 {
		runCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	start := time.Now()
	cmd := exec.CommandContext(runCtx, "sh", "-c", command)
	out, err := cmd.CombinedOutput()
	result.Duration = time.Since(start)
	result.Output = string(out)

	if runCtx.Err() == context.DeadlineExceeded {
		result.TimedOut = true
		result.ExitCode = -1
		result.Err = runCtx.Err()
		return result
	}

	if err != nil {
		result.Err = err
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		return result
	}

	result.ExitCode = 0
	return result
}
