---
name: schema-design
description: Design database schemas with proper normalization, indexing, and advanced patterns
---

# Schema Design Skill

## Purpose

Design database schemas based on requirements with proper normalization, indexing, and scalability considerations.

## When to Use

- When designing new tables
- When modifying schema for new features
- When optimizing database structure
- When designing for scale

## Design Principles

### 1. Normalization

- **1NF**: Atomic values, no repeating groups
- **2NF**: No partial dependencies
- **3NF**: No transitive dependencies
- **Denormalize** only with justification (performance)

### 2. Naming Conventions

- Tables: `snake_case`, plural (`users`, `orders`)
- Columns: `snake_case` (`created_at`, `user_id`)
- Foreign keys: `[table]_id` (`user_id`, `order_id`)
- Indexes: `idx_[table]_[columns]`

### 3. Required Columns

```sql
id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
```

### 4. Soft Delete (optional)

```sql
deleted_at  TIMESTAMP NULL
```

## Advanced Modeling

### Hierarchies (Trees)

Model tree structures efficiently.

**Option A: Adjacency List (Simple)**
```sql
CREATE TABLE categories (
    id UUID PRIMARY KEY,
    parent_id UUID REFERENCES categories(id),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Recursive CTE to fetch full tree (PostgreSQL 8.4+)
WITH RECURSIVE category_tree AS (
    SELECT id, parent_id, name
    FROM categories
    WHERE parent_id IS NULL

    UNION ALL

    SELECT c.id, c.parent_id, c.name
    FROM categories c
    INNER JOIN category_tree ct ON c.parent_id = ct.id
)
SELECT * FROM category_tree;
```

**Option B: Materialized Path (High Performance)**
```sql
CREATE TABLE categories (
    id UUID PRIMARY KEY,
    path VARCHAR(1000), -- e.g., "/1/4/7"
    name VARCHAR(255) NOT NULL
);

-- Fetch subtree efficiently
SELECT * FROM categories WHERE path LIKE '/1/%';
```

**Option C: Closure Table (Auditable, Reorderable)**
```sql
CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE category_closure (
    ancestor_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    descendant_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    depth INT NOT NULL,
    PRIMARY KEY (ancestor_id, descendant_id)
);
```

### Polymorphic Associations

Allow a table to belong to multiple entities.

```sql
CREATE TABLE comments (
    id UUID PRIMARY KEY,
    body TEXT NOT NULL,
    -- Polymorphic relation
    commentable_type VARCHAR(50), -- 'Blog', 'Series'
    commentable_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_comments_polymorphic ON comments(commentable_type, commentable_id);
```

### Many-to-Many with Attributes

Store data on the join table (Pivot data).

```sql
-- ❌ Bad: No metadata on user_bookmarks
CREATE TABLE user_bookmarks (
    user_id UUID NOT NULL,
    blog_id UUID NOT NULL,
    PRIMARY KEY (user_id, blog_id)
);

-- ✅ Good: Added metadata directly to bookmark
CREATE TABLE user_bookmarks (
    user_id UUID NOT NULL,
    blog_id UUID NOT NULL,
    folder_name VARCHAR(50), -- Custom folder for user
    is_favorite BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, blog_id)
);
```

## Performance Patterns

### Partitioning

Split large tables for query performance.

**Range Partitioning (Time-based)**
```sql
CREATE TABLE logs_2024 PARTITION OF logs FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
CREATE TABLE logs_2025 PARTITION OF logs FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');
```

**List Partitioning (Hash-based)**
```sql
CREATE TABLE user_events (
    user_id UUID NOT NULL,
    event_data JSONB
) PARTITION BY HASH (user_id);
```

### Materialized Views

Pre-compute expensive aggregations.

```sql
CREATE TABLE user_rankings_mv AS
SELECT
    user_id,
    SUM(score) as total_score,
    RANK() OVER (ORDER BY SUM(score) DESC) as ranking
FROM user_velocity_scores
GROUP BY user_id
WITH DATA;

-- Refresh strategy via cron job
REFRESH MATERIALIZED VIEW CONCURRENTLY user_rankings_mv;
```

### CQRS (Command Query Responsibility Segregation)

Separate read and write models.

**Write Model (Command)**
```go
type UserCommand struct {
    ID       uuid.UUID
    Email    string
    Password string
}
```

**Read Model (Query)**
```sql
CREATE TABLE user_read_model (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    -- Optimized for display/search, no heavy joins needed
    INDEX idx_name_trgm ON name USING GIN (gin_trgm_ops);
);
```

## PostgreSQL Specifics

### Partial Indexes

Index only a subset of rows (e.g., only active users).

