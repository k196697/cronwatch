// Package watcher provides the top-level daemon lifecycle management for
// cronwatch. It wires together the scheduler and notifier, starts the cron
// engine, and handles OS signals for graceful shutdown.
package watcher
