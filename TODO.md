# Admin Dashboard Stats Implementation Plan

- [ ] **Step 1: Define DTOs**
  - Create `internal/application/dto/admin_stats.go`.
  - Define `MonthlyStat` struct (Month string, Count int64).
  - Define `DashboardStats` struct (Users, Blogs, Comments []MonthlyStat).

- [ ] **Step 2: Update Repository Interfaces**
  - Edit `internal/domain/repository/user_repository.go`: Add `GetRegistrationCountsByMonth(ctx)`.
  - Edit `internal/domain/repository/blog_repository.go`: Add `GetCreationCountsByMonth(ctx)`.
  - Edit `internal/domain/repository/comment_repository.go`: Add `GetCreationCountsByMonth(ctx)`.

- [ ] **Step 3: Implement Repositories (GORM)**
  - Edit `internal/infrastructure/persistence/postgres/repository/user_repository.go`: Implement aggregation query using `DATE_TRUNC`.
  - Edit `internal/infrastructure/persistence/postgres/repository/blog_repository.go`: Implement aggregation query.
  - Edit `internal/infrastructure/persistence/postgres/repository/comment_repository.go`: Implement aggregation query.

- [ ] **Step 4: Implement UseCase**
  - Create `internal/application/usecase/admin_usecase.go`.
  - Define `AdminUseCase` interface and struct.
  - Inject User, Blog, Comment repositories.
  - Implement `GetDashboardStats(ctx) (*dto.DashboardStats, error)`.

- [ ] **Step 5: Implement HTTP Handler**
  - Create `internal/interfaces/http/handler/admin_handler.go`.
  - Inject `AdminUseCase`.
  - Implement `GetStats(c *gin.Context)`.

- [ ] **Step 6: Configure Router**
  - Edit `internal/interfaces/http/router/router.go`.
  - Register `GET /admin/stats`.

- [ ] **Step 7: Verification**
  - Run integration test or curl check.
