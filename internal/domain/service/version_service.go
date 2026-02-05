package service

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"errors"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
)

var (
	ErrVersionNotFound = errors.New("blog version not found")
	ErrVersionMismatch = errors.New("version does not belong to blog")
)

type VersionService interface {
	// CreateVersion creates a new version from a blog
	CreateVersion(ctx context.Context, blog *entity.Blog, editorID uuid.UUID, changeSummary string) (*entity.BlogVersion, error)

	// ListVersions lists versions for a blog
	ListVersions(ctx context.Context, blogID uuid.UUID, pagination repository.Pagination) (*repository.PaginatedResult[entity.BlogVersion], error)

	// GetVersion gets a specific version
	GetVersion(ctx context.Context, versionID uuid.UUID) (*entity.BlogVersion, error)

	// RestoreVersion restores a blog from a version (creates new version from old data)
	RestoreVersion(ctx context.Context, blogID uuid.UUID, versionID uuid.UUID, editorID uuid.UUID) (*entity.Blog, error)

	// DeleteVersion deletes a version
	DeleteVersion(ctx context.Context, versionID uuid.UUID, requesterID uuid.UUID) error
}