```sql
-- Standard index (includes deleted users)
CREATE INDEX idx_users_email ON users(email);

-- Partial index (excludes deleted users)
CREATE INDEX idx_users_active_email ON users(email) WHERE deleted_at IS NULL;
```

### GIN & GiST Indexes

For full-text search and JSONB.

```sql
-- GIN for full-text search (requires pg_trgm extension)
CREATE INDEX idx_blogs_title_search ON blogs USING GIN (title gin_trgm_ops);

-- GiST for geometric or inclusion operators
CREATE INDEX idx_blogs_tags ON blogs USING GIST (tags);
```

### Upserts with `EXCLUDE`

Insert or update without race conditions.

```sql
INSERT INTO user_bookmarks (user_id, blog_id, created_at)
VALUES ($1, $2, NOW())
ON CONFLICT (user_id, blog_id) DO NOTHING;
```

### JSONB Operations

Querying nested data inside JSONB columns.

```sql
-- Find blogs with specific tag
SELECT * FROM blogs
WHERE tags @> '"golang"'::jsonb;

-- Update a specific key in metadata
UPDATE users
SET metadata = jsonb_set(metadata, '{"theme": "dark"}')
WHERE id = $1;
```

## Schema Migration

### Evolutionary Design

Never `DROP TABLE`. Always use `ALTER TABLE`.

```sql
-- ✅ Good: Add column
ALTER TABLE users ADD COLUMN phone_number VARCHAR(20);

-- ❌ Bad: Recreate table
DROP TABLE users;
CREATE TABLE users (...new schema...);
```

### Backward Compatibility

Default values for new columns to prevent errors.

```sql
-- Add new column with default
ALTER TABLE blogs ADD COLUMN view_count INT DEFAULT 0;
```

### Data Seeding

Populate reference data (countries, default roles).

```sql
-- PostgreSQL upsert for seeding
INSERT INTO roles (name, permissions) VALUES
('admin', '["all"]::jsonb),
('user', '["read"]::jsonb')
ON CONFLICT (name) DO NOTHING;
```

## GORM Integration

### Soft Deletes (Scopes)

Add default scopes to exclude deleted records automatically.

```go
// ❌ Bad: Check manually everywhere
db.Where("deleted_at IS NULL").Find(&blogs)

// ✅ Good: Define scope
func (Blog) Scope(q *gorm.DB) *gorm.DB {
    return q.Where("deleted_at IS NULL")
}

// Usage
db.Scopes(Blog.Scope{}).Find(&blogs)
```

### Preloading (Eager Loading)

Avoid N+1 queries by preloading relations.

```go
// ❌ Bad: N+1 queries
var blogs []Blog
db.Find(&blogs)
for _, blog := range blogs {
    db.Where("blog_id = ?", blog.ID).Find(&blog.Tags) // N queries!
}

// ✅ Good: Single query with preload
var blogs []Blog
db.Preload("Tags").Preload("Author").Find(&blogs)
```

### Polymorphism with GORM

Handle associations to multiple models.

```go
type Comment struct {
    ID          uuid.UUID
    Body        string
    // Polymorphic fields
    EntityType   string `gorm:"size:50"` // "blog", "series"
    EntityID    uuid.UUID `gorm:"size:19;index:idx_entity"` // Composite index
}

// Helper methods
func (c *Comment) GetBlog(db *gorm.DB) *Blog {
    var blog Blog
    db.Where("id = ? AND entity_type = ?", c.EntityID, "blog").First(&blog)
    return &blog
}
```

### Hooks for Business Logic

Enforce rules at the database level.

```go
type User struct {
    // ...
    CreatedAt time.Time
}

// Before save: set defaults
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.CreatedAt.IsZero() {
        u.CreatedAt = time.Now()
    }
    return nil
}
```

## Relationship Types

| Type             | Implementation            |
| ---------------- | ------------------------- |
| One-to-Many      | FK on "many" side         |
| Many-to-Many     | Junction table            |
| One-to-One       | FK with unique constraint |
| Self-referencing | FK points to same table   |

## Index Guidelines

- Always index foreign keys
- Index columns used in WHERE clauses
- Index columns used in ORDER BY
- Consider composite indexes for multi-column queries
- Avoid over-indexing (impacts write performance)

## Validation Checklist

- [ ] All tables have id, created_at, updated_at
- [ ] Foreign keys are indexed
- [ ] Naming conventions are consistent
- [ ] Constraints are defined (NOT NULL, UNIQUE, CHECK)
- [ ] Relationships are properly defined
- [ ] Indexes support query patterns
- [ ] Soft deletes handled (if applicable)
