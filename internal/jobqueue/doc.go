// Package jobqueue provides a bounded, thread-safe priority queue for
// scheduling pending cron job execution requests. Jobs with higher priority
// values are dequeued first. The queue rejects pushes when at capacity.
package jobqueue
