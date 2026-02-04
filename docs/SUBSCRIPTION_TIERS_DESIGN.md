# Multi-Tier Subscription Plans - Design Summary

**Feature**: Multi-Tier Subscription Plans  
**Date**: 2026-02-03  
**Architect**: Claude (Architect Agent)  
**Status**: ✅ Design Approved - Ready for Planning  

---

## 1. Executive Summary

This design implements a 4-tier subscription system (FREE, BRONZE, SILVER, GOLD) allowing authors to monetize content through flexible, tag-based access control. The design follows Clean Architecture principles with clear separation between domain entities, repositories, services, and API layers.

### Key Design Decisions (User-Confirmed)

1. **Tags are globally shared** - tier assignments are per-author (same tag can have different tiers for different authors)
2. **Price validation** - Warn only, allow any prices (no hard enforcement of BRONZE < SILVER < GOLD)
3. **Tier storage** - VARCHAR(20) (not PostgreSQL ENUM for migration flexibility)
4. **Plan API** - Upsert operation (single POST endpoint handles create/update)
5. **Access check** - Public endpoint (works for both authenticated and anonymous users)
6. **Purchase API** - Reuse existing `POST /api/v1/payments` with `plan_id` parameter

---

## 2. Database Schema Design

### 2.1 New Tables

#### `subscription_plans`
Stores author-defined pricing plans for each tier.

**Columns:**
- `id` (UUID, PK) - Unique plan identifier
- `author_id` (UUID, FK → users.id) - Author who created the plan
- `tier` (VARCHAR(20), CHECK IN) - Tier level (FREE, BRONZE, SILVER, GOLD)
- `price` (DECIMAL(15,2), CHECK >= 0) - Price in VND
- `duration_days` (INT, DEFAULT 30) - Subscription duration
- `name` (VARCHAR(100), NULLABLE) - Optional custom tier name
- `description` (TEXT, NULLABLE) - Optional tier description
- `is_active` (BOOLEAN, DEFAULT true) - Active status
- `created_at`, `updated_at` (TIMESTAMP)

**Constraints:**
- UNIQUE(author_id, tier) - One plan per tier per author
- CHECK tier IN ('FREE', 'BRONZE', 'SILVER', 'GOLD')
- CHECK price >= 0
- CHECK duration_days > 0

**Indexes:**
- `idx_subscription_plans_author_tier` ON (author_id, tier) - Composite for fast lookups
- `idx_subscription_plans_author_id` ON (author_id)
- `idx_subscription_plans_is_active` ON (is_active)

---

#### `tag_tier_mappings`
Maps tags to required subscription tiers per author.

**Columns:**
- `id` (UUID, PK)
- `author_id` (UUID, FK → users.id)
- `tag_id` (UUID, FK → tags.id)
- `required_tier` (VARCHAR(20), CHECK IN)
- `created_at`, `updated_at` (TIMESTAMP)

**Constraints:**
- UNIQUE(author_id, tag_id) - One tier requirement per tag per author
- CHECK required_tier IN ('FREE', 'BRONZE', 'SILVER', 'GOLD')

**Indexes:**
- `idx_tag_tier_mappings_author_tag` ON (author_id, tag_id) - Composite for fast lookups
- `idx_tag_tier_mappings_author_id` ON (author_id)
- `idx_tag_tier_mappings_tag_id` ON (tag_id) - For reverse lookups

---

### 2.2 Modified Tables

#### `subscriptions` (Already modified in migration 000013)
- ✅ `tier` column already exists (added in SePay migration)
- ✅ `expires_at` column already exists
- No additional changes needed

---

### 2.3 Migration Files

**Location:**
- `migrations/000014_create_subscription_tiers.up.sql` ✅ Created
- `migrations/000014_create_subscription_tiers.down.sql` ✅ Created

**Triggers:**
- Auto-update `updated_at` on both tables using existing `update_updated_at_column()` function

---

## 3. API Contract Design (OpenAPI 3.0)

