// Package jobpause provides a thread-safe store for tracking paused cron jobs.
//
// A job can be paused indefinitely or for a fixed duration. Timed pauses
// expire automatically when IsPaused is called after the deadline has passed.
package jobpause
