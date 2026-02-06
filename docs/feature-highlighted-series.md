# Feature Specification: Highlighted Series API

## Objective

Tạo API endpoint `GET /api/v1/series/highlighted` trả về top 10 series có lượng subscription cao nhất để hiển thị trên homepage.

---

## Problem Statement

Người dùng cần một cách để khám phá các series phổ biến trên platform. Hiện tại không có tính năng highlight series theo subscription count.

---

## Requirements

### Functional

- [ ] Tạo endpoint `GET /api/v1/series/highlighted`
- [ ] Query `UserSeriesPurchase`, đếm theo `SeriesID`, order DESC, limit 10
- [ ] Join với `Series` và `User` để lấy thông tin đầy đủ
- [ ] Trả về subscriber count cho mỗi series
- [ ] Trả về HTTP 200 với empty array nếu không có data
- [ ] Filter `DeletedAt IS NULL`

### Non-Functional

- [ ] **Performance**: Response time < 200ms
- [ ] **Scalability**: Có thể scale với hàng triệu subscriptions

---

## User Stories

### Story: Xem Series nổi bật

**As a** end user,  
**I want to** see top 10 most subscribed series on homepage,  
**So that** I can discover popular content to follow.

**Acceptance Criteria:**

- [ ] Given user visits homepage, When API returns, Then top 10 series with most UserSeriesPurchase count are displayed
- [ ] Series are sorted by subscriber count in descending order
- [ ] Each series item shows: ID, Title, Slug, Description, Author info, Subscriber count, Blog count, Author avatar
- [ ] If no series have subscriptions, return empty list (no error)
- [ ] Deleted series are not included in results

---

## Technical Context

### Impacted Services

| Layer | File Path |
|-------|-----------|
| **DTOs** | `/Users/mrt/workspaces/boilerplate/aiagent/internal/application/dto/series.go` |
| **Repository Interface** | `/Users/mrt/workspaces/boilerplate/aiagent/internal/domain/repository/series_repository.go` |
| **Repository Impl** | `/Users/mrt/workspaces/boilerplate/aiagent/internal/infrastructure/persistence/postgres/repository/series_repository.go` |
| **UseCase Interface** | `/Users/mrt/workspaces/boilerplate/aiagent/internal/application/usecase/series/usecase.go` |
| **Handler** | `/Users/mrt/workspaces/boilerplate/aiagent/internal/interfaces/http/handler/series/series_handler.go` |
| **Router** | `/Users/mrt/workspaces/boilerplate/aiagent/internal/interfaces/http/router/series_routes.go` |

### Dependencies

- `Series` entity (`/Users/mrt/workspaces/boilerplate/aiagent/internal/domain/entity/series.go`)
- `UserSeriesPurchase` entity (`/Users/mrt/workspaces/boilerplate/aiagent/internal/domain/entity/user_series_purchase.go`)
- `User` entity (cho author info)

### Database Relationships

```
Series 1----N Blog (many2many via series_blogs)
Series 1----N UserSeriesPurchase (via SeriesID)
User   1----N UserSeriesPurchase (as purchaser)
```

---

## Data Transfer Objects

### Response DTO

```go
type HighlightedSeriesResponse struct {
    ID              uuid.UUID  `json:"id"`
    Title           string     `json:"title"`
    Slug            string     `json:"slug"`
    Description     string     `json:"description"`
    AuthorID        uuid.UUID  `json:"authorId"`
    AuthorName      string     `json:"authorName"`
    AuthorAvatarURL *string    `json:"authorAvatarUrl,omitempty"`
    SubscriberCount int        `json:"subscriberCount"`
    BlogCount       int        `json:"blogCount"`
    CreatedAt       time.Time  `json:"createdAt"`
}
```

---

## API Specification

### Endpoint

```
GET /api/v1/series/highlighted
```

### Query Parameters

None (fixed top 10)

### Response

**Success (200 OK)**