### 3.1 Endpoints Overview

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/authors/me/plans` | Required | Upsert subscription plans |
| GET | `/api/v1/authors/{authorId}/plans` | Public | Get author's pricing table |
| POST | `/api/v1/authors/me/tags/{tagId}/tier` | Required | Assign tag to tier |
| DELETE | `/api/v1/authors/me/tags/{tagId}/tier` | Required | Unassign tag from tier |
| GET | `/api/v1/authors/me/tag-tiers` | Required | Get all tag-tier mappings |
| GET | `/api/v1/blogs/{blogId}/access` | Public | Check blog access |
| POST | `/api/v1/payments` | Required | **MODIFIED**: Add `plan_id` support |

---

### 3.2 Endpoint Details

#### 1. POST /api/v1/authors/me/plans (Upsert Plans)

**Request:**
```json
{
  "plans": [
    {
      "tier": "BRONZE",
      "price": 50000,
      "name": "Gói Học Viên",
      "description": "Access to beginner tutorials"
    }
  ]
}
```

**Response (200 OK):**
```json
{
  "plans": [
    {
      "id": "uuid",
      "tier": "BRONZE",
      "price": 50000,
      "durationDays": 30,
      "name": "Gói Học Viên",
      "isActive": true,
      "createdAt": "2024-02-03T10:00:00Z",
      "updatedAt": "2024-02-03T10:00:00Z"
    }
  ],
  "warnings": [
    "Price hierarchy warning: SILVER (120000) should be > BRONZE (50000)"
  ]
}
```

**Business Logic:**
- Validates tier values
- Checks price hierarchy (warns if violations)
- Uses UPSERT: creates new or updates existing plans
- Only authenticated author can manage own plans

---

#### 2. GET /api/v1/authors/{authorId}/plans (Public)

**Response (200 OK):**
```json
{
  "authorId": "uuid",
  "plans": [
    {
      "tier": "FREE",
      "price": 0,
      "durationDays": 0,
      "tagCount": 5,
      "tags": ["Public Tutorials", "News"]
    },
    {
      "tier": "BRONZE",
      "price": 50000,
      "durationDays": 30,
      "name": "Gói Học Viên",
      "tagCount": 3,
      "tags": ["Go Basics", "React Intro"]
    }
  ]
}
```

**Business Logic:**
- Public endpoint (no auth required)
- Returns all 4 tiers (FREE always included)
- Shows tag counts and sample tag names per tier
- Used for pricing page display

---

#### 3. POST /api/v1/authors/me/tags/{tagId}/tier

**Request:**
```json
{
  "requiredTier": "BRONZE"
}
```

**Response (200 OK):**
```json
{
  "tagId": "uuid",
  "tagName": "Go Advanced",
  "requiredTier": "BRONZE",
  "affectedBlogsCount": 12
}
```

**Business Logic:**
- Validates tag ownership (tag must exist and be used by author)
- Upserts tag-tier mapping
- Returns count of affected blogs

---

#### 4. DELETE /api/v1/authors/me/tags/{tagId}/tier

**Response (200 OK):**
```json
{
  "message": "Tag unassigned, content is now FREE",
  "affectedBlogsCount": 12
}
```

**Business Logic:**
- Removes tier requirement (content becomes FREE)
- Returns count of affected blogs

---

#### 5. GET /api/v1/authors/me/tag-tiers

**Response (200 OK):**
```json
{
  "mappings": [
    {
      "tagId": "uuid",
      "tagName": "Go Basics",
      "requiredTier": "BRONZE",
      "blogCount": 8
    }
  ]
}
```

---

#### 6. GET /api/v1/blogs/{blogId}/access (Public/Authenticated)

**Response (200 OK - Accessible):**
```json
{
  "accessible": true,
  "userTier": "SILVER",
  "requiredTier": "BRONZE",
  "reason": "Your SILVER subscription includes BRONZE content"
}
```

**Response (200 OK - Blocked):**
```json
{
  "accessible": false,
  "userTier": "BRONZE",
  "requiredTier": "SILVER",
  "reason": "Upgrade to SILVER to access this content",
  "upgradeOptions": [
    {
      "tier": "SILVER",
      "price": 120000,
      "durationDays": 30,
      "planId": "uuid"
    }
  ]
}
```

**Business Logic:**
- Public endpoint (works for anonymous users → userTier = FREE)
- Finds blog's highest required tier from its tags
- Compares user's tier level vs. required tier level
- Returns upgrade options if blocked

---

#### 7. POST /api/v1/payments (Modified)

**Existing Request (Enhanced):**
```json
{
  "userId": "uuid",
  "amount": 120000,
  "type": "SUBSCRIPTION",
  "gateway": "VIETQR",
  "targetId": "author-uuid",
  "planId": "plan-uuid"  // ← NEW: Required for subscription purchases
}
```

**No API changes** - `planId` field already exists in `CreatePaymentRequest` DTO.

---

## 4. Domain Entity Design

### 4.1 SubscriptionTier Enum

**Location:** `internal/domain/entity/subscription_plan.go`

```go
type SubscriptionTier string

