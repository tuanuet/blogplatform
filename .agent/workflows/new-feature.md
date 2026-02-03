---
description: New feature with migration, feature flags, and rollback plan
---

# New Feature Workflow

User Request: $1

Full pipeline **PLUS** production-ready extras: migration, feature flags, rollback.

## When to Use

- Adding significant new functionality
- Features that need database changes
- Features that need gradual rollout

## Phases

### Phase 1-3: Follow Pipeline

Run the standard 3-phase pipeline first.

### Phase 4: MIGRATION PLAN (if schema changes)

```
1. Check if new tables/columns are needed
2. Create migration file
3. Plan backward compatibility
4. Document rollback migration
```

### Phase 5: FEATURE FLAG (optional)

```
1. Wrap new feature behind flag
2. Default: OFF in production
3. Document flag name and purpose
4. Plan gradual rollout percentage
```

### Phase 6: ROLLBACK PLAN

```
1. Document how to disable feature
2. Migration rollback commands
3. Cache invalidation steps
4. Communication plan if issues
```

## Output Checklist

**From Pipeline:**

- [ ] Refined Spec
- [ ] Schema + API Contract
- [ ] Tests + Implementation

**New Feature Extras:**

- [ ] Migration file created
- [ ] Rollback migration ready
- [ ] Feature flag implemented (if needed)
- [ ] Rollback plan documented
- [ ] Deployment checklist ready

## Example

```
Feature: "Add user subscription tiers"

Pipeline output:
- Schema: subscriptions table
- API: /api/v1/subscriptions

New Feature extras:
- Migration: 001_create_subscriptions.sql
- Rollback: 001_drop_subscriptions.sql
- Feature flag: ENABLE_SUBSCRIPTIONS=false
- Rollback plan: Set flag to false, run rollback migration
```
