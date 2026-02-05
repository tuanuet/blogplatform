---
name: testing
description: Testing strategies for unit test, integration test, and e2e tests
---

# Testing Skill

## Purpose

Write effective tests to ensure code quality and catch bugs early.

## When to Use

- TDD workflow (RED phase)
- Adding tests to existing features
- Verifying bug fixes

## Testing Pyramid

```
        /\
       /  \     E2E Tests (few)
      / E2E\    - Full user flows
     /------\   - Slow, expensive
    /        \
   /Integration\ Integration Tests (some)
  /            \ - Multiple units together
 /--------------\
/     Unit       \ Unit Tests (many)
/                  \ - Single function/class
/--------------------\ - Fast, cheap
```

## Test Types

### Unit Tests

Test a single function/class in isolation.

```typescript
describe("calculatePrice", () => {
  it("should return base price for quantity 1", () => {
    expect(calculatePrice(100, 1)).toBe(100);
  });

  it("should apply 10% discount for quantity > 10", () => {
    expect(calculatePrice(100, 11)).toBe(990); // 100 * 11 * 0.9
  });
});
```

### Integration Tests

Test multiple units working together.

```typescript
describe("UserService + UserRepository", () => {
  it("should create and retrieve user", async () => {
    const user = await userService.create({ email: "test@test.com" });
    const found = await userService.getById(user.id);

    expect(found.email).toBe("test@test.com");
  });
});
```

### E2E Tests

Test full user flows.

```typescript
test("user can login and view dashboard", async ({ page }) => {
  await page.goto("/login");
  await page.fill("[name=email]", "user@test.com");
  await page.fill("[name=password]", "password123");
  await page.click("button[type=submit]");

  await expect(page).toHaveURL("/dashboard");
});
```

## Test Structure (AAA)

```typescript
it("should [expected behavior]", async () => {
  // Arrange - Setup
  const user = createTestUser();
  const input = { amount: 100 };

  // Act - Execute
  const result = await service.process(user, input);

  // Assert - Verify
  expect(result.status).toBe("success");
  expect(result.balance).toBe(100);
});
```

## What to Test

| Priority | Scenario                           |
| -------- | ---------------------------------- |
| High     | Happy path (main flow works)       |
| High     | Validation errors                  |
| High     | Edge cases (empty, null, boundary) |
| High     | Error handling                     |
| Medium   | Permission/auth checks             |
| Low      | Performance (if critical)          |

## Mocking

### When to Mock

- External services (APIs, email)
- Database (for unit tests)
- Time/dates
- Random values

### When NOT to Mock

- Simple utilities
- Pure functions
- The thing you're testing

```typescript
// Mock external service
const mockEmailService = {
  send: vi.fn().mockResolvedValue(true),
};

// Inject mock
const userService = new UserService(mockEmailService);

// Assert mock was called
expect(mockEmailService.send).toHaveBeenCalledWith({
  to: "user@test.com",
  subject: "Welcome",
});
```

## Test Data

### Use Factories

```typescript
function createTestUser(overrides = {}) {
  return {
    id: "user-1",
    email: "test@example.com",
    name: "Test User",
    ...overrides,
  };
}

// Usage
const admin = createTestUser({ role: "admin" });
```

## Test Naming

```typescript
// Format: should [expected] when [condition]
it("should throw error when email is empty");
it("should return null when user not found");
it("should apply discount when quantity exceeds 10");
```

## Checklist

- [ ] Test covers happy path
- [ ] Test covers edge cases
- [ ] Test covers error cases
- [ ] Test names are descriptive
- [ ] Test is independent (no shared state)
- [ ] Test runs fast
- [ ] Minimal mocking