const (
    TierFree   SubscriptionTier = "FREE"
    TierBronze SubscriptionTier = "BRONZE"
    TierSilver SubscriptionTier = "SILVER"
    TierGold   SubscriptionTier = "GOLD"
)

// Level returns numeric level for comparison
func (t SubscriptionTier) Level() int {
    switch t {
    case TierFree: return 0
    case TierBronze: return 1
    case TierSilver: return 2
    case TierGold: return 3
    default: return 0
    }
}

func (t SubscriptionTier) IsValid() bool
func (t SubscriptionTier) String() string
```

**Design Notes:**
- Enum with numeric level mapping for tier comparison
- `Level()` method enables simple comparisons: `userTier.Level() >= requiredTier.Level()`
- `IsValid()` for input validation

---

### 4.2 SubscriptionPlan Entity

**Location:** `internal/domain/entity/subscription_plan.go`

**Fields:**
- ID, AuthorID, Tier, Price, DurationDays
- Name, Description (optional)
- IsActive (for soft deactivation)
- CreatedAt, UpdatedAt

**Relationships:**
- BelongsTo User (Author)

---

### 4.3 TagTierMapping Entity

**Location:** `internal/domain/entity/tag_tier_mapping.go`

**Fields:**
- ID, AuthorID, TagID, RequiredTier
- CreatedAt, UpdatedAt

**Relationships:**
- BelongsTo User (Author)
- BelongsTo Tag

---

### 4.4 Subscription Entity (Modified)

**Location:** `internal/domain/entity/subscription.go`

**Changes:**
- ✅ `Tier` field already exists (added in SePay migration)
- No code changes needed

---

## 5. Repository Interface Design

### 5.1 SubscriptionPlanRepository

**Location:** `internal/domain/repository/subscription_plan_repository.go`

**Interface Methods:**
```go
Create(ctx, *entity.SubscriptionPlan) error
FindByID(ctx, uuid.UUID) (*entity.SubscriptionPlan, error)
FindByAuthorAndTier(ctx, authorID, tier) (*entity.SubscriptionPlan, error)
FindByAuthor(ctx, authorID) ([]entity.SubscriptionPlan, error)
FindActiveByAuthor(ctx, authorID) ([]entity.SubscriptionPlan, error)
Update(ctx, *entity.SubscriptionPlan) error
Upsert(ctx, *entity.SubscriptionPlan) error  // ON CONFLICT DO UPDATE
Delete(ctx, uuid.UUID) error
WithTx(tx interface{}) SubscriptionPlanRepository
```

**Design Notes:**
- `Upsert` uses PostgreSQL `ON CONFLICT (author_id, tier) DO UPDATE`
- `WithTx` for transaction support

---

### 5.2 TagTierMappingRepository

**Location:** `internal/domain/repository/tag_tier_mapping_repository.go`

**Interface Methods:**
```go
Create(ctx, *entity.TagTierMapping) error
FindByID(ctx, uuid.UUID) (*entity.TagTierMapping, error)
FindByAuthorAndTag(ctx, authorID, tagID) (*entity.TagTierMapping, error)
FindByAuthor(ctx, authorID) ([]entity.TagTierMapping, error)
FindByTagIDs(ctx, authorID, []tagID) ([]entity.TagTierMapping, error)
Update(ctx, *entity.TagTierMapping) error
Upsert(ctx, *entity.TagTierMapping) error  // ON CONFLICT DO UPDATE
Delete(ctx, authorID, tagID) error
DeleteByID(ctx, uuid.UUID) error
CountBlogsByTagAndAuthor(ctx, authorID, tagID) (int64, error)
WithTx(tx interface{}) TagTierMappingRepository
```

**Design Notes:**
- `CountBlogsByTagAndAuthor` for "affected blogs count" feature
- `FindByTagIDs` for bulk lookups when checking blog access

---

### 5.3 SubscriptionRepository (Modified)

**Location:** `internal/domain/repository/subscription_repository.go`

**Enhanced Methods:**
- ✅ `FindActiveSubscription(ctx, userID, authorID)` already exists
- ✅ `UpdateExpiry(ctx, userID, authorID, expiresAt, tier)` already exists
- No changes needed

---

## 6. Service Layer Design (Contracts Only)

### 6.1 PlanManagementService

**Location:** `internal/domain/service/plan_management_service.go`

**Interface Methods:**
```go
UpsertPlans(ctx, authorID, []CreatePlanDTO) ([]entity.SubscriptionPlan, []string, error)
GetAuthorPlans(ctx, authorID) ([]PlanWithTags, error)
DeactivatePlan(ctx, authorID, tier) error
ActivatePlan(ctx, authorID, tier) error
```

**Responsibilities:**
- Price hierarchy validation (returns warnings, not errors)
- Plan upsert orchestration
- Aggregating plans with tag data

---

### 6.2 TagTierService

**Location:** `internal/domain/service/tag_tier_service.go`

**Interface Methods:**
```go
AssignTagToTier(ctx, authorID, tagID, tier) (*entity.TagTierMapping, int64, error)
UnassignTagFromTier(ctx, authorID, tagID) (int64, error)
GetAuthorTagTiers(ctx, authorID) ([]TagTierWithCount, error)
GetRequiredTierForBlog(ctx, blogID) (entity.SubscriptionTier, error)
```

**Responsibilities:**
- Tag ownership validation
- Blog count calculation
- Determining highest required tier for a blog

---

### 6.3 ContentAccessService

**Location:** `internal/domain/service/content_access_service.go`

**Interface Methods:**
```go
CheckBlogAccess(ctx, blogID, *userID) (*AccessResult, error)
GetUserTier(ctx, userID, authorID) (entity.SubscriptionTier, error)
```

**Responsibilities:**
- Blog access determination logic
- User tier retrieval (handles anonymous users → FREE tier)
- Upgrade options generation

**Access Check Algorithm:**
```
1. Fetch blog with tags
2. For each tag, get author's tier requirement (from tag_tier_mappings)
3. Determine highest required tier
4. Get user's current subscription tier (FREE if none)
5. Compare: accessible = (userTier.Level() >= requiredTier.Level())
6. If blocked, generate upgrade options from higher tiers
```

---

## 7. DTO Design

**Location:** `internal/application/dto/subscription_plan.go` ✅ Created

### Request DTOs
- `CreatePlanRequest` - Single plan in upsert request
- `UpsertPlansRequest` - Array of plans
- `AssignTagTierRequest` - Tier assignment

### Response DTOs
- `PlanResponse` - Single plan
- `UpsertPlansResponse` - Plans + warnings
- `PlanWithTagsResponse` - Plan with tag info (for public pricing table)
- `GetAuthorPlansResponse` - Full pricing table
- `AssignTagTierResponse` - Tag assignment result
- `UnassignTagTierResponse` - Unassignment result
- `TagTierMappingResponse` - Single mapping
- `GetTagTiersResponse` - All mappings
- `CheckBlogAccessResponse` - Access check result
- `UpgradeOption` - Upgrade option details

---

## 8. Handler Interface Design (Contracts Only)

### PlanHandler Interface

**Location:** `internal/interfaces/http/handler/plan/handler.go` (to be created)

```go
type PlanHandler interface {
    // Plan Management
    UpsertPlans(c *gin.Context)         // POST /api/v1/authors/me/plans
    GetAuthorPlans(c *gin.Context)       // GET /api/v1/authors/{authorId}/plans
    
    // Tag-Tier Management
    AssignTagToTier(c *gin.Context)      // POST /api/v1/authors/me/tags/{tagId}/tier
    UnassignTagFromTier(c *gin.Context)  // DELETE /api/v1/authors/me/tags/{tagId}/tier
    GetAuthorTagTiers(c *gin.Context)    // GET /api/v1/authors/me/tag-tiers
    
    // Access Control
    CheckBlogAccess(c *gin.Context)      // GET /api/v1/blogs/{blogId}/access
}
```

**Dependencies:**
- PlanManagementService
- TagTierService
- ContentAccessService

---

## 9. Routing Design

**Location:** `internal/interfaces/http/router/router.go` (to be modified)

```go
// In router.go, add to v1 group:

