# Tag Feature Design

## Overview
The Tag feature allows for granular classification of blog posts. Unlike categories which are broad, tags are specific descriptors that help users find content related to specific keywords or themes.

## Architecture

### Components
- **Entity**: `Tag` (Domain Layer)
- **UseCase**: `TagUseCase` (Application Layer)
- **Service**: `TagService` (Domain Service Layer)
- **Repository**: `TagRepository` (Infrastructure Layer)

## Core Logic & Features

### 1. Management (CRUD)
- **Creation**: Tags are created with a Name and automatically generated Slug.
- **Retrieval**: Retrieve tags by ID or list them with pagination.
- **Updates**: Modify tag name or slug.
- **Deletion**: Remove tags from the system.

### 2. Validation
- **Uniqueness**: Ensures tag names and slugs are unique across the system.

## Data Model

### Entity: `Tag`
| Field | Type | Description |
|-------|------|-------------|
| `id` | UUID | Primary Key |
| `name` | String | Display name (Unique) |
| `slug` | String | URL-friendly identifier (Unique) |
| `created_at` | Timestamp | Creation time |
| `updated_at` | Timestamp | Last update time |

### Relationships
- **Many-to-Many**: A `Tag` can be associated with multiple `Blogs` via the `blog_tags` join table.

## API Reference

### Interface: `TagUseCase`

```go
type TagUseCase interface {
    Create(ctx context.Context, req *dto.CreateTagRequest) (*dto.TagResponse, error)
    GetByID(ctx context.Context, id uuid.UUID) (*dto.TagResponse, error)
    List(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[dto.TagResponse], error)
    Update(ctx context.Context, id uuid.UUID, req *dto.UpdateTagRequest) (*dto.TagResponse, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```
