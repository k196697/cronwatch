package metrics

import (
	"sync"
	"time"
)

// JobMetrics holds aggregated runtime statistics for a single cron job.
type JobMetrics struct {
	JobName      string
	TotalRuns    int
	SuccessCount int
	FailureCount int
	LastRun      time.Time
	LastDuration time.Duration
	AvgDuration  time.Duration
	totalNanos   int64
}

// Collector accumulates metrics across all monitored jobs.
type Collector struct {
	mu   sync.RWMutex
	jobs map[string]*JobMetrics
}

// New returns an initialised Collector.
func New() *Collector {
	return &Collector{
		jobs: make(map[string]*JobMetrics),
	}
}

// Record updates metrics for the named job after a run completes.
func (c *Collector) Record(name string, success bool, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	m, ok := c.jobs[name]
	if !ok {
		m = &JobMetrics{JobName: name}
		c.jobs[name] = m
	}

	m.TotalRuns++
	m.LastRun = time.Now()
	m.LastDuration = duration
	m.totalNanos += duration.Nanoseconds()
	m.AvgDuration = time.Duration(m.totalNanos / int64(m.TotalRuns))

	if success {
		m.SuccessCount++
	} else {
		m.FailureCount++
	}
}

// Get returns a copy of the metrics for the named job and a boolean
// indicating whether the job has been seen before.
func (c *Collector) Get(name string) (JobMetrics, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	m, ok := c.jobs[name]
	if !ok {
		return JobMetrics{}, false
	}
	return *m, true
}

// All returns a snapshot of metrics for every job seen so far.
func (c *Collector) All() []JobMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make([]JobMetrics, 0, len(c.jobs))
	for _, m := range c.jobs {
		out = append(out, *m)
	}
	return out
}
