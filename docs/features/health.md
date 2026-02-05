# Health Feature

## Overview
The Health feature provides system monitoring capabilities, allowing external services (like load balancers or orchestrators) and administrators to verify the operational status of the application and its critical dependencies.

## Architecture
- **UseCase**: `HealthUseCase` serves as the entry point for health checks.
- **Domain Service**: Delegates the actual verification logic to `domainService.SystemService`, which interacts with infrastructure components.

## Core Logic & Features
### System Health Check
The `Check` operation performs a real-time status check of the system.
- **Dependency Verification**: Checks the connectivity and status of critical infrastructure components:
  - **Database**: Verifies connection to the primary data store.
  - **Redis**: Verifies connection to the caching/messaging layer.
- **Status Aggregation**: Returns an overall system status (`OK` or `ERROR`) along with granular status for each service.
- **Timestamping**: Includes the server timestamp of the check.

## Data Model
The feature uses DTOs to structure the health report.

### HealthResponse
```go
type HealthResponse struct {
    Status    string        // "OK" or "ERROR"
    Timestamp time.Time     // Time of check
    Services  ServiceHealth // Detailed status
}
```

### ServiceHealth
```go
type ServiceHealth struct {
    Database ServiceStatus // "connected" or "disconnected"
    Redis    ServiceStatus // "connected" or "disconnected"
}
```

## API Reference (Internal)
### HealthUseCase
- `Check(ctx context.Context) *dto.HealthResponse`
  - Returns: A comprehensive health report containing overall status and individual dependency states.
