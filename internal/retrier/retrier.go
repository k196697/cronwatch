// Package retrier provides retry logic for cron job execution.
package retrier

import (
	"context"
	"fmt"
	"time"
)

// Result holds the outcome of a retried operation.
type Result struct {
	Attempts int
	Err      error
	Output   string
}

// Config holds retry configuration.
type Config struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64 // backoff multiplier; 1.0 = constant delay
}

// Func is the function signature for a retryable operation.
type Func func(ctx context.Context) (string, error)

// Retrier executes a function with retry logic.
type Retrier struct {
	cfg Config
}

// New creates a new Retrier with the given config.
// MaxAttempts defaults to 1 if zero or negative.
func New(cfg Config) *Retrier {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}
	if cfg.Multiplier <= 0 {
		cfg.Multiplier = 1.0
	}
	return &Retrier{cfg: cfg}
}

// Run executes fn up to MaxAttempts times, backing off between attempts.
// It returns a Result containing the number of attempts made and the final error.
func (r *Retrier) Run(ctx context.Context, fn Func) Result {
	var (
		lastErr error
		output  string
		delay   = r.cfg.Delay
	)

	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return Result{Attempts: attempt - 1, Err: ctx.Err()}
		default:
		}

		output, lastErr = fn(ctx)
		if lastErr == nil {
			return Result{Attempts: attempt, Err: nil, Output: output}
		}

		if attempt < r.cfg.MaxAttempts {
			select {
			case <-ctx.Done():
				return Result{Attempts: attempt, Err: ctx.Err()}
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * r.cfg.Multiplier)
		}
	}

	return Result{
		Attempts: r.cfg.MaxAttempts,
		Err:      fmt.Errorf("all %d attempt(s) failed: %w", r.cfg.MaxAttempts, lastErr),
		Output:   output,
	}
}
