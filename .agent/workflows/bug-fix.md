---
description: Bug fix with regression test (skips Architect phase)
---

# Bug Fix Workflow

User Request: $1

Simplified workflow: **Skip Architect**, focus on reproducing and fixing.

## When to Use

- Fixing reported bugs
- No schema/API changes needed
- Quick turnaround required

## Phases

### Phase 1: GATEKEEPER (Lite)

Gather bug info:

- [ ] Steps to reproduce?
- [ ] Expected behavior?
- [ ] Actual behavior?
- [ ] Error logs/screenshots?

### ~~Phase 2: ARCHITECT~~ — SKIPPED

No design needed for bug fixes.

### Phase 3: BUILDER (Bug Focus)

#### Step 1: Reproduce

```
1. Create failing test that reproduces bug
2. Run test → MUST fail
3. This becomes your regression test
```

#### Step 2: Fix

```
1. Write minimal fix
2. Run new test → should pass
3. Run ALL tests → should pass
```

#### Step 3: Verify

```
1. Manual verification if possible
2. Check for side effects
3. Review fix for code smells
```

## Output Checklist

- [ ] Bug reproduced in test
- [ ] Root cause identified
- [ ] Fix implemented
- [ ] Regression test passing
- [ ] All other tests passing
- [ ] No new code smells

## Example

```
Bug: "Login fails for emails with + character"

Phase 1 (Gatekeeper):
- Reproduce: Login with "user+test@email.com"
- Expected: Login succeeds
- Actual: Returns 400 error

Phase 3 (Builder):
- Test: login.test.ts → "should accept + in email"
- Fix: Update email validation regex
- Verify: All auth tests pass
```
