# Architecture Design Complete ✅

## Multi-Tier Subscription Plans Feature

**Date:** 2026-02-03  
**Architect:** Claude (Architect Agent)  
**Status:** Design Complete - Ready for Implementation Planning

---

## 📋 Design Artifacts Created

### 1. Database Schema ✅
**Location:** `migrations/000014_create_subscription_tiers.{up,down}.sql`

**New Tables:**
- `subscription_plans` - Author-defined pricing plans (BRONZE/SILVER/GOLD)
- `tag_tier_mappings` - Tag-to-tier assignments per author

**Key Features:**
- UNIQUE constraints for data integrity
- Composite indexes for fast lookups
- CHECK constraints for tier validation
- Auto-update triggers for `updated_at` columns
- CASCADE deletes for referential integrity

---

### 2. API Contracts ✅
**Location:** `docs/api/subscription_tiers_api.yaml`

**Endpoints Designed:**
1. `POST /api/v1/authors/me/plans` - Upsert subscription plans
2. `GET /api/v1/authors/{authorId}/plans` - Public pricing table
3. `POST /api/v1/authors/me/tags/{tagId}/tier` - Assign tag to tier
4. `DELETE /api/v1/authors/me/tags/{tagId}/tier` - Unassign tag from tier
5. `GET /api/v1/authors/me/tag-tiers` - Get all tag-tier mappings
6. `GET /api/v1/blogs/{blogId}/access` - Check blog access (public)
7. `POST /api/v1/payments` - Modified to support `plan_id` parameter

**Full OpenAPI 3.0 specification included with:**
- Request/response schemas
- Validation rules
- Example payloads
- Error responses
- Authentication requirements

---

### 3. Domain Entities ✅
**Location:** `internal/domain/entity/`

**New Entities:**
- `SubscriptionTier` enum with `Level()` method for tier comparison
- `SubscriptionPlan` - Pricing plan entity
- `TagTierMapping` - Tag-tier assignment entity

**Key Methods:**
```go
func (t SubscriptionTier) Level() int  // Returns 0-3 for FREE-GOLD
func (t SubscriptionTier) IsValid() bool
func (t SubscriptionTier) String() string
```

---

### 4. Repository Interfaces ✅
**Location:** `internal/domain/repository/`

**New Repositories:**
- `SubscriptionPlanRepository` - 10 methods including Upsert
- `TagTierMappingRepository` - 11 methods including blog count

**Key Features:**
- Upsert support for ON CONFLICT handling
- Transaction support via `WithTx()`
- Bulk operations (FindByTagIDs)
- Blog count queries for "affected blogs" feature

---

### 5. Service Interfaces ✅
**Location:** `internal/domain/service/`

**New Services:**
- `PlanManagementService` - Plan CRUD + price validation
- `TagTierService` - Tag-tier mapping + blog access logic
- `ContentAccessService` - Blog access checks + upgrade options

