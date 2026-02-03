---
description: Safe refactoring (Builder phase only)
---

# Refactor Workflow

User Request: $1

Minimal workflow: **Builder only**, focus on safety.

## When to Use

- Code cleanup without behavior change
- Improving readability/performance
- Reducing technical debt

## Phases

### ~~Phase 1: GATEKEEPER~~ — SKIPPED

Requirement is clear: improve code without changing behavior.

### ~~Phase 2: ARCHITECT~~ — SKIPPED

No new design needed for refactoring.

### Phase 3: BUILDER (Refactor Focus)

#### Pre-conditions

```
⚠️ REQUIRED: Tests must exist before refactoring
   If no tests → Write tests first!
```

#### Step 1: Verify Baseline

```
1. Run all tests → MUST pass
2. Note current coverage
3. Identify refactoring targets
```

#### Step 2: Refactor Loop

```
For each change:
1. Make ONE small change
2. Run tests immediately
3. If fail → Undo
4. If pass → Commit
5. Repeat
```

#### Step 3: Verify

```
1. All tests still pass
2. Coverage not decreased
3. No new code smells
4. Behavior unchanged
```

## Refactoring Techniques

| Technique        | When                    |
| ---------------- | ----------------------- |
| Extract Function | Long function           |
| Extract Variable | Complex expression      |
| Rename           | Unclear names           |
| Inline           | Unnecessary abstraction |
| Move             | Wrong location          |

## Output Checklist

- [ ] Tests existed before (or added)
- [ ] All tests pass after
- [ ] Code is cleaner
- [ ] No behavior change
- [ ] Each change committed separately

## Example

```
Refactor: "Clean up UserService"

Before:
- 300 line function
- Magic numbers
- Duplicate code

After:
- 5 small functions
- Named constants
- Shared utilities
- Same behavior, same tests pass
```
