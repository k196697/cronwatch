// Package jobpriority provides a thread-safe registry for assigning and
// querying priority levels for cron jobs.
//
// Priority levels range from 1 (Low) to 10 (High). Jobs without an explicit
// assignment default to Normal (5). The registry is safe for concurrent use.
package jobpriority
