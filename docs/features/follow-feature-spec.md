# Refined Spec: Follow/Following Model

## Summary

Implement a bidirectional follow relationship system (Twitter-style) that allows users to follow each other, view their followers/following lists, and see follow counts on user profiles. This feature enables social connections and will support future feed functionality.

---

## User Stories

### US-1: Follow a User
**As a** logged-in user,  
**I want to** follow another user,  
**So that** I can see their public content in my feed and stay updated with their activity.

### US-2: Unfollow a User
**As a** logged-in user,  
**I want to** unfollow a user I'm currently following,  
**So that** I can stop seeing their content in my feed.

### US-3: View My Followers
**As a** logged-in user,  
**I want to** see a list of users who follow me,  
**So that** I can understand my audience and potentially follow them back.

### US-4: View Users I'm Following
**As a** logged-in user,  
**I want to** see a list of users I follow,  
**So that** I can manage my following list and unfollow if needed.

### US-5: Check Follow Status
**As a** logged-in user viewing another user's profile,  
**I want to** see if I'm following them or if they follow me,  
**So that** I know the relationship status between us.

### US-6: View Follow Counts (Public)
**As a** any user (guest or logged-in) viewing a profile,  
**I want to** see the follower and following counts,  
**So that** I can understand the user's social presence and influence.

---

## Functional Requirements

### FR-1: Follow Management
- **FR-1.1**: Users can follow other users via API
- **FR-1.2**: Users can unfollow users they are following
- **FR-1.3**: System must prevent duplicate follow relationships
- **FR-1.4**: System must prevent users from following themselves

### FR-2: List Management
- **FR-2.1**: Retrieve paginated list of followers for any user
- **FR-2.2**: Retrieve paginated list of users being followed (following) for any user
- **FR-2.3**: Support pagination with configurable page size (default: 20, max: 100)
- **FR-2.4**: Support sorting by follow date (newest first)

### FR-3: Status & Counts
- **FR-3.1**: Check if user A follows user B (boolean check)
- **FR-3.2**: Get follower count for any user
- **FR-3.3**: Get following count for any user
- **FR-3.4**: Follow counts must be publicly accessible (no auth required)

### FR-4: Data Storage
- **FR-4.1**: Store follow relationship with unique composite key (follower_id, following_id)
- **FR-4.2**: Store creation timestamp for each follow relationship
- **FR-4.3**: Soft delete is NOT required (hard delete on unfollow)

### FR-5: API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/users/{userId}/follow` | Required | Follow a user |
| POST | `/api/v1/users/{userId}/unfollow` | Required | Unfollow a user |
| GET | `/api/v1/users/{userId}/followers` | Optional | Get user's followers |
| GET | `/api/v1/users/{userId}/following` | Optional | Get users being followed |
| GET | `/api/v1/users/{userId}/followers/count` | None | Get follower count |
| GET | `/api/v1/users/{userId}/following/count` | None | Get following count |
| GET | `/api/v1/me/following` | Required | Get current user's following list |
| GET | `/api/v1/me/followers` | Required | Get current user's followers |

---

## Non-Functional Requirements

### Performance Requirements
- **NFR-1.1**: Follow/unfollow operations must complete within 200ms (p95)
- **NFR-1.2**: Follower/following list queries must complete within 300ms (p95) for users with up to 10K followers
- **NFR-1.3**: Count queries must complete within 100ms (p95)
- **NFR-1.4**: Support users with up to 1M followers (with appropriate pagination)

### Security Requirements
- **NFR-2.1**: Only authenticated users can follow/unfollow
- **NFR-2.2**: Users cannot follow themselves (self-follow prevention)
- **NFR-2.3**: Follow lists for private accounts may be restricted (future consideration)
- **NFR-2.4**: Rate limiting: Max 100 follow/unfollow actions per hour per user

### Scalability Requirements
- **NFR-3.1**: Database schema must support indexing for efficient queries
- **NFR-3.2**: Composite index on (follower_id, following_id) for uniqueness and lookups
- **NFR-3.3**: Index on following_id for follower list queries
- **NFR-3.4**: Index on follower_id for following list queries
- **NFR-3.5**: Support for database sharding if user base grows beyond 10M users

### Data Consistency
- **NFR-4.1**: Follow relationship must be atomic (no partial writes)
- **NFR-4.2**: Counts should be accurate (consider count caching for high-traffic users)

---

## Edge Cases

### EC-1: Self-Follow Prevention
**Scenario**: User attempts to follow themselves  
**Expected Behavior**: 
- API returns 400 Bad Request
- Error message: "cannot follow yourself"
- No database record created

### EC-2: Duplicate Follow Handling
**Scenario**: User attempts to follow someone they already follow  
**Expected Behavior**:
- API returns 409 Conflict
- Error message: "already following this user"
- No duplicate database record created
- Existing relationship unchanged

### EC-3: Unfollow Non-Existing Relationship
**Scenario**: User attempts to unfollow someone they don't follow  
**Expected Behavior**:
- API returns 404 Not Found
- Error message: "follow relationship not found"
- No database changes

