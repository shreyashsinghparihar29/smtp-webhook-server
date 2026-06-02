# SMTP Webhook Server

SMTP Webhook Server is a lightweight event-driven SMTP server written in Go that receives incoming emails and forwards them to a configured HTTP webhook endpoint as a structured JSON payload.

The system is designed with **reliability and scalability in mind**, supporting asynchronous processing, persistent storage, and retry mechanisms.

---

## 🚀 Features

- 📩 Receive emails through SMTP
- 🌐 Forward email content to any HTTP webhook
- ⚙️ Configurable SMTP listener address
- 🔗 Configurable webhook endpoint
- 🗄️ SQLite-based persistent email storage
- 🧵 Asynchronous worker queue for processing emails
- 🔁 Retry mechanism with exponential backoff
- 💀 Dead-letter handling for failed deliveries
- 🐳 Docker and native deployment support
- 🧪 Simple local testing workflow

---

## 🏗️ System Architecture

🔹 Core Flow : 
SMTP -> Parser -> SQLite Storage -> Queue -> Worker Pool -> Webhook

🔹 Failure Handling Flow : 
Webhook Failure ->
Retry Engine (Exponential Backoff) ->
Dead Letter Queue (FAILED status in DB)


## 📡 Scalable Design (Kafka-ready architecture)

This system is designed in a way that can be extended to distributed systems:

- Instead of in-memory queue → Kafka can be used
- Workers can be horizontally scaled
- Email ingestion becomes event-driven pipeline

```

SMTP → Kafka Topic → Worker Consumers → Webhook Services

```

---

## ⚙️ Development

### Install dependencies

```bash
go mod vendor
```

### Build the application

```bash
go build -o smtp-webhook-server .
```

---

## 🐳 Development with Docker

### Build development image

```bash
docker build -f Dockerfile.dev -t smtp-webhook-server-dev .
```

### Run container

```bash
docker run -p 2525:2525 smtp-webhook-server-dev \
  --webhook=http://host.docker.internal:8080/my/webhook
```

---

## 🚀 Production Docker

### Build production image

```bash
docker build -t smtp-webhook-server .
```

### Run container

```bash
docker run -p 2525:2525 smtp-webhook-server \
  --webhook=http://host.docker.internal:8080/my/webhook
```

---

## 💻 Native Usage

### Run directly

```bash
smtp-webhook-server \
  --listen=:2525 \
  --webhook=http://localhost:8080/my/webhook
```

### View all options

```bash
smtp-webhook-server --help
```

---

## 🧪 Local Testing

### 1. Start example webhook

```bash
go run examples/webhook.go
```

Expected output:

```
Listening on :8080
```

---

### 2. Start SMTP server

```bash
go run . \
  --listen=:2525 \
  --webhook=http://localhost:8080/my/webhook \
  --timeout.read=60 \
  --timeout.write=60
```

---

### 3. Connect using Telnet

```bash
telnet localhost 2525
```

---

### 4. Send test email

```
HELO localhost
MAIL FROM:<test@example.com>
RCPT TO:<receiver@example.com>
DATA
Subject: Test

Hello SMTP webhook server
.
QUIT
```

---

### 5. Verify system behavior

When email is processed:

- Email is stored in SQLite
- Job is queued to worker pool
- Worker processes webhook delivery
- Retry happens automatically on failure
- Failed messages go into FAILED state (dead letter)

---

## 🗄️ Database Schema (SQLite)

```sql
CREATE TABLE emails (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  message_id TEXT,
  from_addr TEXT,
  to_addr TEXT,
  subject TEXT,
  body TEXT,
  status TEXT,
  retry_count INTEGER DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## 🔁 Retry System Behavior

- Max retries: 3
- Exponential backoff
- On success → status = SUCCESS
- On failure → status = FAILED

---

## 📌 Notes

- Webhook must be running before SMTP server sends data
- Use `host.docker.internal` for Docker on Windows/macOS
- Port `2525` used to avoid root permissions
- System is designed to be easily extended to Kafka-based architecture
```

---

