# cronwatch

Lightweight daemon that monitors cron job execution and sends alerts on failures or timeouts.

## Installation

```bash
go install github.com/yourname/cronwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/cronwatch.git && cd cronwatch && make build
```

## Usage

Define your monitored jobs in a YAML config file:

```yaml
# cronwatch.yaml
jobs:
  - name: daily-backup
    schedule: "0 2 * * *"
    timeout: 30m
    alert:
      email: ops@example.com

  - name: hourly-sync
    schedule: "0 * * * *"
    timeout: 5m
    alert:
      slack: "#alerts"
```

Start the daemon:

```bash
cronwatch --config cronwatch.yaml
```

Wrap an existing cron job to report its status:

```bash
# In your crontab
0 2 * * * cronwatch exec --job daily-backup -- /usr/local/bin/backup.sh
```

cronwatch will send an alert if the job exits with a non-zero status or exceeds the configured timeout.

## Configuration

| Field | Description | Default |
|-------|-------------|---------|
| `timeout` | Max allowed runtime | `1h` |
| `retries` | Attempts before alerting | `0` |
| `alert.email` | Alert recipient email | — |
| `alert.slack` | Slack channel for alerts | — |

## License

MIT © 2024 yourname