```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Series Title",
      "slug": "series-title",
      "description": "Series description",
      "authorId": "550e8400-e29b-41d4-a716-446655440001",
      "authorName": "Author Name",
      "authorAvatarUrl": "https://cdn.example.com/avatars/author.jpg",
      "subscriberCount": 150,
      "blogCount": 5,
      "createdAt": "2026-01-01T00:00:00Z"
    }
  ]
}
```

**Empty Response (200 OK)**

```json
{
  "success": true,
  "data": []
}
```

**Error Response (500 Internal Server Error)**

```json
{
  "success": false,
  "message": "Internal server error"
}
```

---

## SQL Query Strategy

```sql
-- Pseudocode for repository implementation
SELECT
    s.id,
    s.title,
    s.slug,
    s.description,
    s.author_id,
    s.created_at,
    u.name as author_name,
    u.avatar_url as author_avatar_url,
    COALESCE(subscriber_counts.subscriber_count, 0) as subscriber_count,
    COALESCE(blog_counts.blog_count, 0) as blog_count
FROM series s
LEFT JOIN users u ON s.author_id = u.id
LEFT JOIN (
    SELECT series_id, COUNT(*) as subscriber_count
    FROM user_series_purchases
    GROUP BY series_id
) subscriber_counts ON s.id = subscriber_counts.series_id
LEFT JOIN (
    SELECT sb.series_id, COUNT(*) as blog_count
    FROM series_blogs sb
    JOIN blogs b ON sb.blog_id = b.id
    WHERE b.deleted_at IS NULL
    GROUP BY sb.series_id
) blog_counts ON s.id = blog_counts.series_id
WHERE s.deleted_at IS NULL
ORDER BY subscriber_count DESC
LIMIT 10;
```

---

## User Flows

### Happy Path

```
User → Homepage → GET /api/v1/series/highlighted → 
Backend (Query + Count + Join + Limit 10) → 
Frontend displays highlighted series section
```

### Edge Cases

| Edge Case | Handling |
|-----------|----------|
| No subscriptions | Return `[]`, HTTP 200 |
| < 10 series | Return all available |
| DB error | Return `500` with error |
| Series with deleted author | Skip hoặc set authorName = "Unknown" |
| Series is soft-deleted | Filter out with `DeletedAt IS NULL` |

---

## Out of Scope

- Pagination (fixed top 10)
- Filtering (no authorId, no search)
- Caching (có thể add sau nếu cần)
- Admin management (không có tính năng manually highlight)
- Customizable limit (always return top 10)

---

## Implementation Tasks

### 1. Add DTO

- Thêm `HighlightedSeriesResponse` trong `internal/application/dto/series.go`

### 2. Add Repository Method

- Interface: `GetHighlightedSeries(ctx context, limit int) ([]*Series, []int, error)` (returns series + subscriber counts)
- Implementation: Raw SQL query với joins và counting

### 3. Add UseCase Method

- `GetHighlightedSeries(ctx context) ([]*dto.HighlightedSeriesResponse, error)`

### 4. Add Handler Method

- `GetHighlightedSeries(c *gin.Context)`

### 5. Add Route

- `router.GET("/series/highlighted", h.GetHighlightedSeries)`

---

## Acceptance Checklist

- [ ] API endpoint trả về đúng 10 series với subscriber count cao nhất
- [ ] Series được sort theo subscriber count giảm dần
- [ ] Response chứa tất cả fields: id, title, slug, description, authorId, authorName, authorAvatarUrl, subscriberCount, blogCount, createdAt
- [ ] Deleted series không xuất hiện trong kết quả
- [ ] Empty array được trả về khi không có series nào
- [ ] API response time < 200ms
- [ ] Unit tests cho repository, usecase, và handler
- [ ] Integration test cho API endpoint

---

## Notes

- Subscription count đếm từ `UserSeriesPurchase` table
- Blog count đếm từ `series_blogs` join `blogs` (với `deleted_at IS NULL`)
- Author avatar có thể null nên dùng pointer và `omitempty` trong JSON
