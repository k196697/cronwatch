// Package notifier implements alert delivery for cronwatch.
//
// It formats job failure details into human-readable messages and
// dispatches them via SMTP. Additional backends (webhook, Slack, etc.)
// can be added by extending the Notifier interface.
package notifier
