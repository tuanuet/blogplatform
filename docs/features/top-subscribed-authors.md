# Feature: top subscribed authors API

This document describes the top subscribed authors feature for the API. It
covers endpoint behavior and high-level implementation impact so the feature is
ready for implementation.

## Objective

Expose a public endpoint that returns authors ranked by all-time subscriber
count so clients can build discovery surfaces and improve subscription
conversion.

## Requirements

- [ ] Add `GET /api/v1/authors/top-subscribed`.
- [ ] Make the endpoint public (no authentication required) and apply existing
      public middleware and rate limiting.
- [ ] Support pagination with defaults `page=1`, `pageSize=20`.
- [ ] Enforce `pageSize` maximum of `100`.
- [ ] Rank by `subscriber_count` descending, with deterministic tie-breaker by
      `author_id` ascending.
- [ ] Exclude soft-deleted or disabled authors.
- [ ] Return summary-only author fields (`authorId`, `username`, `displayName`,
      `avatarUrl`) plus `subscriberCount`.
- [ ] Reuse existing validation and error envelope conventions (`400` for
      invalid query params, `500` for internal failures).
- [ ] Add Swagger annotations and tests for handler, use case, and repository.
- [ ] Ship database-first in v1 (no Redis cache in scope).

## Non-goals

- Time-window ranking (for example, 7-day or 30-day filters).
- Personalized ranking or recommendation logic.
- Real-time ranking updates via streams/websockets.
- Subscription business-rule changes (tiers, billing, eligibility changes).
- Admin ranking management UI.

## Acceptance criteria

1. `GET /api/v1/authors/top-subscribed?page=1&pageSize=20` returns `200` with
   authors sorted by all-time `subscriberCount` descending.
2. Each returned item includes `authorId`, `username`, `displayName`,
   `avatarUrl`, and `subscriberCount`.
3. Response includes pagination metadata (`page`, `pageSize`, `totalItems`,
   `totalPages`).
4. Invalid `page` or `pageSize` returns `400` using the standard error shape.
5. Soft-deleted or disabled authors are not returned.
6. Existing subscribe/unsubscribe endpoints continue to behave as before.
7. Swagger documentation includes the new endpoint, query params, and response.

## Technical context

### Impacted areas

- `internal/interfaces/http/router/subscription_routes.go`
- `internal/interfaces/http/handler/subscription/*`
- `internal/application/usecase/subscription/*`
- `internal/domain/repository/subscription_repository.go`
- `internal/infrastructure/persistence/postgres/repository/subscription_repository.go`
- `docs/` Swagger artifacts

### Data changes

- No mandatory schema change expected in v1.
- Optional optimization after v1: add or verify index support for aggregation
  on subscription author keys used by ranking queries.
