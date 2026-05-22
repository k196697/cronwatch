// Package scheduler wires together the config, runner, and notifier to
// execute cron jobs on their defined schedules and report failures.
package scheduler

import (
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"cronwatch/internal/config"
	"cronwatch/internal/notifier"
	"cronwatch/internal/runner"
)

// Scheduler manages cron job lifecycle.
type Scheduler struct {
	cfg     *config.Config
	cron    *cron.Cron
	runner  *runner.Runner
	notify  notifier.Notifier
	mu      sync.Mutex
	entries map[string]cron.EntryID
}

// New creates a Scheduler from the given config, runner, and notifier.
func New(cfg *config.Config, r *runner.Runner, n notifier.Notifier) *Scheduler {
	return &Scheduler{
		cfg:     cfg,
		cron:    cron.New(),
		runner:  r,
		notify:  n,
		entries: make(map[string]cron.EntryID),
	}
}

// Start registers all jobs and begins the cron loop.
func (s *Scheduler) Start() error {
	for _, job := range s.cfg.Jobs {
		job := job // capture
		id, err := s.cron.AddFunc(job.Schedule, func() {
			s.execute(job)
		})
		if err != nil {
			return err
		}
		s.mu.Lock()
		s.entries[job.Name] = id
		s.mu.Unlock()
		log.Printf("scheduler: registered job %q (%s)", job.Name, job.Schedule)
	}
	s.cron.Start()
	return nil
}

// Stop gracefully halts the cron loop.
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("scheduler: stopped")
}

func (s *Scheduler) execute(job config.Job) {
	timeout := time.Duration(job.TimeoutSeconds) * time.Second
	result := s.runner.Run(job.Command, timeout)

	if !result.Success {
		alert := notifier.Alert{
			JobName:   job.Name,
			Command:   job.Command,
			ExitCode:  result.ExitCode,
			Output:    result.Output,
			Duration:  result.Duration,
			TimedOut:  result.TimedOut,
			Timestamp: time.Now(),
		}
		if err := s.notify.Send(alert); err != nil {
			log.Printf("scheduler: failed to send alert for job %q: %v", job.Name, err)
		}
	}
}
