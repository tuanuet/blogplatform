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
