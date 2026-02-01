---
name: schema-design
description: Design database schemas with proper normalization and indexing
---

# Schema Design Skill

## Purpose

Design database schemas based on requirements with proper normalization, indexing, and constraints.

## When to Use

- When designing new tables
- When modifying schema for new features
- When optimizing database structure

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

## Schema Templates

### PostgreSQL

```sql
CREATE TABLE [table_name] (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  [columns...]
  created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_[table]_[column] ON [table]([column]);

-- Foreign Keys
ALTER TABLE [table]
  ADD CONSTRAINT fk_[table]_[ref]
  FOREIGN KEY ([column]) REFERENCES [ref_table](id)
  ON DELETE CASCADE;
```

### Prisma

```prisma
model [ModelName] {
  id        String   @id @default(uuid())
  [fields...]
  createdAt DateTime @default(now()) @map("created_at")
  updatedAt DateTime @updatedAt @map("updated_at")

  @@map("[table_name]")
}
```

### Drizzle

```typescript
export const [tableName] = pgTable('[table_name]', {
  id: uuid('id').primaryKey().defaultRandom(),
  [columns...],
  createdAt: timestamp('created_at').notNull().defaultNow(),
  updatedAt: timestamp('updated_at').notNull().defaultNow(),
});
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
