# Project Information

## Purpose
A production-ready Go API boilerplate implementing Clean Architecture.
It provides a foundation for building scalable web services with user management, authentication, and other common features.

## Tech Stack
- **Language**: Go 1.24+
- **Framework**: Gin
- **ORM**: GORM (PostgreSQL)
- **Cache**: Redis
- **Config**: Viper
- **Logging**: Zerolog
- **Docs**: Swagger

## Structure
- `cmd/api`: Entry point
- `internal/domain`: Entities and repository interfaces (Pure Go, no external deps)
- `internal/application`: Use cases and DTOs
- `internal/infrastructure`: DB implementations, external services
- `internal/interfaces`: HTTP handlers, routes, middleware
- `pkg`: Shared packages