**Key Business Logic:**
- Price hierarchy validation (warns, doesn't block)
- Tier comparison using `Level()` method
- Upgrade option generation for blocked content
- Anonymous user support (userTier = FREE)

---

### 6. DTOs ✅
**Location:** `internal/application/dto/subscription_plan.go`

**Request DTOs:**
- `CreatePlanRequest`, `UpsertPlansRequest`
- `AssignTagTierRequest`

**Response DTOs:**
- `PlanResponse`, `UpsertPlansResponse` (with warnings)
- `PlanWithTagsResponse` (for public pricing table)
- `GetAuthorPlansResponse`
- `AssignTagTierResponse`, `UnassignTagTierResponse`
- `TagTierMappingResponse`, `GetTagTiersResponse`
- `CheckBlogAccessResponse` (with upgrade options)
- `UpgradeOption`

---

### 7. Design Documentation ✅
**Location:** `docs/SUBSCRIPTION_TIERS_DESIGN.md`

**Comprehensive 500+ line design document including:**
- Executive summary
- Database schema design with constraints
- API contract specifications
- Domain entity design
- Repository interface design
- Service layer architecture
- DTO definitions
- Integration with existing payment flow
- Dependency injection setup
- Testing strategy
- Architecture diagrams
- Success criteria
- Handoff checklist for Planner

---

## 🏗️ Architecture Overview

```
HTTP Layer (Handlers)
    ↓
Service Layer (Business Logic)
  - PlanManagementService
  - TagTierService
  - ContentAccessService
  - PaymentService (modified)
    ↓
Repository Layer (Data Access)
  - SubscriptionPlanRepository
  - TagTierMappingRepository
  - SubscriptionRepository
  - TagRepository
  - BlogRepository
    ↓
Database Layer
  - subscription_plans
  - tag_tier_mappings
  - subscriptions (existing)
  - tags (existing)
  - blogs (existing)
```

---

## ✅ Design Principles Applied

### SOLID Principles
- **Single Responsibility:** Each service has one clear purpose
- **Dependency Inversion:** Services depend on repository interfaces
- **Interface Segregation:** Small, focused interfaces

### Design Patterns
- **Repository Pattern:** Data access isolated
- **Service Pattern:** Business logic in services
- **DTO Pattern:** Clear API contracts

---

## 🔑 Key Design Decisions (User-Confirmed)

1. ✅ **Tags are globally shared** - tier assignments per-author
2. ✅ **Price validation** - Warn only, no hard enforcement
3. ✅ **Tier storage** - VARCHAR(20) for migration flexibility
4. ✅ **Plan API** - Single POST endpoint for upsert
5. ✅ **Access check** - Public endpoint (anonymous users = FREE tier)
6. ✅ **Purchase API** - Reuse existing with `plan_id` parameter

---

## 📝 Integration Points

### Modified: PaymentService.HandleSePayWebhook
**Changes Required:**
1. Add `planRepo` dependency
2. Fetch plan by `transaction.PlanID`
3. Set `subscription.Tier = plan.Tier` on create/update
4. Reset expiry to `NOW() + plan.DurationDays`

**Upgrade Logic:**
- Full price payment (no prorating)
- Tier replacement (not additive)
- Expiry reset to NOW() + 30 days

---

## 🧪 Testing Strategy

### Unit Tests
- Tier enum Level() method
- Price hierarchy validation
- Access control logic
- Webhook tier update

### Integration Tests
- Complete subscription purchase workflow
- Tag-tier mapping workflow
- Upgrade workflow
- Access check workflow

**Coverage Target:** ≥ 80%

---

## 📦 Files Summary

### Created (10 files)
✅ Migrations (2)
✅ Domain Entities (2)
✅ Repository Interfaces (2)
✅ Service Interfaces (3)
✅ DTOs (1)
✅ API Spec (1)
✅ Design Doc (1)

### To Modify (9 files in implementation phase)
- Repository implementations (2 new)
- Service implementations (3 new, 1 modified)
- Handler (1 new)
- Router (1 modified)
- DI modules (4 modified)

---

## ✅ Handoff Checklist

**Design Phase Complete:**
- [x] All design questions answered
- [x] User approved schema design
- [x] User approved API contracts
- [x] Database migrations created
- [x] Domain entities defined
- [x] Repository interfaces defined
- [x] Service interfaces defined
- [x] DTOs created
- [x] Integration points identified
- [x] Testing strategy outlined
- [x] No implementation code (contracts only)

**Ready for Next Phase:**
- [ ] Planner Agent to create task breakdown
- [ ] Builder Agent to implement via TDD
- [ ] Reviewer Agent to verify implementation

---

## 🎯 Success Criteria

### Functional
- Authors can create/update subscription plans
- Authors can assign tags to tiers
- Readers can view pricing tables
- Readers can purchase subscriptions
- Blog access is correctly gated
- Upgrades work with tier replacement

### Technical
- Clean Architecture maintained
- SOLID principles applied
- Test coverage ≥ 80%
- Swagger docs updated

---

## 🚀 Next Steps

**Immediate Next Step:**
→ **Planner Agent** will create atomic implementation tasks

**Task Breakdown Will Include:**
1. Run database migrations
2. Implement repository layer
3. Implement service layer
4. Implement handler layer
5. Update router
6. Update DI containers
7. Write tests
8. Update Swagger docs

**Implementation Approach:**
- TDD workflow: RED → GREEN → REFACTOR
- One task per TDD cycle
- Sequential execution by dependency order

---

**Design Status:** ✅ COMPLETE  
**Next Agent:** Planner  
**Blocked By:** None  
**Ready to Proceed:** Yes

---

**End of Architecture Design Phase**
