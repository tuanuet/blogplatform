# Suggested Commands

## Build & Run
- `make run`: Run the application
- `make dev`: Run with hot reload (air)
- `make build`: Build the application binary

## Testing & Quality
- `make test`: Run unit tests
- `make test-coverage`: Run tests with coverage report
- `make lint`: Run golangci-lint

## Docker
- `make docker-up`: Start PostgreSQL and Redis
- `make docker-down`: Stop services

## Database
- `make migrate-up`: Run migrations
- `make migrate-down`: Rollback migrations
- `make migrate-create name=migration_name`: Create a new migration file

## Documentation
- `make swagger`: Generate Swagger API docs
