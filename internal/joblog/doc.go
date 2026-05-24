// Package joblog provides a per-job in-memory log store for cronwatch.
//
// Each job accumulates a bounded ring of log entries (info, warn, error)
// produced during execution. Entries are safe for concurrent access and
// can be retrieved or cleared at any time.
package joblog
