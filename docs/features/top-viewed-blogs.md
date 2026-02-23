# Feature: top viewed blogs API

This document describes the top viewed blogs feature for the API. It covers the
endpoint behavior and high-level implementation impact so the feature is ready
for implementation.

## Objective

Help users discover standout content quickly by exposing a ranked list of blogs
with the highest all-time view count.

## Requirements

- [ ] Add `GET /blogs/top-viewed`.
- [ ] Return only blogs with `status=published` and `visibility=public`.
- [ ] Sort results by `view_count` in descending order.
- [ ] Support pagination with defaults `page=1` and `pageSize=10`.
- [ ] Enforce `pageSize` maximum of `50`.
- [ ] Reuse existing blog list response shape where practical (id, author, title,
      slug, excerpt, thumbnail, publishedAt, reaction counts).
- [ ] Add view tracking data source for blogs:
  - Add `view_count` to `blogs` with default `0`.
  - Increment `view_count` when a blog detail view is valid and accessible.
- [ ] Add Swagger annotations and tests for handler, use case, and repository.

## Non-goals

- Time-window ranking (for example, 7-day or 30-day trending).
- Personalized ranking by viewer.
- Fraud-resistant analytics (bot filtering, de-duplication by device/IP).
- Changes to the existing related/recommended blog logic.

## Acceptance criteria

1. `GET /blogs/top-viewed` returns `200` and includes only `published` +
   `public` blogs.
2. Returned items are ordered by `view_count` descending.
3. With no query params, API returns at most `10` items.
4. If `pageSize` exceeds `50`, the API enforces the configured maximum based on
   the current project pagination convention.
5. When a valid blog detail endpoint is called, `view_count` increases by `1`
   and affects top-viewed ranking.
6. Swagger documentation includes the new endpoint and related params/response.

## Technical context

### Impacted areas

- `internal/interfaces/http/router/blog_routes.go`
- `internal/interfaces/http/handler/blog/blog_handler.go`
- `internal/application/usecase/blog/usecase.go`
- `internal/domain/service/blog_service.go`
- `internal/domain/repository/blog_repository.go`
- `internal/infrastructure/persistence/postgres/repository/blog_repository.go`
- `internal/application/dto/blog.go` (if response includes `viewCount`)
- `docs/` Swagger artifacts

### Data changes

- Add `view_count` column to `blogs` (integer, non-negative, default `0`).
- Add or review indexing strategy to support ordered top-viewed queries.
