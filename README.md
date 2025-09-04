# Todo Service

A simple **Todo Service** built with **Go**, **PostgreSQL**, and **Hexagonal Architecture** principles.  
The service exposes a REST API to manage `TodoItem` entities and runs inside Docker with automatic database migrations.

---

## Features

- Manage `TodoItem` entities with:
    - `id` (UUID)
    - `description` (string)
    - `dueDate` (timestamp)
- REST API for creating `TodoItem`:
    - `POST /todo`
- PostgreSQL persistence
- Hexagonal architecture (ports & adapters)
- Dependency injection for easier testing
- Automatic database migrations
- Swagger documentation
- Unit tests with mocked repository

---

## Prerequisites

Make sure you have the following installed:

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go 1.24+](https://go.dev/dl/)
- [Make](https://www.gnu.org/software/make/)

---

## Project Structure

```
.
├── cmd/executor        # Application entrypoint
│       └── main.go
    └──  migrations/    # Database migrations
├── internal/           # Application source (hexagonal architecture)
│   ├── domain/         # Entities and business logic
│   ├── ports/          # Inbound and outbound ports
│   └── adapters/       # Adapters (e.g. PostgreSQL, REST)
├── docker-compose.yml  # Docker setup
├── Dockerfile          # Build the Project
├── Makefile            # Automation commands
└── README.md           # Documentation
```

---

## Running the Service

Start the project with Docker:

```bash
make run-docker
```

This will:

- Start PostgresSQL
- Run the Todo service on **port 1212**
- Apply database migrations automatically

---

## API Documentation

Once the service is running, Swagger documentation is available at:

👉 [http://localhost:1212/swagger/index.html](http://localhost:1212/swagger/index.html)

---

## Endpoints

### Create TodoItem
**POST** `/todo-items`

**Request Body:**
```json
{
  "description": "Buy groceries",
  "dueDate": "2025-09-05T18:00:00Z"
}
```

**Response:**
```json
{
  "id": "b1b8f44c-6c2c-4f0c-92e5-9c1a6e8f7c8f",
  "description": "Buy groceries",
  "dueDate": "2025-09-05T18:00:00Z"
}
```

---

## Development

### Install dependencies
```bash
make install
```

### Build the project
```bash
make build
```

### Run locally
```bash
make run
```
### Run With Docker
```bash
make run-docker
```

---

## Database Migrations

Migration files are located in `migrations/`.

To apply migrations manually:
```bash
make migrate
```

---

## Testing

Unit tests mock the repository (no real database needed):

```bash
make test
```

---

## Additional Commands

- **Lint & format code**
  ```bash
  make lint
  ```

- **Generate Swagger docs**
  ```bash
  make doc
  ```

- **Check vulnerabilities**
  ```bash
  make vulncheck
  ```

---

## Tech Stack

- **Language:** Go (1.25+)
- **Database:** PostgreSQL
- **Frameworks/Tools:**
    - `swaggo/swag` (Swagger docs)
    - `testify` (unit testing & mocks)
    - `golangci-lint`, `gci`, `gofumpt` (lint & formatting)
    - `govulncheck` (security scanning)
- **Architecture:** Hexagonal (Ports & Adapters)

---
