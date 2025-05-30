# Interactive Polling Platform

A highly scalable and performant polling platform for web and mobile applications. This system allows users to interact with polls in a vertical feed with smooth UX and robust backend support for concurrency, caching, and observability.

---

## 🚀 Features

- Vote or skip polls with real-time feedback
- Filter and search polls by tags
- Enforced daily voting limits (100 votes/day per user)
- Feed excludes previously voted/skipped polls
- High-throughput ingestion and processing using event streaming
- Inbox/Outbox pattern to ensure consistency and eventual delivery
- Extensible hexagonal architecture with protocol-agnostic processing

---

## 🧠 Architecture Overview

The platform follows **Hexagonal Architecture (Ports & Adapters)** to promote separation of concerns and adaptability to multiple interfaces (HTTP, Kafka, etc).

### ✅ Core Components

- **API Layer**: Handles HTTP requests for creating polls, voting, skipping, retrieving stats, and fetching feed.
- **Application Layer**: Business logic that enforces voting limits, handles poll visibility, and manages core workflows.
- **Infrastructure Layer**:
  - **PostgreSQL**: Persistent storage for polls, votes, skips, and outbox messages.
  - **Redis**: Caching for hot data such as popular polls and aggregated stats.
  - **Kafka**: Asynchronous event streaming for vote and poll-related events.
  - **Debezium**: CDC tool that reads the `outbox` table and publishes messages to Kafka.
  - **Inbox Table**: Ensures idempotent message consumption across multiple service instances.

---

## 📦 Tech Stack

- **Golang**: Backend service logic
- **PostgreSQL**: Primary relational database
- **Redis**: In-memory cache for performance optimization
- **Kafka**: Event-driven processing
- **Debezium**: Change Data Capture (CDC) for reliable Kafka integration

---

## 🧰 Design Highlights

### Inbox/Outbox Pattern

- **Outbox**: All important write-side events (like votes) are saved into the `outbox` table. Debezium reads from here and forwards them to Kafka, decoupling our core logic from Kafka’s availability.
- **Inbox**: Incoming Kafka messages are checked against the `inbox` table to ensure **exactly-once processing**, even with **multi-instance deployments**.
- This architecture supports **high availability** and **idempotent message handling** across service restarts.

### Decoupled Processing

- The core logic is decoupled from transport protocols. This means poll voting logic can be triggered by **HTTP**, **Kafka**, or any other future adapters, without modifying the core logic.

---

## 🔥 Performance & Scale

- Write-heavy operations like voting and skipping are fast due to async processing and Redis caching.
- Feed fetching is optimized via indexed queries and Redis-backed hot poll sets.
- Stress-tested with increasing RPS to observe response time degradation.
- Kafka Consumer Groups ensure scalability across multiple service replicas.

---

## 📊 Observability

- `/metrics` endpoint exposes:
  - HTTP request latencies
  - DB query durations
  - Cache hits/misses
  - Kafka publish and consumption rates

---

## 🧪 Testing

- **Unit Tests**:
  - Vote limit enforcement
  - Feed exclusion logic
  - Outbox record generation
- **Integration Tests**:
  - Poll creation to Kafka publishing via outbox & Debezium
  - End-to-end vote + aggregation flow
- **Performance Tests**:
  - k6/Go benchmarks to evaluate system under increasing load
  - Includes RPS vs. latency report

---

## 🐳 Docker Compose

```bash
docker-compose up --build -d
