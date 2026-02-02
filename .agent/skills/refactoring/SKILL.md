---
name: refactoring
description: Safe refactoring techniques to improve code without changing behavior
---

# Refactoring Skill

## Purpose

Improve code structure without changing behavior.

## When to Use

- REFACTOR phase of TDD
- When code smells are detected
- Before adding new features to legacy code

## Golden Rule

> Tests MUST pass before and after refactoring.
> If no tests exist → Write tests before refactoring.

## Safe Refactoring Process

```
1. Ensure tests exist and pass
2. Make ONE small change
3. Run tests
4. If fail → Undo and try again
5. If pass → Commit
6. Repeat
```

## Common Refactorings

### Extract Function

Separate code into its own function.

```typescript
// Before
function processOrder(order) {
  // validate
  if (!order.items.length) throw new Error("Empty");
  if (!order.customer) throw new Error("No customer");

  // calculate
  let total = 0;
  for (const item of order.items) {
    total += item.price * item.quantity;
  }
  return total;
}

// After
function processOrder(order) {
  validateOrder(order);
  return calculateTotal(order);
}

function validateOrder(order) {
  if (!order.items.length) throw new Error("Empty");
  if (!order.customer) throw new Error("No customer");
}

function calculateTotal(order) {
  return order.items.reduce((sum, item) => sum + item.price * item.quantity, 0);
}
```

### Extract Variable

Name complex expressions.

```typescript
// Before
if (user.age >= 18 && user.country === "VN" && !user.isBanned) {
  allowAccess();
}

// After
const isAdult = user.age >= 18;
const isVietnamese = user.country === "VN";
const isNotBanned = !user.isBanned;
const canAccess = isAdult && isVietnamese && isNotBanned;

if (canAccess) {
  allowAccess();
}
```

### Inline Variable

Remove unnecessary variables.

```typescript
// Before
const basePrice = order.basePrice();
return basePrice;

// After
return order.basePrice();
```

### Rename

Use more meaningful names.

```typescript
// Before
const d = new Date();
const arr = getUsers();

// After
const currentDate = new Date();
const activeUsers = getUsers();
```

### Replace Conditional with Polymorphism

```typescript
// Before
function getSpeed(vehicle) {
  switch (vehicle.type) {
    case "car":
      return vehicle.baseSpeed * 1.2;
    case "bike":
      return vehicle.baseSpeed * 0.8;
    default:
      return vehicle.baseSpeed;
  }
}

// After
type Vehicle interface {
  getSpeed() float64;
}

type Car struct {
  baseSpeed float64;
}

func (c Car) getSpeed() float64 {
  return c.baseSpeed * 1.2;
}

type Bike struct {
  baseSpeed float64;
}

func (b Bike) getSpeed() float64 {
  return b.baseSpeed * 0.8;
}
```

### Replace Magic Number with Constant

```typescript
// Before
if (password.length < 8) {
  /* ... */
}
if (age >= 18) {
  /* ... */
}

// After
const MIN_PASSWORD_LENGTH = 8;
const MINIMUM_AGE = 18;

if (password.length < MIN_PASSWORD_LENGTH) {
  /* ... */
}
if (age >= MINIMUM_AGE) {
  /* ... */
}
```

## Golang Specific Refactorings

### Error Wrapping

Use `fmt.Errorf` with `%w` to preserve stack traces.

```go
// ❌ Bad
return errors.New("failed to save user")

// ✅ Good
return fmt.Errorf("failed to save user: %w", err)
```

### Context Handling

Always pass `context.Context` to long-running operations.

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

### Interface Satisfaction

Extract interfaces to satisfy dependency inversion.

```go
// Before - Direct dependency
type OrderService struct {
    emailer Emailer
    saver   DatabaseSaver
}

func (s *OrderService) Process(order *Order) {
    // Uses s.emailer and s.saver directly
}

// After - Interface dependency
type NotificationService interface {
    Send(to string, msg string) error
}

type DataRepository interface {
    Save(entity interface{}) error
}

type OrderService struct {
    notifier NotificationService
    repo     DataRepository
}

func (s *OrderService) Process(order *Order) {
    s.notifier.Send(order.User.Email, "Order confirmed")
    s.repo.Save(order)
}
```

### Defer for Cleanup

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

### Reduce Allocations

Use slice tricks or `sync.Pool` for high-frequency allocations.

```go
// ❌ Bad - Allocates every iteration
func joinStrings(strings []string) string {
    var result string
    for _, s := range strings {
        result += s
    }
    return result
}

// ✅ Good - Pre-allocates
func joinStrings(strings []string) string {
    var builder strings.Builder
    builder.Grow(len(strings) * 10)
    for _, s := range strings {
        builder.WriteString(s)
    }
    return builder.String()
}
```

