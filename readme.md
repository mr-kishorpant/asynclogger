# Async JSON Logger for Go

## Overview

This library provides a high-performance, thread-safe, **asynchronous logging mechanism** in **JSON format** that writes structured logs to disk. It is optimized for applications with high log volume and is compatible with the **ELK Stack** (Elasticsearch, Logstash, Kibana) or **Filebeat**.

---

## Key Benefits

| Feature                  | Benefit                                                                 |
|--------------------------|-------------------------------------------------------------------------|
| ✅ Asynchronous logging  | Keeps your main application fast and responsive                         |
| 📦 JSON format           | Structured, searchable, and ELK-friendly logs                           |
| 🕒 Date-wise log files   | Automatic log rotation by date                                          |
| 🚀 No manual init needed | Logger auto-initializes on first use                                   |
| 🧹 Graceful shutdown     | Captures OS signals (`SIGINT`, `SIGTERM`) to flush and close cleanly    |
| 🔌 Pluggable             | Easily extendable for Elasticsearch, CloudWatch, etc.                   |

---

## Setup

```bash
go get github.com/your-org/asynclogger
```

> Assumes your project uses Go modules (`go.mod`).

---

## Directory Structure

```
/asynclogger
  ├── logger.go         # Core logic for async JSON logging
  ├── go.mod
  └── logs/             # Generated log files (e.g., 2025-05-05.log)
```

---

## Usage

```go
import "github.com/your-org/asynclogger"

func main() {
    asynclogger.Info("User registered", userID, email)
    asynclogger.Error("Failed DB connection", "retry in", 5, "seconds")
}
```

✅ No manual init  
✅ Auto date-based log file (`logs/YYYY-MM-DD.log`)  
✅ Each line is a structured JSON object

---

## Sample Log Output

```json
{"timestamp":"2025-05-05T15:49:32.891Z","level":"INFO","message":"User registered user123 john@example.com"}
{"timestamp":"2025-05-05T15:49:33.200Z","level":"ERROR","message":"Failed DB connection retry in 5 seconds"}
```

---

## Graceful Shutdown

- On `SIGINT`/`SIGTERM`, logger:
  - Flushes all logs in buffer
  - Closes file handlers safely
  - Prevents data loss

Automatically handled using Go's `os/signal` package.

---

## Extending This Logger

- Add custom fields (`request_id`, `user_id`, etc.) to the `logEntry` struct
- Push logs directly to:
  - Elasticsearch (via REST or Go client)
  - Cloud logging platforms (Loki, CloudWatch)
- Integrate with OpenTelemetry for full observability

---

## Contribution Guidelines

- Maintain non-blocking behavior for logging
- Keep log format NDJSON (1 JSON object per line)
- Follow consistent field naming for structured logging

---

## License

MIT – feel free to use and extend with attribution.
