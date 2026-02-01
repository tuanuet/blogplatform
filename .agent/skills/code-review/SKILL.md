---
name: code-review
description: Code review checklist covering security, performance, and best practices
---

# Code Review Skill

## Purpose

Review code to detect issues before merging.

## When to Use

- Self-review before committing
- Reviewing PRs from teammates
- Final check after implementation

## Review Checklist

### 1. Correctness

- [ ] Code does what it's supposed to do
- [ ] Edge cases are handled
- [ ] Error handling is appropriate
- [ ] No off-by-one errors
- [ ] Nulls/undefineds are handled

### 2. Security

- [ ] No SQL injection vulnerabilities
- [ ] No XSS vulnerabilities
- [ ] Inputs are validated/sanitized
- [ ] Auth/Authorization is checked
- [ ] No secrets in code
- [ ] Sensitive data is not logged

### 3. Performance

- [ ] No N+1 queries
- [ ] No unnecessary loops
- [ ] Expensive operations are optimized
- [ ] Pagination for large datasets
- [ ] Proper indexing (if DB changes)

### 4. Code Quality

- [ ] Names are meaningful
- [ ] Functions are small and focused
- [ ] No code duplication
- [ ] No dead code
- [ ] Comments explain WHY

### 5. Testing

- [ ] Tests exist for new code
- [ ] Tests cover edge cases
- [ ] Tests are readable
- [ ] No flaky tests

### 6. Architecture

- [ ] Follows existing patterns
- [ ] Dependencies go in right direction
- [ ] No circular dependencies
- [ ] Appropriate abstractions

## Common Issues

### Security

```typescript
// ❌ SQL Injection
const query = `SELECT * FROM users WHERE id = ${id}`;

// ✅ Parameterized query
const query = "SELECT * FROM users WHERE id = $1";
db.query(query, [id]);
```

```typescript
// ❌ XSS vulnerable
element.innerHTML = userInput;

// ✅ Safe
element.textContent = userInput;
```

### Performance

```typescript
// ❌ N+1 query
const orders = await getOrders();
for (const order of orders) {
  order.customer = await getCustomer(order.customerId); // N queries!
}

// ✅ Batch query
const orders = await getOrders();
const customerIds = orders.map((o) => o.customerId);
const customers = await getCustomersByIds(customerIds); // 1 query
```

### Error Handling

```typescript
// ❌ Swallowing error
try {
  await riskyOperation();
} catch (e) {
  // silent fail
}

// ✅ Handle appropriately
try {
  await riskyOperation();
} catch (e) {
  logger.error("Operation failed", { error: e });
  throw new AppError("OPERATION_FAILED", e.message);
}
```

## Review Template

```markdown
## Summary

[Brief description of changes]

## Changes Reviewed

- [ ] Correctness
- [ ] Security
- [ ] Performance
- [ ] Code Quality
- [ ] Testing

## Issues Found

1. [Issue 1]
2. [Issue 2]

## Suggestions

1. [Suggestion 1]

## Verdict

- [ ] Approved
- [ ] Needs changes
```

## Giving Good Feedback

### Do

- Be specific ("Line 42: this could throw null")
- Suggest alternatives
- Explain reasoning
- Acknowledge good code

### Don't

- Be vague ("this is bad")
- Be personal ("you always...")
- Nitpick style (use linter)
- Block for minor issues
