# Router Split Pattern (Go/Gin)

**Extracted:** 2026-02-04
**Context:** Refactoring large Go Gin router setups to improve maintainability.

## Problem
In Go applications using Gin, the main `router.go` file often accumulates all route definitions, becoming a "god object" that is difficult to read, maintain, and merge (prone to conflicts).

## Solution
Split the monolithic router by extracting related route groups (e.g., Auth, User, Blog) into separate files within the `router` package.

1.  **Create Separate Files:** Create `[domain]_routes.go` files (e.g., `auth_routes.go`).
2.  **Define Registration Functions:** Create public functions like `Register[Domain]Routes` that accept:
    *   The `*gin.RouterGroup` or `*gin.Engine`.
    *   The specific handlers or a dependencies struct (like `Params`).
    *   Necessary middleware.
3.  **Update Main Router:** Replace inline code in `router.go` with calls to these registration functions.

## Example

**Before (`router.go`):**
```go
func New(p Params) *gin.Engine {
    r := gin.New()
    v1 := r.Group("/api/v1")
    
    // Auth
    auth := v1.Group("/auth")
    auth.POST("/login", p.AuthHandler.Login)
    
    // Users
    users := v1.Group("/users")
    users.GET("/:id", p.UserHandler.Get)
    
    return r
}
```

**After:**

`auth_routes.go`:
```go
func RegisterAuthRoutes(rg *gin.RouterGroup, p Params) {
    g := rg.Group("/auth")
    g.POST("/login", p.AuthHandler.Login)
}
```

`router.go`:
```go
func New(p Params) *gin.Engine {
    r := gin.New()
    v1 := r.Group("/api/v1")
    
    RegisterAuthRoutes(v1, p)
    RegisterUserRoutes(v1, p)
    
    return r
}
```

## When to Use
*   When `router.go` exceeds ~200 lines or contains mixed concerns.
*   When multiple developers frequently touch `router.go` causing merge conflicts.
*   To improve code organization by grouping routes by domain/feature.
