# Uber FX Dependency Injection - Rollback Plan

## Feature Overview

Integrated Uber FX for dependency injection, replacing manual wiring in `main.go`.

## Changes Made

1. Added `go.uber.org/fx` dependency
2. Created FX modules in `cmd/api/modules/`:
   - `config_module.go` - Configuration providers
   - `database_module.go` - Database + Redis with lifecycle hooks
   - `repository_module.go` - All repository implementations
   - `service_module.go` - All service implementations
   - `handler_module.go` - All HTTP handlers
   - `http_module.go` - Router + HTTP server with lifecycle

3. Refactored `cmd/api/main.go` to use `fx.New()` composition

## Benefits

- ✅ Cleaner dependency injection
- ✅ Automatic lifecycle management (OnStart/OnStop hooks)
- ✅ Graceful shutdown handled by FX
- ✅ Better testability with FX validation
- ✅ Modular architecture

## Rollback Instructions

If issues occur, rollback with these steps:

### Step 1: Restore Original main.go

Copy the backup main.go and replace the current one:

```bash
git checkout HEAD~1 -- cmd/api/main.go
```

Or use this command to see the original:

```bash
git show HEAD~1:cmd/api/main.go > cmd/api/main.go.bak
```

### Step 2: Remove FX Modules (optional)

```bash
rm -rf cmd/api/modules/
```

### Step 3: Remove FX Dependency (optional)

```bash
go mod edit -droprequire go.uber.org/fx
go mod edit -droprequire go.uber.org/dig
go mod tidy
```

### Step 4: Rebuild

```bash
go build ./cmd/api
```

## Verification Checklist

- [ ] Build succeeds: `go build ./...`
- [ ] All tests pass: `go test ./... -v`
- [ ] Health check works: `curl http://localhost:8080/api/v1/health`
- [ ] Graceful shutdown works: Send SIGTERM and verify clean shutdown
- [ ] All endpoints functional

## Pre-Roll Main.go (Backup Reference)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aiagent/boilerplate/internal/application/service"
	"github.com/aiagent/boilerplate/internal/infrastructure/cache"
	"github.com/aiagent/boilerplate/internal/infrastructure/config"
	"github.com/aiagent/boilerplate/internal/infrastructure/persistence/postgres"
	pgRepo "github.com/aiagent/boilerplate/internal/infrastructure/persistence/postgres/repository"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler"
	"github.com/aiagent/boilerplate/internal/interfaces/http/router"
	"github.com/aiagent/boilerplate/pkg/logger"
	"github.com/aiagent/boilerplate/pkg/validator"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.New(&cfg.Logger)
	logger.Info("Starting application...")
	validator.Init()

	db, err := postgres.NewDatabase(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", err, nil)
	}
	defer postgres.Close(db)

	redisClient, err := cache.NewRedisClient(&cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", err, nil)
	}
	defer redisClient.Close()

	// Repositories
	blogRepo := pgRepo.NewBlogRepository(db)
	categoryRepo := pgRepo.NewCategoryRepository(db)
	tagRepo := pgRepo.NewTagRepository(db)
	commentRepo := pgRepo.NewCommentRepository(db)
	subscriptionRepo := pgRepo.NewSubscriptionRepository(db)

	// Services
	healthService := service.NewHealthService(db, redisClient)
	blogService := service.NewBlogService(blogRepo, subscriptionRepo, tagRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)
	commentService := service.NewCommentService(commentRepo, blogRepo, subscriptionRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)

	// Handlers
	healthHandler := handler.NewHealthHandler(healthService)
	blogHandler := handler.NewBlogHandler(blogService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	tagHandler := handler.NewTagHandler(tagService)
	commentHandler := handler.NewCommentHandler(commentService)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	// Router
	r := router.NewRouter(
		healthHandler, blogHandler, categoryHandler,
		tagHandler, commentHandler, subscriptionHandler,
		cfg.Server.Mode,
	)
	engine := r.Setup()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("Server starting", map[string]interface{}{
			"port": cfg.Server.Port,
			"mode": cfg.Server.Mode,
		})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", err, nil)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", err, nil)
	}
	logger.Info("Server exited properly")
}
```

---

# Follow Feature - Rollback Plan

## Feature Overview

Implemented Follow/Following Model (Twitter-style social relationships) for the blog platform.

## Changes Made

### 1. Database Migration
- **File**: `migrations/000005_follow_feature.up.sql`
- **Table**: `follows` with unique constraint on (follower_id, following_id)
- **Indexes**: follower_id, following_id, created_at
- **Constraints**: Foreign keys to users table, prevents self-follows

### 2. Domain Layer
- **Entity**: `internal/domain/entity/follow.go`
- **Repository Interface**: `internal/domain/repository/follow_repository.go`
- **Service**: `internal/domain/service/follow_service.go`

### 3. Infrastructure Layer
- **Repository Implementation**: `internal/infrastructure/persistence/postgres/repository/follow_repository.go`

### 4. Application Layer
- **DTOs**: `internal/application/dto/follow.go`
- **UseCase**: `internal/application/usecase/follow_usecase.go`

### 5. Interface Layer
- **Handler**: `internal/interfaces/http/handler/follow_handler.go`
- **Routes**: Added to `internal/interfaces/http/router/router.go`

### 6. Dependency Injection
- Updated: `cmd/api/modules/domain_service_module.go`
- Updated: `cmd/api/modules/repository_module.go`
- Updated: `cmd/api/modules/usecase_module.go`
- Updated: `cmd/api/modules/handler_module.go`

## API Endpoints Added

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/users/:userId/follow` | Follow a user | Yes |
| DELETE | `/api/v1/users/:userId/follow` | Unfollow a user | Yes |
| GET | `/api/v1/users/:userId/followers` | Get user's followers | No |
| GET | `/api/v1/users/:userId/following` | Get users being followed | No |
| GET | `/api/v1/users/:userId/follow-counts` | Get follower/following counts | No |
| GET | `/api/v1/users/:userId/follow-status` | Check if current user follows target | Yes |

