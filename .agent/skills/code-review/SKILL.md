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
- [ ] Nulls/nils are handled

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

## Golang Specific Issues

### Context Handling

Always pass `context.Context` to long-running operations (DB calls, HTTP requests).

```go
// ❌ Bad - No timeout control
func GetUser(id string) (*User, error) {
    return db.Where("id = ?", id).First(&user).Error
}

// ✅ Good - Timeout/Cancel support
func GetUser(ctx context.Context, id string) (*User, error) {
    return db.WithContext(ctx).Where("id = ?", id).First(&user).Error
}
```

### Goroutine Leaks

Always handle goroutine termination and channel draining.

```go
// ❌ Bad - Potential leak
for _, url := range urls {
    go fetch(url) // Can spawn 1000s of goroutines
}

// ✅ Good - Concurrency limit
maxWorkers := 10
sem := make(chan struct{}, maxWorkers)

for _, url := range urls {
    sem <- struct{}{} // Block if full
    go func(u string) {
        defer func() { <-sem }()
        fetch(u)
    }(url)
}
```

### Race Conditions

Avoid concurrent writes to shared state.

```go
// ❌ Bad - Data race
var counter int

for i := 0; i < 100; i++ {
    go func() { counter++ }()
}

// ✅ Good - Mutex protection
var (
    counter int
    mu     sync.Mutex
)

func increment() {
    mu.Lock()
    counter++
    mu.Unlock()
}
```

### Panic Recovery

Recover from panics, especially in goroutines.

```go
func worker(jobs <-chan Job) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered: %v", r)
        }
    }()
    for job := range jobs {
        process(job)
    }
}
```

### Interface Satisfaction

Check if struct implements required interface methods.

```go
type Saver interface {
    Save() error
}

// ❌ Bad - Method pointer receiver for struct that implements interface
func (u *User) Save() error { ... }

type Saver interface {
    Save() error
}

// ✅ Good
func (u User) Save() error { ... }
```

### Error Wrapping

Use `fmt.Errorf` with `%w` to preserve stack traces.

```go
// ❌ Bad
return errors.New("failed to save user")

// ✅ Good
return fmt.Errorf("failed to save user: %w", err)
```

### Defer Statements

Ensure resources are released even on error.

```go
// ❌ Bad
f, _ := os.Open("file.txt")
// Do work
f.Close() // Might not run if panic

// ✅ Good
f, err := os.Open("file.txt")
if err != nil {
    return err
}
defer f.Close()
// Do work
```

### nil Pointer Checks

Validate pointers before dereferencing.

```go
// ❌ Bad
var user *User
fmt.Println(user.Name) // Panic if nil

// ✅ Good
if user != nil {
    fmt.Println(user.Name)
}
```

### Error Checks

Never ignore errors.

```go
// ❌ Bad
json.Unmarshal(data, &result)

// ✅ Good
if err := json.Unmarshal(data, &result); err != nil {
    return fmt.Errorf("invalid JSON: %w", err)
}
```

### SQL Injection (GORM)

Use parameters, never string interpolation.

```go
// ❌ Bad
db.Raw("SELECT * FROM users WHERE id = " + userID)

// ✅ Good
db.Raw("SELECT * FROM users WHERE id = ?", userID)
// Or better with GORM:
db.Where("id = ?", userID).First(&user)
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
