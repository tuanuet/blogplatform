# Bookmark Feature Design

## Overview
The Bookmark feature allows users to save blog posts for later reading. This is a personal collection tool that helps users curate content they find valuable or want to return to.

## Architecture

### Components
- **UseCase**: `BookmarkUseCase` (Application Layer)
- **Repository**: `BookmarkRepository` (Infrastructure Layer)

## Core Logic & Features

### 1. Toggling Bookmarks
- **Add Bookmark**: Associates a specific blog post with the current user.
- **Remove Bookmark**: Removes the association.

### 2. Retrieval
- **User Library**: Fetches a paginated list of blog posts bookmarked by a specific user.
- **Filtering**: Supports standard blog list parameters (page, pageSize).

## Data Model

### Relationships
The implementation relies on a Many-to-Many relationship between `Users` and `Blogs`.

- **User**: The entity performing the action.
- **Blog**: The entity being bookmarked.
- **Storage**: Typically handled via a join table (e.g., `user_bookmarks` or handled implicitly within repository logic).

## API Reference

### Interface: `BookmarkUseCase`

```go
type BookmarkUseCase interface {
    // Add a bookmark
    BookmarkBlog(ctx context.Context, userID, blogID uuid.UUID) error

    // Remove a bookmark
    UnbookmarkBlog(ctx context.Context, userID, blogID uuid.UUID) error

    // List user's bookmarks
    GetUserBookmarks(ctx context.Context, userID uuid.UUID, params *dto.BlogFilterParams) (*repository.PaginatedResult[dto.BlogListResponse], error)
}
```