authors := v1.Group("/authors")
{
    // Public
    authors.GET("/:authorId/plans", planHandler.GetAuthorPlans)
    
    // Authenticated
    authorsAuth := authors.Group("")
    authorsAuth.Use(sessionAuth)
    {
        // Plan management
        authorsAuth.POST("/me/plans", planHandler.UpsertPlans)
        
        // Tag-tier management
        authorsAuth.POST("/me/tags/:tagId/tier", planHandler.AssignTagToTier)
        authorsAuth.DELETE("/me/tags/:tagId/tier", planHandler.UnassignTagFromTier)
        authorsAuth.GET("/me/tag-tiers", planHandler.GetAuthorTagTiers)
    }
}

blogs := v1.Group("/blogs")
{
    // Public/Authenticated (optional auth)
    blogs.GET("/:blogId/access", planHandler.CheckBlogAccess)
}
```

---

## 10. Integration with Existing Payment Flow

### Modified: PaymentService.HandleSePayWebhook

**Location:** `internal/domain/service/payment_service.go` (to be modified)

**Pseudocode:**
```go
func (s *PaymentService) HandleSePayWebhook(ctx, dto) error {
    // ... existing transaction validation ...
    
    // NEW: Fetch plan to get tier information
    plan, err := s.planRepo.FindByID(ctx, transaction.PlanID)
    if err != nil {
        return err
    }
    
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Update transaction status
        transaction.Status = TransactionStatusSuccess
        s.txRepo.Update(ctx, transaction)
        
        switch transaction.Type {
        case TransactionTypeSubscription:
            // NEW: Update subscription with tier
            subscription, _ := s.subRepo.FindActiveSubscription(ctx, transaction.UserID, transaction.AuthorID)
            if subscription == nil {
                // Create new subscription
                subscription = &Subscription{
                    UserID: transaction.UserID,
                    AuthorID: transaction.AuthorID,
                    Tier: plan.Tier,  // ← NEW
                    ExpiresAt: time.Now().Add(plan.DurationDays * 24h),
                }
                return s.subRepo.Create(ctx, subscription)
            } else {
                // Upgrade: Replace tier and reset expiry
                subscription.Tier = plan.Tier  // ← NEW
                subscription.ExpiresAt = time.Now().Add(plan.DurationDays * 24h)
                return s.subRepo.Update(ctx, subscription)
            }
        
        case TransactionTypeSeries:
            // ... existing series purchase logic ...
        }
        
        return nil
    })
}
```

**Changes Required:**
1. Add `planRepo SubscriptionPlanRepository` dependency
2. Fetch plan by `transaction.PlanID`
3. Set `subscription.Tier = plan.Tier` on create/update

---

## 11. Dependency Injection (Fx)

**Modules to Update:**

### 11.1 Repository Module
**Location:** `cmd/api/modules/repository.go`

```go
// Add to Provides:
fx.Provide(
    postgresrepo.NewSubscriptionPlanRepository,
    postgresrepo.NewTagTierMappingRepository,
)
```

---

### 11.2 Service Module
**Location:** `cmd/api/modules/service.go`

```go
// Add to Provides:
fx.Provide(
    service.NewPlanManagementService,
    service.NewTagTierService,
    service.NewContentAccessService,
)
```

---

### 11.3 Handler Module
**Location:** `cmd/api/modules/handler.go`

```go
// Add to Provides:
fx.Provide(
    planhandler.NewPlanHandler,
)
```

---

### 11.4 Router Module
**Location:** `cmd/api/modules/router.go`

```go
// Update Params struct:
type RouterParams struct {
    fx.In
    // ... existing handlers ...
    PlanHandler plan.PlanHandler  // ← NEW
}
```

---

## 12. Design Patterns Applied

### 12.1 SOLID Principles

**Single Responsibility:**
- `PlanManagementService` - Plan CRUD and validation
- `TagTierService` - Tag-tier mapping logic
- `ContentAccessService` - Access control logic

**Dependency Inversion:**
- All services depend on repository interfaces, not concrete implementations
- Handlers depend on service interfaces

**Interface Segregation:**
- Small, focused interfaces (no fat interfaces)

---

### 12.2 Repository Pattern
- Data access logic isolated in repositories
- Service layer never touches database directly

---

### 12.3 Service Pattern
- Business logic lives in services
- Handlers are thin: validate input → call service → return response

---

## 13. Testing Strategy (For Planner Reference)

### 13.1 Unit Tests (with mocks)

**Entity Tests:**
- `SubscriptionTier.Level()` returns correct hierarchy
- `SubscriptionTier.IsValid()` validates tier values

**Service Tests:**
- `PlanManagementService.UpsertPlans` validates price hierarchy
- `TagTierService.GetRequiredTierForBlog` returns highest tier
- `ContentAccessService.CheckBlogAccess` correctly determines access
- `PaymentService.HandleSePayWebhook` updates subscription tier

---

### 13.2 Integration Tests

**Workflow Tests:**
```go
func TestSubscriptionTierWorkflow(t *testing.T) {
    // 1. Author creates BRONZE plan
    // 2. Reader purchases BRONZE subscription
    // 3. Verify: Subscription.Tier = BRONZE, ExpiresAt set
    // 4. Test: Access BRONZE blog → allowed
    // 5. Test: Access SILVER blog → blocked
    // 6. Reader upgrades to SILVER
    // 7. Verify: Tier = SILVER, ExpiresAt reset
    // 8. Test: Access SILVER blog → allowed
}

