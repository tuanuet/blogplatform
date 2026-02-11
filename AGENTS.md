# AI Agent

This file defines the project context, coding standards, and agent configurations for the OpenCode AI system. It serves as the primary source of truth for all AI agents working on this repository.

## Project Overview

**Go Boilerplate API** is a production-ready Go API boilerplate implementing Clean Architecture.
It uses **GORM**, **Gin**, **PostgreSQL**, **Redis**, and **Swagger** to provide a robust foundation for building scalable web services.

## Tech Stack & Standards

- **Language**: Go 1.24+
- **Web Framework**: Gin (`github.com/gin-gonic/gin`)
- **ORM**: GORM (`gorm.io/gorm`)
- **Database**: PostgreSQL (Driver: `pgx`)
- **Cache**: Redis (`go-redis/v9`)
- **Configuration**: Viper
- **Logging**: Zerolog
- **Documentation**: Swagger (swaggo)
- **Testing**: Testcontainers, testify

## Project Structure

The project follows a **Clean Architecture** layout:

- `cmd/api/`: Application entry point (main.go).
- `internal/`: Private application logic.
  - `domain/`: Business entities & repository interfaces.
  - `application/`: Use cases & DTOs.
  - `infrastructure/`: Database, cache, external service implementations.
  - `interfaces/`: HTTP handlers, middleware, routes.
- `pkg/`: Shared public packages (logger, response, validator).
- `docs/`: Swagger documentation.
- `migrations/`: Database migrations.

## Coding Conventions

- **Architecture**: Strictly adhere to Clean Architecture boundaries. `domain` should have no external dependencies.
- **Error Handling**: Use explicit error handling. Wrap errors with context where appropriate.
- **Testing**:
  - Unit tests for domain logic.
  - Integration tests (using Testcontainers) for infrastructure.
  - End-to-end tests for API endpoints.
- **Documentation**: All public functions and types must have GoDoc comments. API handlers must have Swagger annotations.

---

### Tooling Configuration

#### Serena MCP (High Priority)

**Repository**: [oraios/serena](https://github.com/oraios/serena)
**Description**: Provides deep code intelligence via LSP (Symbol resolution, references, etc.).

**Usage Rule**: Always prefer Serena's semantic tools over basic `grep`/`glob` for code navigation and understanding.

---

**Instructions**:

- Do NOT preemptively load all references - use lazy loading based on actual need.
- When loaded, treat content as mandatory instructions that override defaults.
- Follow references recursively when needed.
