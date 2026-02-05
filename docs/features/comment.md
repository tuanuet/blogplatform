# Comment Feature Design

## Overview
The Comment feature enables social interaction on blog posts. Users can post comments, reply to existing comments (threaded discussions), and manage their contributions.

## Architecture

### Components
- **Entity**: `Comment` (Domain Layer)
- **UseCase**: `CommentUseCase` (Application Layer)
- **Service**: `CommentService` (Domain Service Layer)
- **Repository**: `CommentRepository` (Infrastructure Layer)

## Core Logic & Features

### 1. Posting Comments
- **Top-Level Comments**: Users can post comments directly on a blog post.
- **Replies**: Users can reply to specific comments, creating a parent-child relationship.

### 2. Retrieval
- **Paginated Lists**: Comments are retrieved by Blog ID with pagination support.
- **Hierarchy**: The system supports nested comments (replies are typically loaded with the parent or lazily).

### 3. Moderation & Management
- **Update**: Users can edit the content of their own comments.
- **Delete**: Users can delete their own comments.
- **Access Control**: Validates that the user modifying the comment is the author.

## Data Model

### Entity: `Comment`
| Field | Type | Description |
|-------|------|-------------|
| `id` | UUID | Primary Key |
| `blog_id` | UUID | FK to Blog |
| `user_id` | UUID | FK to User (Author) |
| `parent_id` | UUID | FK to Comment (Nullable, for replies) |
| `content` | Text | The comment text |
| `created_at` | Timestamp | Creation time |
| `updated_at` | Timestamp | Last update time |
| `deleted_at` | Timestamp | Soft delete timestamp |

### Relationships
- **Belongs To**: `Blog`, `User`.
- **Self-Referencing**: `Parent` (Comment) and `Replies` (List of Comments).

## API Reference

### Interface: `CommentUseCase`

```go
type CommentUseCase interface {
    Create(ctx context.Context, userID, blogID uuid.UUID, req *dto.CreateCommentRequest) (*dto.CommentResponse, error)
    GetByID(ctx context.Context, id uuid.UUID) (*dto.CommentResponse, error)
    GetByBlogID(ctx context.Context, blogID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[dto.CommentResponse], error)
    Update(ctx context.Context, id, userID uuid.UUID, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error)
    Delete(ctx context.Context, id, userID uuid.UUID) error
}
```
