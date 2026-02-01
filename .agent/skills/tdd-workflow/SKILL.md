---
name: tdd-workflow
description: Test-Driven Development cycle - RED, GREEN, REFACTOR
---

# TDD Workflow Skill

## Purpose

Implement features following the TDD cycle to ensure code always has tests and works correctly.

## When to Use

- When implementing new features
- When fixing bugs
- When refactoring code

## The TDD Cycle

```
    ┌──────────────────┐
    │     1. RED       │  Write failing test
    │  (Write Test)    │
    └────────┬─────────┘
             │
             ▼
    ┌──────────────────┐
    │    2. GREEN      │  Write minimal code
    │  (Make It Pass)  │
    └────────┬─────────┘
             │
             ▼
    ┌──────────────────┐
    │   3. REFACTOR    │  Clean up
    │   (Clean Up)     │
    └────────┬─────────┘
             │
             └──────────► Repeat
```

## Phase 1: RED (Write Failing Test)

### Rules

- Write test BEFORE implementation
- Test MUST fail initially
- Only write ONE test at a time
- Test the behavior, not implementation

### Process

```
1. Write test case for one behavior
2. Run test → MUST fail
3. If test passes immediately → Test is wrong or feature already exists
```

### Example

```typescript
// Test first - this MUST fail initially
describe("UserService", () => {
  it("should create a user with hashed password", async () => {
    const result = await userService.create({
      email: "test@example.com",
      password: "plain123",
    });

    expect(result.email).toBe("test@example.com");
    expect(result.password).not.toBe("plain123"); // hashed
  });
});
```

## Phase 2: GREEN (Make It Pass)

### Rules

- Write MINIMAL code to pass the test
- Don't optimize yet
- Don't refactor yet
- It's OK to be "ugly" temporarily

### Process

```
1. Write the simplest code to pass test
2. Run test → MUST pass
3. Commit code (WIP commit)
```

### Example

```typescript
// Minimal implementation - just make it pass
class UserService {
  async create(input: CreateUserInput) {
    const hashedPassword = await hash(input.password);
    return { email: input.email, password: hashedPassword };
  }
}
```

## Phase 3: REFACTOR (Clean Up)

### Rules

- Tests MUST still pass after refactoring
- Remove duplication
- Apply clean code principles
- Extract functions/classes if needed

### Process

```
1. Identify code smells
2. Refactor incrementally
3. Run test after EACH change
4. Commit when clean
```

### What to Refactor

- [ ] Magic numbers → Constants
- [ ] Long functions → Extract
- [ ] Duplicate code → DRY
- [ ] Poor names → Rename
- [ ] Complex conditionals → Simplify

## Test Structure (AAA Pattern)

```typescript
it("should [behavior]", async () => {
  // Arrange - Setup
  const input = {
    /* ... */
  };

  // Act - Execute
  const result = await service.method(input);

  // Assert - Verify
  expect(result).toEqual(expected);
});
```

## What to Test

| Priority | Test Type                           |
| -------- | ----------------------------------- |
| High     | Happy path (main flow)              |
| High     | Edge cases (empty, null, max)       |
| High     | Error cases (validation, not found) |
| Medium   | Boundary conditions                 |
| Low      | Performance (if critical)           |

## Anti-Patterns to Avoid

❌ Writing implementation first, tests later
❌ Writing multiple tests before any implementation
❌ Testing implementation details (private methods)
❌ Not running tests after refactoring
❌ Skipping REFACTOR phase
