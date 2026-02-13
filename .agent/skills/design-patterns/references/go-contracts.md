# Go Contracts (Entity + Repository Interface)

Use this as a contract example in a Clean Architecture Go backend.

```go
package domain

import (
    "context"
    "time"
)

type Example struct {
    ID        string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type ExampleQuery struct {
    Limit  int
    Offset int
}

type ExampleRepository interface {
    Create(ctx context.Context, e *Example) error
    GetByID(ctx context.Context, id string) (*Example, error)
    List(ctx context.Context, q ExampleQuery) ([]Example, error)
}
```

Notes:

- Keep domain types free of Gin/GORM/Redis.
- Accept `context.Context` as the first parameter.
- Keep queries bounded (limit/pagination decided).
