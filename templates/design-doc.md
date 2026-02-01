# Design Doc: Blog Reactions

## 1. Schema Design

### New Entity: `BlogReaction`

File: `internal/domain/entity/blog_reaction.go`

```go
package entity

import (
	"time"

	"github.com/google/uuid"
)

type ReactionType string

const (
	ReactionTypeUpvote   ReactionType = "upvote"
	ReactionTypeDownvote ReactionType = "downvote"
)

// BlogReaction tracks user reactions to blogs
type BlogReaction struct {
	ID        uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BlogID    uuid.UUID    `gorm:"type:uuid;not null;index:idx_blog_user,unique" json:"blogId"`
	UserID    uuid.UUID    `gorm:"type:uuid;not null;index:idx_blog_user,unique" json:"userId"`
	Type      ReactionType `gorm:"type:varchar(20);not null" json:"type"`
	CreatedAt time.Time    `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time    `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	Blog *Blog `gorm:"foreignKey:BlogID" json:"blog,omitempty"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (BlogReaction) TableName() string {
	return "blog_reactions"
}
```

### Modified Entity: `Blog`

File: `internal/domain/entity/blog.go`

Add fields for denormalized counts to optimize read performance.

```go
type Blog struct {
    // ... existing fields ...
    UpvoteCount   int `gorm:"not null;default:0" json:"upvoteCount"`
    DownvoteCount int `gorm:"not null;default:0" json:"downvoteCount"`
}
```

## 2. API Contract

### Endpoint: Manage Reaction

**POST** `/blogs/{id}/reaction`

- **Summary**: Upvote, downvote, or remove reaction for a blog.
- **Auth**: Required (Bearer Token)

**Request Body**:

```json
{
  "reaction": "upvote" | "downvote" | "none"
}
```

**Response (200 OK)**:
Returns the updated counts and the user's current reaction status.

```json
{
  "success": true,
  "data": {
    "blogId": "uuid...",
    "upvoteCount": 15,
    "downvoteCount": 2,
    "userReaction": "upvote" // or "downvote" or null
  }
}
```

**Response (400 Bad Request)**:

- Invalid reaction type.
- Invalid Blog ID.

**Response (401 Unauthorized)**:

- User not logged in.

**Response (404 Not Found)**:

- Blog not found.

## 3. Implementation Details

### Service Logic (`BlogService.React`)

1.  **Input**: `blogID`, `userID`, `newReactionType`.
2.  **Transaction**:
    - **Lock**: Select `Blog` (optional but good for consistency) or just proceed.
    - **Fetch Existing**: Check if `BlogReaction` exists for `(blogID, userID)`.
    - **Logic**:
      - **Case 1: No previous reaction**:
        - If `new == none`: Do nothing.
        - If `new != none`: Insert `BlogReaction`, Increment `Blog.Count`.
      - **Case 2: Previous reaction exists (e.g., Upvote)**:
        - If `new == same (Upvote)`: Do nothing (or toggle off -> treat as `none`?). _Decision: Idempotent - do nothing._
        - If `new == none`: Delete `BlogReaction`, Decrement `Blog.UpvoteCount`.
        - If `new == different (Downvote)`: Update `BlogReaction` type. Decrement `UpvoteCount`, Increment `DownvoteCount`.
3.  **Return**: New counts and status.

### Repository Interface

```go
type BlogRepository interface {
    // ... existing ...
    React(ctx context.Context, blogID, userID uuid.UUID, reaction entity.ReactionType) (upvotes, downvotes int, err error)
    RemoveReaction(ctx context.Context, blogID, userID uuid.UUID) (upvotes, downvotes int, err error)
    // Helper to get reaction if needed, though React can handle upsert logic
    GetReaction(ctx context.Context, blogID, userID uuid.UUID) (*entity.BlogReaction, error)
}
```

_Refinement_: To keep Service pure, the `Repository` might just provide atomic primitives `UpsertReaction` and `DeleteReaction` which return the new counts, or the Service orchestrates the transaction. Given the requirement for "optimization", doing it in one db transaction/query is best.

**Optimized Query Strategy (Postgres)**:
Use a CTE or stored procedure equivalent logic in Go transaction:

1.  Read existing reaction.
2.  Calculate deltas.
3.  Insert/Update/Delete reaction table.
4.  Update blog counts `SET upvote_count = upvote_count + ?`.
