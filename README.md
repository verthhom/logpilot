# logpilot

> Lightweight CLI for tailing and filtering structured JSON logs from multiple sources

---

## Installation

```bash
go install github.com/youruser/logpilot@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/logpilot.git
cd logpilot
go build -o logpilot .
```

---

## Usage

Tail a single log file and filter by log level:

```bash
logpilot tail --file /var/log/app.log --level error
```

Watch multiple sources and filter by a specific field:

```bash
logpilot tail --file /var/log/app.log --file /var/log/worker.log --filter "service=auth"
```

Pretty-print JSON output:

```bash
logpilot tail --file /var/log/app.log --pretty
```

Pipe from stdin:

```bash
kubectl logs -f my-pod | logpilot tail --level warn --pretty
```

### Flags

| Flag | Description |
|------|-------------|
| `--file` | Path to a log file (repeatable) |
| `--level` | Minimum log level to display (`debug`, `info`, `warn`, `error`) |
| `--filter` | Filter by field value (e.g. `key=value`) |
| `--pretty` | Pretty-print JSON output |

---

## Requirements

- Go 1.21+

---

## License

[MIT](LICENSE) © 2024 youruser