### Goroutine Management

Limit goroutines to prevent resource exhaustion.

```go
// ❌ Bad - Unbounded goroutines
func ProcessItems(items []Item) {
    for _, item := range items {
        go handle(item) // Can spawn 1000s
    }
}

// ✅ Good - Bounded concurrency
func ProcessItems(items []Item) {
    maxWorkers := runtime.NumCPU() * 2
    sem := make(chan struct{}, maxWorkers)
    var wg sync.WaitGroup

    for _, item := range items {
        sem <- struct{}{}
        wg.Add(1)
        go func(i Item) {
            defer wg.Done()
            defer func() { <-sem }()
            handle(i)
        }(item)
    }
    wg.Wait()
}
```

## Design Patterns (Go Edition)

### Strategy Pattern

```go
type PaymentStrategy interface {
    Process(amount int) error
}

type CreditCardPayment struct{}

func (c *CreditCardPayment) Process(amount int) error {
    // Process credit card logic
}

type PayPalPayment struct{}

func (p *PayPalPayment) Process(amount int) error {
    // Process PayPal logic
}

type PaymentProcessor struct {
    strategy PaymentStrategy
}

func (p *PaymentProcessor) Pay(amount int) error {
    return p.strategy.Process(amount)
}
```

### Factory Pattern

```go
// ❌ Bad - Tight coupling to concrete types
func GetUser(userType string) (interface{}, error) {
    if userType == "admin" {
        return &Admin{}, nil
    }
    return &User{}, nil
}

// ✅ Good - Factory function
type UserFactory interface {
    CreateUser() (interface{}, error)
}

type AdminFactory struct{}
func (f *AdminFactory) CreateUser() (interface{}, error) {
    return &Admin{}, nil
}

func GetFactory(userType string) UserFactory {
    switch userType {
    case "admin":
        return &AdminFactory{}
    case "user":
        return &UserFactory{}
    }
}
```

### Builder Pattern

```go
// ❌ Bad - Many parameters
type Server struct{}

func NewServer(host string, port int, timeout time.Duration, maxConns int) *Server {
    return &Server{
        host:       host,
        port:       port,
        timeout:    timeout,
        maxConns:  maxConns,
    }
}

// ✅ Good - Builder pattern
type ServerBuilder struct {
    config ServerConfig
}

func (b *ServerBuilder) Host(host string) *ServerBuilder {
    b.config.Host = host
    return b
}

func (b *ServerBuilder) Port(port int) *ServerBuilder {
    b.config.Port = port
    return b
}

func (b *ServerBuilder) Build() *Server {
    return &Server{
        host:   b.config.Host,
        port:   b.config.Port,
        timeout: b.config.Timeout,
    }
}

s := NewServerBuilder().
    Host("localhost").
    Port(8080).
    Timeout(30 * time.Second).
    Build()
```

## Test Refactorings

### Table-Driven Tests

Replace repeated test code with data tables.

```go
// ❌ Bad - Duplicated test logic
func TestValidateEmail(t *testing.T) {
    err := validateEmail("test@test.com")
    assert.Nil(t, err)

    err = validateEmail("invalid")
    assert.NotNil(t, err)
    // ... 10 more variations
}

// ✅ Good - Table-driven
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name  string
        email string
        want  error
    }{
        {"valid", "test@test.com", nil},
        {"invalid", "invalid", ErrInvalidEmail},
        {"empty", "", ErrEmptyEmail},
        {"too-long", strings.Repeat("a", 300), ErrTooLong},
        // ... 7 more tests
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := validateEmail(tt.email)
            if got != tt.want {
                t.Errorf("validateEmail(%q) = %v, want %v", tt.email, got)
            }
        })
    }
}
```

## When to Refactor

| Signal                | Action                        |
| --------------------- | ----------------------------- |
| Before adding feature | Clean area first          |
| After making it work  | "Make it work, make it right" |
| During code review    | Small improvements            |
| Duplicate code found  | Extract shared function       |
| Hard to understand    | Rename, extract, simplify     |

## When NOT to Refactor

- ❌ No tests exist (write tests first)
- ❌ Deadline tomorrow (risky)
- ❌ Code will be deleted soon
- ❌ Just for "beauty" (must have benefit)

## Checklist

- [ ] Tests pass before refactoring
- [ ] Make small, incremental changes
- [ ] Run tests after each change
- [ ] Commit after each successful change
- [ ] Tests still pass after refactoring
- [ ] No behavior changes
