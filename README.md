# Go Boilerplate API

A production-ready Go API boilerplate with Clean Architecture, GORM, Gin, PostgreSQL, Redis, and Swagger.

## 🏗️ Architecture

```
├── cmd/api/             # Application entry point
├── internal/
│   ├── domain/          # Business entities & repository interfaces
│   ├── application/     # Use cases & DTOs
│   ├── infrastructure/  # Database, cache, config implementations
│   └── interfaces/      # HTTP handlers, middleware, routes
├── pkg/                 # Shared packages (logger, response, validator)
├── docs/                # Swagger documentation
└── migrations/          # Database migrations
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make

### Setup

1. **Start services:**

   ```bash
   make docker-up
   ```

2. **Install dependencies:**

   ```bash
   make deps
   ```

3. **Generate Swagger docs:**

   ```bash
   make swagger
   ```

4. **Run application:**

   ```bash
   make run
   # Or with hot reload:
   make dev
   ```

5. **Access API:**
   - API: http://localhost:8080
   - Swagger: http://localhost:8080/swagger/index.html
   - Health: http://localhost:8080/api/v1/health

## 📋 Available Commands

| Command            | Description               |
| ------------------ | ------------------------- |
| `make run`         | Run the application       |
| `make dev`         | Run with hot reload (air) |
| `make test`        | Run tests                 |
| `make lint`        | Run linter                |
| `make swagger`     | Generate Swagger docs     |
| `make docker-up`   | Start Docker services     |
| `make docker-down` | Stop Docker services      |
| `make migrate-up`  | Run migrations            |

## 🛠️ Tech Stack

- **Framework:** Gin
- **ORM:** GORM
- **Database:** PostgreSQL
- **Cache:** Redis
- **Config:** Viper
- **Logger:** Zerolog
- **Docs:** Swagger (swaggo)

## 📖 API Endpoints

| Method | Endpoint         | Description          |
| ------ | ---------------- | -------------------- |
| GET    | `/ping`          | Load balancer health |
| GET    | `/api/v1/health` | Full health check    |
| GET    | `/swagger/*`     | API documentation    |

## 📄 License

MIT
