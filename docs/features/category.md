# Category Feature Design

## Overview
The Category feature provides a hierarchical organization system for blog posts. Each category represents a broad topic under which multiple blogs can be classified. This system helps users navigate content and discover relevant articles.

## Architecture

### Components
- **Entity**: `Category` (Domain Layer)
- **UseCase**: `CategoryUseCase` (Application Layer)
- **Service**: `CategoryService` (Domain Service Layer)
- **Repository**: `CategoryRepository` (Infrastructure Layer)

### Data Flow
1.  **API Layer** receives HTTP requests (Create, Read, Update, Delete).
2.  **UseCase Layer** validates request DTOs and delegates business logic to the Domain Service.
3.  **Domain Service** enforces domain rules (e.g., uniqueness of slugs, existence checks).
4.  **Repository Layer** handles database persistence.

## Core Logic & Features

### 1. Management (CRUD)
- **Creation**: Admins can create new categories with a Name, Slug, and optional Description.
- **Retrieval**: Categories can be retrieved by ID or listed with pagination.
- **Updates**: Name, Slug, and Description can be modified.
- **Deletion**: Categories can be removed. *Note: Impact on existing blogs should be considered (e.g., set null or cascading delete).*

### 2. Validation
- **Slug Uniqueness**: Ensures two categories cannot share the same URL-friendly identifier.
- **Name Uniqueness**: Prevents duplicate category names.

## Data Model

### Entity: `Category`
| Field | Type | Description |
|-------|------|-------------|
| `id` | UUID | Primary Key |
| `name` | String | Display name (Unique) |
| `slug` | String | URL-friendly identifier (Unique) |
| `description` | Text | Optional description |
| `created_at` | Timestamp | Creation time |
| `updated_at` | Timestamp | Last update time |
| `deleted_at` | Timestamp | Soft delete timestamp (optional) |

### Relationships
- **One-to-Many**: A `Category` can have multiple `Blogs`.

## API Reference

### Interface: `CategoryUseCase`

```go
type CategoryUseCase interface {
    Create(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
    GetByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error)
    List(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[dto.CategoryResponse], error)
    Update(ctx context.Context, id uuid.UUID, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```