## Rollback Instructions

### Step 1: Revert Database Migration

```bash
# Run down migration to drop follows table
psql -U your_db_user -d your_db_name -f migrations/000005_follow_feature.down.sql
```

Or manually:
```sql
DROP TABLE IF EXISTS follows;
```

### Step 2: Remove Source Files

```bash
# Remove domain layer
rm internal/domain/entity/follow.go
rm internal/domain/repository/follow_repository.go
rm internal/domain/service/follow_service.go

# Remove infrastructure layer
rm internal/infrastructure/persistence/postgres/repository/follow_repository.go

# Remove application layer
rm internal/application/dto/follow.go
rm internal/application/usecase/follow_usecase.go

# Remove interface layer
rm internal/interfaces/http/handler/follow_handler.go

# Remove migration files
rm migrations/000005_follow_feature.up.sql
rm migrations/000005_follow_feature.down.sql
```

### Step 3: Revert Dependency Injection

Edit `cmd/api/modules/domain_service_module.go`:
- Remove: `service.NewFollowService,`

Edit `cmd/api/modules/repository_module.go`:
- Remove: `pgRepo.NewFollowRepository,`

Edit `cmd/api/modules/usecase_module.go`:
- Remove: `usecase.NewFollowUseCase,`

Edit `cmd/api/modules/handler_module.go`:
- Remove: `handler.NewFollowHandler,`

### Step 4: Revert Router

Edit `internal/interfaces/http/router/router.go`:
1. Remove `FollowHandler *handler.FollowHandler` from Params struct
2. Remove follow routes from v1 group:
   ```go
   v1.GET("/users/:userId/followers", ...)
   v1.GET("/users/:userId/following", ...)
   v1.GET("/users/:userId/follow-counts", ...)
   v1.GET("/users/:userId/follow-status", ...)
   v1.POST("/users/:userId/follow", ...)
   v1.DELETE("/users/:userId/follow", ...)
   ```

### Step 5: Rebuild

```bash
go mod tidy
go build ./cmd/api
```

## Verification Checklist

- [ ] Build succeeds: `go build ./...`
- [ ] Database migration reverted: `follows` table no longer exists
- [ ] No follow-related API endpoints respond
- [ ] All other features still work (subscriptions, blogs, etc.)
- [ ] Tests pass: `go test ./... -v`

## Quick Rollback Script

```bash
#!/bin/bash
# rollback_follow.sh

echo "Rolling back Follow Feature..."

# Database rollback
psql -U $DB_USER -d $DB_NAME -c "DROP TABLE IF EXISTS follows;"

# Remove files
rm -f internal/domain/entity/follow.go
rm -f internal/domain/repository/follow_repository.go
rm -f internal/domain/service/follow_service.go
rm -f internal/infrastructure/persistence/postgres/repository/follow_repository.go
rm -f internal/application/dto/follow.go
rm -f internal/application/usecase/follow_usecase.go
rm -f internal/interfaces/http/handler/follow_handler.go
rm -f migrations/000005_follow_feature.*.sql

echo "Follow Feature rolled back successfully!"
echo "Remember to revert changes in:"
echo "  - cmd/api/modules/*.go"
echo "  - internal/interfaces/http/router/router.go"
```