func TestTagTierMapping(t *testing.T) {
    // 1. Author assigns "Go Basics" → BRONZE
    // 2. Blog has tags ["Go Basics", "Public"]
    // 3. Verify: CheckBlogAccess returns RequiredTier=BRONZE
    // 4. Unassign "Go Basics"
    // 5. Verify: RequiredTier=FREE
}
```

---

## 14. Files Created/Modified

### ✅ Created Files

**Database:**
- `migrations/000014_create_subscription_tiers.up.sql`
- `migrations/000014_create_subscription_tiers.down.sql`

**Domain Entities:**
- `internal/domain/entity/subscription_plan.go`
- `internal/domain/entity/tag_tier_mapping.go`

**Repository Interfaces:**
- `internal/domain/repository/subscription_plan_repository.go`
- `internal/domain/repository/tag_tier_mapping_repository.go`

**Service Interfaces:**
- `internal/domain/service/plan_management_service.go`
- `internal/domain/service/tag_tier_service.go`
- `internal/domain/service/content_access_service.go`

**DTOs:**
- `internal/application/dto/subscription_plan.go`

**Documentation:**
- `docs/api/subscription_tiers_api.yaml` (OpenAPI 3.0 spec)

---

### 📝 Files to Modify (Implementation Phase)

**Repository Implementations:**
- `internal/infrastructure/persistence/postgres/repository/subscription_plan_repository_impl.go` (new)
- `internal/infrastructure/persistence/postgres/repository/tag_tier_mapping_repository_impl.go` (new)

**Service Implementations:**
- `internal/domain/service/plan_management_service_impl.go` (new)
- `internal/domain/service/tag_tier_service_impl.go` (new)
- `internal/domain/service/content_access_service_impl.go` (new)
- `internal/domain/service/payment_service.go` (modify webhook handler)

**Handlers:**
- `internal/interfaces/http/handler/plan/handler.go` (new)
- `internal/interfaces/http/handler/plan/routes.go` (new)

**Router:**
- `internal/interfaces/http/router/router.go` (modify)

**Dependency Injection:**
- `cmd/api/modules/repository.go` (modify)
- `cmd/api/modules/service.go` (modify)
- `cmd/api/modules/handler.go` (modify)
- `cmd/api/modules/router.go` (modify)

**Swagger:**
- `docs/swagger.yaml` (regenerate after handler annotations)

---

## 15. Handoff Checklist for Planner

### ✅ Design Artifacts Completed

- [x] Database schema designed (migration files created)
- [x] API contracts defined (OpenAPI spec created)
- [x] Domain entities defined (Go structs)
- [x] Repository interfaces defined (contracts only)
- [x] Service interfaces defined (contracts only)
- [x] DTO definitions created
- [x] Integration points identified (payment webhook)
- [x] Routing structure planned
- [x] Testing strategy outlined

---

### ✅ User Confirmations Received

- [x] All design questions answered
- [x] Schema design approved
- [x] API contract approved
- [x] Tier storage approach confirmed (VARCHAR vs ENUM)
- [x] Price validation strategy confirmed (warn only)
- [x] Access check endpoint visibility confirmed (public)
- [x] Payment integration approach confirmed (reuse existing)

---

### 📋 Ready for Planning Phase

**Next Steps:**
1. **Planner Agent** will break down this design into atomic, sequential tasks
2. Tasks will follow TDD workflow: RED → GREEN → REFACTOR
3. Each task will be implementable in one TDD cycle
4. Tasks will be organized by dependency order (migrations → repositories → services → handlers)

---

## 16. Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         HTTP Layer                              │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  PlanHandler                                             │  │
│  │  - UpsertPlans()                                         │  │
│  │  - GetAuthorPlans()                                      │  │
│  │  - AssignTagToTier()                                     │  │
│  │  - UnassignTagFromTier()                                 │  │
│  │  - GetAuthorTagTiers()                                   │  │
│  │  - CheckBlogAccess()                                     │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                       Service Layer                             │
│  ┌──────────────────┐  ┌─────────────────┐  ┌────────────────┐ │
│  │ PlanManagement   │  │ TagTierService  │  │ ContentAccess  │ │
│  │ Service          │  │                 │  │ Service        │ │
│  └──────────────────┘  └─────────────────┘  └────────────────┘ │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ PaymentService (Modified - webhook handler)             │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                     Repository Layer                            │
│  ┌───────────────────┐  ┌──────────────────┐  ┌──────────────┐ │
│  │ SubscriptionPlan  │  │ TagTierMapping   │  │ Subscription │ │
│  │ Repository        │  │ Repository       │  │ Repository   │ │
│  └───────────────────┘  └──────────────────┘  └──────────────┘ │
│  ┌───────────────────┐  ┌──────────────────┐                   │
│  │ Tag Repository    │  │ Blog Repository  │                   │
│  └───────────────────┘  └──────────────────┘                   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                       Database Layer                            │
│  ┌─────────────────┐  ┌──────────────────┐  ┌────────────────┐ │
│  │ subscription_   │  │ tag_tier_        │  │ subscriptions  │ │
│  │ plans           │  │ mappings         │  │                │ │
│  └─────────────────┘  └──────────────────┘  └────────────────┘ │
│  ┌─────────────────┐  ┌──────────────────┐                     │
│  │ tags            │  │ blogs            │                     │
│  └─────────────────┘  └──────────────────┘                     │
└─────────────────────────────────────────────────────────────────┘
```

