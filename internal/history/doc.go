// Package history provides an in-memory, concurrency-safe store for
// recording and querying cron job execution results, including helpers
// for detecting consecutive failures.
//
// Basic usage:
//
//	store := history.NewStore()
//	store.Record(history.Entry{
//		JobName:   "backup",
//		StartedAt: time.Now(),
//		Success:   false,
//		Output:    "connection refused",
//	})
//
//	// Check if the job has failed too many times in a row.
//	if store.ConsecutiveFailures("backup") >= 3 {
//		// trigger an alert
//	}
package history
