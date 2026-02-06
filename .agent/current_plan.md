# Implementation Plan: Highlighted Series API

## Tasks

- [ ] **Task 1: DTOs & Domain Interfaces** <!-- id: 1 -->
  - [ ] Add `HighlightedSeriesResponse` struct to `internal/application/dto/series.go`
  - [ ] Add `HighlightedSeriesResult` struct to `internal/domain/repository/series_repository.go`
  - [ ] Add `GetHighlighted(ctx context.Context, limit int) ([]HighlightedSeriesResult, error)` to `SeriesRepository` interface in `internal/domain/repository/series_repository.go`
  - [ ] Run `go generate ./...` to update mocks
  - [ ] Verify: `go build ./...`

- [ ] **Task 2: Repository Implementation** <!-- id: 2 -->
  - [ ] Implement `GetHighlighted` in `internal/infrastructure/persistence/postgres/repository/series_repository.go` using raw SQL
  - [ ] Create `internal/infrastructure/persistence/postgres/repository/series_repository_test.go`
  - [ ] Add test `TestSeriesRepository_GetHighlighted` to verify raw SQL execution and mapping
  - [ ] Verify: `go test ./internal/infrastructure/persistence/postgres/repository/...`

- [ ] **Task 3: UseCase Implementation** <!-- id: 3 -->
  - [ ] Add `GetHighlightedSeries(ctx context.Context) ([]*dto.HighlightedSeriesResponse, error)` to `SeriesUseCase` interface in `internal/application/usecase/series/usecase.go`
  - [ ] Implement `GetHighlightedSeries` in `internal/application/usecase/series/usecase.go` (call repo, map Result to DTO)
  - [ ] Create/Update `internal/application/usecase/series/usecase_test.go`
  - [ ] Add test `TestSeriesUseCase_GetHighlightedSeries`
  - [ ] Verify: `go test ./internal/application/usecase/series/...`

- [ ] **Task 4: HTTP Handler** <!-- id: 4 -->
  - [ ] Implement `GetHighlightedSeries(c *gin.Context)` in `internal/interfaces/http/handler/series/series_handler.go`
  - [ ] Create/Update `internal/interfaces/http/handler/series/series_handler_test.go`
  - [ ] Add test `TestSeriesHandler_GetHighlightedSeries`
  - [ ] Verify: `go test ./internal/interfaces/http/handler/series/...`

- [ ] **Task 5: Router & Integration** <!-- id: 5 -->
  - [ ] Register `GET /series/highlighted` in `internal/interfaces/http/router/series_routes.go`
  - [ ] Create `tests/integration/highlighted_series_test.go`
  - [ ] Add E2E test `TestGetHighlightedSeries_E2E` (setup data, call API, verify JSON)
  - [ ] Verify: `go test ./tests/integration/...`

Approve to start building?
