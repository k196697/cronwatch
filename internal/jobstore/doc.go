// Package jobstore maintains an in-memory, thread-safe registry of cron job
// execution states, allowing other components to query or update job status
// at any point during the daemon lifecycle.
package jobstore
