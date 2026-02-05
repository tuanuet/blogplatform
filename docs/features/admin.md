# Admin Feature

## Overview
The Admin feature provides administrative capabilities for system oversight. Currently, its primary function is to aggregate and present dashboard statistics, giving administrators a high-level view of platform activity and growth.

## Architecture
The Admin feature follows the Clean Architecture pattern:
- **UseCase**: `AdminUseCase` orchestrates the retrieval and aggregation of data from various domain repositories.
- **Repositories**: Direct dependencies on `UserRepository`, `BlogRepository`, and `CommentRepository` to fetch statistical data.

## Core Logic & Features
### Dashboard Statistics
The `GetDashboardStats` operation provides a 12-month retrospective of system activity.
- **Data Aggregation**: Fetches monthly counts for:
  - New User registrations
  - New Blog posts created
  - New Comments posted
- **Timeline Construction**: Generates a continuous 12-month timeline, ensuring months with zero activity are correctly represented.
- **Chronological Ordering**: Returns data sorted chronologically (Oldest to Newest) to facilitate chart rendering on the frontend.

## Data Model
The feature relies on Data Transfer Objects (DTOs) for its response structure:

### MonthlyStat
Represents activity for a single month.
```go
type MonthlyStat struct {
    Month       string // Format: "YYYY-MM"
    NewUsers    int64
    NewBlogs    int64
    NewComments int64
}
```

### DashboardStatsResponse
The wrapper for the statistics list.
```go
type DashboardStatsResponse struct {
    Stats []MonthlyStat
}
```

## API Reference (Internal)
### AdminUseCase
- `GetDashboardStats(ctx context.Context) (*dto.DashboardStatsResponse, error)`
  - Returns: A 12-month statistical history of platform activity.
