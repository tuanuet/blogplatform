# Authentication System Implementation Plan

## Phase 1: Dependencies & Schema (High Priority)
- [ ] [High] Install Dependencies
    - Run `go get golang.org/x/oauth2`
    - Verify `golang.org/x/crypto` is available
- [ ] [High] Database Schema Migration
    - Create `migrations/000011_auth_schema.up.sql`
    - Create `migrations/000011_auth_schema.down.sql`
    - Add `social_accounts` table
    - Add `is_verified` and `verification_token` to `users` table
- [ ] [High] Domain Entities
    - Update `internal/domain/entity/user.go`: Add fields
    - Create `internal/domain/entity/social_account.go`: Define struct

## Phase 2: Repository Layer (High Priority)
- [ ] [High] Update User Repository
    - Update `internal/domain/repository/user_repository.go`: Add `Create` method
    - Update `internal/infrastructure/persistence/postgres/repository/user_repository.go`: Implement `Create`
- [ ] [High] Social Account Repository
    - Create `internal/domain/repository/social_account_repository.go` (Interface)
    - Create `internal/infrastructure/persistence/postgres/repository/social_account_repository.go` (Implementation)
- [ ] [High] Session Repository (Redis)
    - Create `internal/domain/repository/session_repository.go` (Interface)
    - Create `internal/infrastructure/cache/session_repository.go` (Implementation using Redis)

## Phase 3: Use Case Layer (High Priority)
- [ ] [High] Auth UseCase - Core
    - Create `internal/application/usecase/auth_usecase.go`
    - Implement `Register(ctx, email, password)` -> Hashing, User Creation
    - Implement `Login(ctx, email, password)` -> Verification, Session Creation
    - Implement `Logout(ctx, sessionID)`
- [ ] [Medium] Auth UseCase - Social & Verification
    - Implement `SocialLogin(ctx, provider, token)` -> Find/Create User, Link Account
    - Implement `VerifyEmail(ctx, token)`

## Phase 4: Interface Layer (Handlers & Middleware) (Medium Priority)
- [ ] [High] Auth Handler - Core
    - Create `internal/interfaces/http/dto/auth_dto.go`: Define Request/Response structs
    - Create `internal/interfaces/http/handler/auth/auth_handler.go`
    - Implement `Register` endpoint
    - Implement `Login` endpoint
    - Implement `Logout` endpoint
- [ ] [Medium] Auth Middleware
    - Create `internal/interfaces/http/middleware/auth_middleware.go`: Session validation
    - Create `internal/interfaces/http/middleware/rate_limit.go`: Redis-based limiting
- [ ] [Medium] Auth Handler - Social
    - Implement `SocialLogin` endpoint (OAuth redirect)
    - Implement `SocialCallback` endpoint
- [ ] [Low] Router Integration
    - Update `internal/interfaces/http/router/router.go`: Register Auth routes
    - Register Middleware

## Phase 5: Configuration & wiring (Low Priority)
- [ ] [Low] Configuration
    - Update `config.yaml`: Add OAuth providers (Google, GitHub)
    - Update `internal/infrastructure/config/config.go`: Load OAuth config
- [ ] [Low] Module Wiring
    - Create `cmd/api/modules/auth_module.go`: Wire Handler, UseCase, Repos
