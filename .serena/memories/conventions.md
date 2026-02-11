# Coding Conventions

## Architecture
- **Clean Architecture**: strict separation of concerns.
- **Dependencies**: `domain` must not depend on other layers. `application` depends on `domain`. `infrastructure` and `interfaces` depend on `application` and `domain`.

## Code Style
- Use explicit error handling with context wrapping.
- Public functions/types must have GoDoc comments.
- API handlers must have Swagger annotations.

## Testing
- Unit tests for domain logic.
- Integration tests (Testcontainers) for infrastructure.
- E2E tests for API endpoints.
