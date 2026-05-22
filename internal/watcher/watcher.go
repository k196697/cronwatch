package watcher

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/cronwatch/internal/config"
	"github.com/user/cronwatch/internal/notifier"
	"github.com/user/cronwatch/internal/scheduler"
)

// Watcher ties together the scheduler and notifier, managing the daemon lifecycle.
type Watcher struct {
	cfg       *config.Config
	sched     *scheduler.Scheduler
	notify    notifier.Notifier
	logger    *log.Logger
	stopCh    chan struct{}
}

// New creates a new Watcher from the provided config.
func New(cfg *config.Config) (*Watcher, error) {
	n, err := notifier.New(cfg)
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "[cronwatch] ", log.LstdFlags)
	s := scheduler.New(cfg, n, logger)

	return &Watcher{
		cfg:    cfg,
		sched:  s,
		notify: n,
		logger: logger,
		stopCh: make(chan struct{}),
	}, nil
}

// Run starts the watcher and blocks until a termination signal is received.
func (w *Watcher) Run() error {
	w.logger.Printf("starting cronwatch with %d job(s)", len(w.cfg.Jobs))

	if err := w.sched.Start(); err != nil {
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		w.logger.Printf("received signal %s, shutting down", sig)
	case <-w.stopCh:
		w.logger.Println("stop requested, shutting down")
	}

	w.sched.Stop()
	w.logger.Println("cronwatch stopped")
	return nil
}

// Stop signals the watcher to shut down gracefully.
func (w *Watcher) Stop() {
	close(w.stopCh)
}
