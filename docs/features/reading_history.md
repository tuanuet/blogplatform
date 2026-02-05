# Reading History Feature Design

## Overview
The Reading History feature tracks which blog posts a user has read. It supports "Resume Reading" functionality (implicitly) and allows users to view their recently read content.

## Architecture

### Components
- **Entity**: `UserReadingHistory` (Domain Layer)
- **UseCase**: `ReadingHistoryUseCase` (Application Layer)
- **Repository**: `ReadingHistoryRepository` (Infrastructure Layer)

## Core Logic & Features

### 1. Tracking
- **Mark as Read**: Updates the timestamp (`LastReadAt`) when a user reads a blog. This is an "Upsert" operation—if the record exists, it updates the time; otherwise, it creates it.

### 2. Retrieval
- **Recent History**: Retrieves the most recently read blogs for a user, sorted by `LastReadAt` descending.

## Data Model

### Entity: `UserReadingHistory`
This entity uses a composite primary key.

| Field | Type | Description |
|-------|------|-------------|
| `user_id` | UUID | PK, FK to User |
| `blog_id` | UUID | PK, FK to Blog |
| `last_read_at` | Timestamp | The last time the user accessed the blog |

### Relationships
- **Belongs To**: `User`, `Blog`.

## API Reference

### Interface: `ReadingHistoryUseCase`

```go
type ReadingHistoryUseCase interface {
    // Record that a user read a blog
    MarkAsRead(ctx context.Context, userID, blogID uuid.UUID) error

    // Get recently read blogs
    GetHistory(ctx context.Context, userID uuid.UUID, limit int) (*dto.ReadingHistoryListResponse, error)
}
```