### EC-4: Large Follower Lists (Pagination)
**Scenario**: User has 100K+ followers  
**Expected Behavior**:
- API returns paginated results (default 20 per page)
- Maximum page size: 100
- Efficient query performance via database indexing
- Cursor-based pagination considered for future (offset is fine for MVP)

### EC-5: User Deletion - Cascade Behavior
**Scenario**: User account is deleted  
**Expected Behavior**:
- All follow relationships where user is follower are deleted (cascade)
- All follow relationships where user is followed are deleted (cascade)
- Follower/following counts for other users are updated accordingly

### EC-6: Following Deleted User
**Scenario**: User A follows User B, then User B is deleted  
**Expected Behavior**:
- Follow relationship is automatically removed via cascade
- User A's following count decrements
- User A can no longer see User B in their following list

### EC-7: Concurrent Follow Requests
**Scenario**: Two simultaneous requests to follow the same user  
**Expected Behavior**:
- Database unique constraint prevents duplicate
- One request succeeds (201 Created)
- Other request fails with 409 Conflict

### EC-8: Invalid User ID
**Scenario**: Request with non-existent or invalid user ID  
**Expected Behavior**:
- API returns 400 Bad Request for invalid UUID format
- API returns 404 Not Found for non-existent user

### EC-9: Guest User Access
**Scenario**: Unauthenticated user attempts to follow  
**Expected Behavior**:
- API returns 401 Unauthorized
- Error message: "authentication required"

### EC-10: Rate Limiting
**Scenario**: User exceeds follow/unfollow rate limit (100/hour)  
**Expected Behavior**:
- API returns 429 Too Many Requests
- Retry-After header included
- Error message: "rate limit exceeded"

---

## Acceptance Criteria

### AC-1: Follow User
```gherkin
Given I am an authenticated user
And the target user exists
And I am not following the target user
And the target user is not me
When I send a POST request to /api/v1/users/{userId}/follow
Then I receive a 201 Created response
And the response contains the follow relationship details
And I am now following the target user
```

### AC-2: Unfollow User
```gherkin
Given I am an authenticated user
And I am following the target user
When I send a POST request to /api/v1/users/{userId}/unfollow
Then I receive a 204 No Content response
And I am no longer following the target user
```

### AC-3: View Followers List
```gherkin
Given a user has followers
When I send a GET request to /api/v1/users/{userId}/followers
Then I receive a 200 OK response
And the response contains a paginated list of followers
And each follower includes user details (id, name, avatar)
And pagination metadata is included (page, pageSize, total, totalPages)
```

### AC-4: View Following List
```gherkin
Given a user follows other users
When I send a GET request to /api/v1/users/{userId}/following
Then I receive a 200 OK response
And the response contains a paginated list of followed users
And each user includes basic profile details
And pagination metadata is included
```

### AC-5: View Follow Counts (Public)
```gherkin
Given any user (authenticated or guest)
When I send a GET request to /api/v1/users/{userId}/followers/count
Then I receive a 200 OK response
And the response contains the follower count

When I send a GET request to /api/v1/users/{userId}/following/count
Then I receive a 200 OK response
And the response contains the following count
```

### AC-6: Self-Follow Prevention
```gherkin
Given I am an authenticated user
When I attempt to follow myself via POST /api/v1/users/{myUserId}/follow
Then I receive a 400 Bad Request response
And the error message indicates "cannot follow yourself"
And no follow relationship is created
```

### AC-7: Duplicate Follow Prevention
```gherkin
Given I am an authenticated user
And I am already following the target user
When I attempt to follow the same user again
Then I receive a 409 Conflict response
And the error message indicates "already following this user"
And no duplicate relationship is created
```

### AC-8: Check Follow Status
```gherkin
Given I am an authenticated user
When I view another user's profile (via existing profile endpoint)
Then I can see if I am following them (isFollowing: true/false)
And I can see if they are following me (isFollower: true/false)
```

---

## Tech Stack (Auto-Detected)

| Component | Technology |
|-----------|------------|
| **Language** | Go 1.24.0 |
| **Framework** | Gin (github.com/gin-gonic/gin) |
| **Database** | PostgreSQL |
| **ORM** | GORM (gorm.io/gorm) |
| **Testing** | Go testing + testify |
| **Architecture** | Clean Architecture / Domain-Driven Design |
| **API Documentation** | Swagger (swaggo) |
| **DI Framework** | Uber FX (go.uber.org/fx) |
| **Cache** | Redis (go-redis) |

### Existing Patterns to Follow
- **Entity**: UUID primary keys, GORM tags, TableName() method
- **Repository**: Interface in domain layer, implementation in persistence layer
- **Service**: Domain logic with error variables exported
- **UseCase**: Orchestration, DTO conversion
- **Handler**: Gin handlers with Swagger annotations
- **Pagination**: repository.Pagination and repository.PaginatedResult[T]

---

## Data Model

### Follow Entity