---

## 17. Success Criteria

### Functional Requirements
- [ ] Authors can create 3 paid plans (BRONZE/SILVER/GOLD) via API
- [ ] Authors can assign tags to tiers via API
- [ ] Readers can view pricing table for any author
- [ ] Readers can purchase plans via existing SePay flow
- [ ] Blog access is correctly gated based on user tier vs. required tier
- [ ] Upgrades work: tier replaced, expiry reset to NOW()+30 days
- [ ] All tests pass (unit + integration)
- [ ] Swagger docs updated with new endpoints

### Non-Functional Requirements
- [ ] Clean Architecture maintained
- [ ] SOLID principles applied
- [ ] Repository pattern used
- [ ] Service pattern used
- [ ] All interfaces defined before implementation
- [ ] Test coverage ≥ 80%

---

## 18. Out of Scope (Non-Goals)

- Multiple duration options (90/365 days) - MVP uses 30 days only
- Prorated upgrades - MVP uses simple full-price replacement
- Downgrades - Not allowed in MVP
- Refunds - Not supported
- Admin-defined tier names - System uses fixed FREE/BRONZE/SILVER/GOLD
- Per-blog pricing - Only tier-based pricing via tags

---

**End of Design Summary**

**Handoff Status:** ✅ READY FOR PLANNER  
**Next Agent:** Planner Agent  
**Next Action:** Break down into atomic implementation tasks
