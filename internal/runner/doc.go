// Package runner provides functionality for executing shell commands
// as part of cron job monitoring. It wraps os/exec to support configurable
// timeouts and captures combined stdout/stderr output along with exit codes
// and execution duration, returning structured Result values for downstream
// alerting and logging.
package runner