```go
type Follow struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    FollowerID  uuid.UUID `gorm:"type:uuid;not null;index"`  // The user who follows
    FollowingID uuid.UUID `gorm:"type:uuid;not null;index"`  // The user being followed
    CreatedAt   time.Time `gorm:"not null;default:now()"`
    
    // Relationships
    Follower  *User `gorm:"foreignKey:FollowerID"`
    Following *User `gorm:"foreignKey:FollowingID"`
}
```

### Database Schema

```sql
CREATE TABLE follows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(follower_id, following_id)
);

-- Indexes for efficient queries
CREATE INDEX idx_follows_follower_id ON follows(follower_id);
CREATE INDEX idx_follows_following_id ON follows(following_id);
CREATE INDEX idx_follows_created_at ON follows(created_at);
```

### User Entity Updates

Add to existing User entity:
```go
// In internal/domain/entity/user.go
// Add to User struct:
Followers []Follow `gorm:"foreignKey:FollowingID" json:"followers,omitempty"`  // Users following me
Following []Follow `gorm:"foreignKey:FollowerID" json:"following,omitempty"`   // Users I follow
```

---

## Dependencies

### Required Components
1. **User Repository** - To validate user existence
2. **Authentication Middleware** - To get current user ID from JWT
3. **Database Migration System** - To create follows table

### New Components to Create
1. **Entity**: `internal/domain/entity/follow.go`
2. **Repository Interface**: `internal/domain/repository/follow_repository.go`
3. **Repository Implementation**: `internal/infrastructure/persistence/postgres/repository/follow_repository.go`
4. **Domain Service**: `internal/domain/service/follow_service.go`
5. **UseCase**: `internal/application/usecase/follow_usecase.go`
6. **DTOs**: `internal/application/dto/follow.go`
7. **Handler**: `internal/interfaces/http/handler/follow_handler.go`
8. **Migration**: `migrations/000005_follows.up.sql` and `.down.sql`

---

## Out of Scope

The following features are explicitly NOT included in this iteration:

1. **Mutual Follow Indication** - Showing "follows you" badge (can be added later)
2. **Follow Requests** - Private account follow approval system
3. **Follow Suggestions** - AI/algorithm-based user recommendations
4. **Follow Notifications** - Push/email notifications when someone follows
5. **Bulk Follow Operations** - Import/follow multiple users at once
6. **Follow Analytics** - Detailed stats on follower growth/loss
7. **Block Functionality** - Preventing specific users from following
8. **Feed Algorithm** - Content feed based on follows (separate feature)
9. **Follower Privacy Settings** - Hide followers/following lists
10. **Verified/Bot Labels** - Special follower categorization

---

## Open Questions

> **Note**: The following questions need clarification before implementation begins:

1. **Q1: Relationship with Subscriptions**  
   The platform already has a "Subscriptions" feature (for subscribing to authors). Should "Follow" be a separate social feature, or should we consolidate?  
   **Assumption**: Keep separate - Subscriptions = content access, Follows = social connection

2. **Q2: Follower List Privacy**  
   Should users be able to hide their followers/following lists?  
   **Assumption**: Public by default for MVP, privacy settings in future iteration

3. **Q3: Soft Delete vs Hard Delete**  
   Should unfollow be a soft delete (keep history) or hard delete?  
   **Assumption**: Hard delete - no business need to track unfollow history for MVP

4. **Q4: Rate Limiting Values**  
   Is 100 follow/unfollow actions per hour reasonable?  
   **Assumption**: 100/hour is standard, adjustable based on usage patterns

5. **Q5: Follow Count Caching**  
   Should we cache follower counts in Redis for high-traffic users?  
   **Assumption**: Not for MVP - add caching if performance issues arise

---

## Implementation Checklist

### Phase 1: Database & Domain Layer
- [ ] Create migration files for follows table
- [ ] Create Follow entity
- [ ] Create FollowRepository interface
- [ ] Create FollowRepository PostgreSQL implementation
- [ ] Write repository unit tests

### Phase 2: Business Logic Layer
- [ ] Create FollowService with business rules
- [ ] Create FollowUseCase
- [ ] Create Follow DTOs
- [ ] Write service and usecase unit tests

### Phase 3: API Layer
- [ ] Create FollowHandler with all endpoints
- [ ] Add Swagger annotations
- [ ] Wire up dependencies in FX module
- [ ] Write handler integration tests

### Phase 4: Integration & Testing
- [ ] Run all tests (unit + integration)
- [ ] Test edge cases manually
- [ ] Verify API documentation
- [ ] Performance test with 10K+ followers

---

## References

- **Similar Implementation**: See `internal/domain/entity/subscription.go` and related files
- **Pagination Pattern**: `internal/domain/repository/base.go`
- **Error Handling**: `internal/domain/service/subscription_service.go`
- **Handler Pattern**: `internal/interfaces/http/handler/subscription_handler.go`
- **Migration Pattern**: `migrations/000001_init.up.sql`

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-02-01 | Gatekeeper Agent | Initial specification |
