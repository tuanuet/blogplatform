# Series Feature Design

## Overview
The Series feature allows authors to organize related blog posts into a cohesive collection or sequence. This is useful for multi-part tutorials, guides, or thematic collections.

## Architecture

### Components
- **Entity**: `Series` (Domain Layer)
- **UseCase**: `SeriesUseCase` (Application Layer)
- **Repository**: `SeriesRepository` (Infrastructure Layer)

## Core Logic & Features

### 1. Series Management
- **Create**: Authors can create a series with a title, slug, and description.
- **Update**: Authors can modify series metadata.
- **Delete**: Authors can delete a series (does not delete the blogs within it).

### 2. Content Curation
- **Add Blog**: Add an existing blog post to a series.
- **Remove Blog**: Remove a blog post from a series.
- **Ownership Check**: Enforces that only the author of the series can modify it or add blogs to it.

### 3. Discovery
- **Get By ID/Slug**: Retrieve series details.
- **List**: Search and filter series by author or keywords.

## Data Model

### Entity: `Series`
| Field | Type | Description |
|-------|------|-------------|
| `id` | UUID | Primary Key |
| `author_id` | UUID | FK to User |
| `title` | String | Series title |
| `slug` | String | URL identifier |
| `description` | Text | Series summary |
| `created_at` | Timestamp | Creation time |
| `updated_at` | Timestamp | Update time |

### Relationships
- **Belongs To**: `User` (Author).
- **Many-to-Many**: `Blogs` (via `series_blogs` join table).

## API Reference

### Interface: `SeriesUseCase`

```go
type SeriesUseCase interface {
    CreateSeries(ctx context.Context, userID uuid.UUID, req *dto.CreateSeriesRequest) (*dto.SeriesResponse, error)
    UpdateSeries(ctx context.Context, userID, seriesID uuid.UUID, req *dto.UpdateSeriesRequest) (*dto.SeriesResponse, error)
    DeleteSeries(ctx context.Context, userID, seriesID uuid.UUID) error
    GetSeriesByID(ctx context.Context, id uuid.UUID) (*dto.SeriesResponse, error)
    GetSeriesBySlug(ctx context.Context, slug string) (*dto.SeriesResponse, error)
    ListSeries(ctx context.Context, params *dto.SeriesFilterParams) ([]dto.SeriesResponse, int64, error)
    AddBlogToSeries(ctx context.Context, userID, seriesID, blogID uuid.UUID) error
    RemoveBlogFromSeries(ctx context.Context, userID, seriesID, blogID uuid.UUID) error
}
```
