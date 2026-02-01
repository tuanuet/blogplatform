package repository

import (
	"context"

	"github.com/google/uuid"
)

// Pagination holds pagination parameters
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// PaginatedResult holds paginated query results
type PaginatedResult[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
}

// BaseRepository defines the base repository interface using generics
type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id uuid.UUID) (*T, error)
	FindAll(ctx context.Context, pagination Pagination) (*PaginatedResult[T], error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uuid.UUID) error
}
