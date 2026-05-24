package scheduler

import (
	"context"
	"log"
	"sync"

	"github.com/robfig/cron/v3"

	"github.com/example/cronwatch/internal/alertmanager"
	"github.com/example/cronwatch/internal/config"
	"github.com/example/cronwatch/internal/jobqueue"
	"github.com/example/cronwatch/internal/runner"
)

// Scheduler manages cron job registration and execution.
type Scheduler struct {
	cfg   *config.Config
	cron  *cron.Cron
	alert *alertmanager.Manager
	queue *jobqueue.Queue
	mu    sync.Mutex
}

// New creates a Scheduler from the provided config and alert manager.
func New(cfg *config.Config, alert *alertmanager.Manager) *Scheduler {
	q, _ := jobqueue.New(256)
	return &Scheduler{
		cfg:   cfg,
		cron:  cron.New(),
		alert: alert,
		queue: q,
	}
}

// Start registers all configured jobs and begins the cron loop.
func (s *Scheduler) Start(ctx context.Context) error {
	for _, job := range s.cfg.Jobs {
		job := job
		_, err := s.cron.AddFunc(job.Schedule, func() {
			_ = s.queue.Push(jobqueue.Entry{JobName: job.Name, Priority: job.Priority})
			s.execute(ctx, job)
		})
		if err != nil {
			return err
		}
	}
	s.cron.Start()
	go func() {
		<-ctx.Done()
		s.Stop()
	}()
	return nil
}

// Stop halts the cron scheduler gracefully.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cron.Stop()
}

// execute runs a single job and sends an alert on failure.
func (s *Scheduler) execute(ctx context.Context, job config.Job) {
	r := runner.New(job)
	result := r.Run(ctx)
	if !result.Success {
		if err := s.alert.Send(alertmanager.Alert{
			JobName: job.Name,
			Error:   result.Error,
			Output:  result.Output,
		}); err != nil {
			log.Printf("scheduler: alert send failed for job %s: %v", job.Name, err)
		}
	}
}
