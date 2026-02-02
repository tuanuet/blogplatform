---
name: design-patterns
description: Apply SOLID principles, DDD, and Clean Architecture patterns
---

# Design Patterns Skill

## Purpose

Apply design patterns and principles for clean, maintainable, testable code.

## When to Use

- When designing system architecture
- When refactoring code structure
- When reviewing code quality

## SOLID Principles

### S - Single Responsibility

> Each class/module should have only ONE reason to change.

```typescript
// ❌ Bad
class UserService {
  createUser() {
    /* ... */
  }
  sendEmail() {
    /* ... */
  }
  generateReport() {
    /* ... */
  }
}

// ✅ Good
class UserService {
  createUser() {
    /* ... */
  }
}

class EmailService {
  send() {
    /* ... */
  }
}

class ReportService {
  generate() {
    /* ... */
  }
}
```

### O - Open/Closed

> Open for extension, closed for modification.

```typescript
// ❌ Bad - must modify code when adding payment method
class PaymentProcessor {
  process(type: string) {
    if (type === "credit") {
      /* ... */
    }
    if (type === "paypal") {
      /* ... */
    }
  }
}

// ✅ Good - adding payment method = adding new class
interface PaymentMethod {
  process(): void;
}

class CreditPayment implements PaymentMethod {
  /* ... */
}

class PayPalPayment implements PaymentMethod {
  /* ... */
}
```

### L - Liskov Substitution

> Subclass must be substitutable for parent class.

```typescript
// ❌ Bad - subclass violates behavior
class Rectangle {
  width: number;
  height: number;
}

class Square extends Rectangle {
  constructor() {
    super();
  }
  getArea(): number {
    return this.width * this.width; // Square's width is side, not width
  }
}

// ✅ Good - shape hierarchy
interface Shape {
  getArea(): number;
}

class Rectangle implements Shape {
  constructor(w: number, h: number) {
    this.width = w;
    this.height = h;
  }
  getArea(): number {
    return this.width * this.height;
  }
}

class Square implements Shape {
  constructor(side: number) {
    this.side = side;
  }
  getArea(): number {
    return this.side * this.side;
  }
}
```

### I - Interface Segregation

> Clients should not depend on interfaces they don't use.

```typescript
// ❌ Bad - fat interface
interface Worker {
  work(): void;
  eat(): void;
  sleep(): void;
}

// ✅ Good - segregated interfaces
interface Workable {
  work(): void;
}

interface Eatable {
  eat(): void;
}

interface Sleepable {
  sleep(): void;
}
```

### D - Dependency Inversion

> Depend on abstractions, not concretions.

```typescript
// ❌ Bad
class UserService {
  private db = new PostgresDB(); // Concrete dependency
}

// ✅ Good
interface Database {
  findUser(id: string): User;
}

class UserService {
  constructor(private db: Database) {} // Abstract dependency
}
```

## Common Patterns

### Repository Pattern

Separate data access logic from business logic.

```typescript
interface IUserRepository {
  findById(id: string): Promise<User | null>;
  save(user: User): Promise<void>;
}

class PostgresUserRepository implements IUserRepository {
  // Implementation
}

class UserService {
  constructor(private repo: IUserRepository) {}

  async getUser(id: string) {
    return await this.repo.findById(id);
  }
}
```

### Service Pattern

Business logic lives in services, not controllers.

```typescript
class UserService {
  async create(data: CreateUserInput) {
    // validation, business rules
    return await this.repo.save(data);
  }
}
```

### Factory Pattern

Create complex objects.

```typescript
interface NotificationFactory {
  create(type: string): Notification;
}

class EmailNotification implements Notification {
  send() { /* ... */ }
}

class SMSNotification implements Notification {
  send() { /* ... */ }
}

class NotificationFactory implements NotificationFactory {
  create(type: string): Notification {
    switch (type) {
      case "email": return new EmailNotification();
      case "sms": return new SMSNotification();
    }
  }
}
```

### Strategy Pattern

Select algorithm at runtime.

```typescript
interface PaymentStrategy {
  process(amount: number): void;
}

class CreditCardPayment implements PaymentStrategy {
  process(amount: number) { /* ... */ }
}

class PayPalPayment implements PaymentStrategy {
  process(amount: number) { /* ... */ }
}

class PaymentProcessor {
  constructor(private strategy: PaymentStrategy) {}

  pay(amount: number) {
    this.strategy.process(amount);
  }
}
```

### Observer Pattern

Event-driven communication.

```typescript
interface Observer {
  update(event: Event): void;
}

interface Subject {
  subscribe(observer: Observer): void;
  unsubscribe(observer: Observer): void;
  notify(event: Event): void;
}

class EventPublisher implements Subject {
  private observers: Observer[] = [];

  subscribe(observer: Observer): void {
    this.observers.push(observer);
  }

  notify(event: Event): void {
    this.observers.forEach(o => o.update(event));
  }
}
```

### Decorator Pattern

Add behavior to objects dynamically.

```go
// Base Component
type NotificationSender interface {
    Send(msg string) error
}

// Base Implementation
type EmailSender struct{}

func (e *EmailSender) Send(msg string) error {
    return fmt.Sprintf("Sending via Email: %s", msg)
}

// Decorator: Logging
type LoggingSender struct {
    wrapped NotificationSender
}

func (l *LoggingSender) Send(msg string) error {
    fmt.Printf("[LOG] Sending message: %s\n", msg)
    return l.wrapped.Send(msg)
}

// Decorator: Rate Limiting
type RateLimitSender struct {
    wrapped   NotificationSender
    limit     time.Duration
    lastSent  time.Time
}

func (r *RateLimitSender) Send(msg string) error {
    if time.Since(r.lastSent) < r.limit {
        return fmt.Errorf("rate limit exceeded")
    }
    r.lastSent = time.Now()
    return r.wrapped.Send(msg)
}

// Usage
sender := &EmailSender{}
logSender := &LoggingSender{sender}
rateLimitSender := &RateLimitSender{wrapped: logSender, limit: time.Minute}

rateLimitSender.Send("Hello World")
```

### Adapter Pattern

Convert incompatible interfaces.

```go
// External Library (incompatible)
type ExternalClient interface {
    ProcessData(data []byte) error
}

// Our Interface (expected)
type DataProcessor interface {
    Process(data []byte) error
}

// Adapter
type ExternalAdapter struct {
    client ExternalClient
}

func (a *ExternalAdapter) Process(data []byte) error {
    // Convert data format if needed
    return a.client.ProcessData(data)
}
```

## Golang Specific Patterns

### Dependency Injection (DI)

Wire dependencies cleanly using interfaces.

**Option 1: Constructor Injection (Simple)**
```go
type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

**Option 2: Function Options (Idiomatic)**
```go
type Option func(*UserService)

func WithRepo(repo UserRepository) Option {
    return func(s *UserService) {
        s.repo = repo
    }
}

func NewUserService(opts ...Option) *UserService {
    s := &UserService{}
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

### Middleware Pattern

Cross-cutting concerns (logging, auth, recovery).

```go
// Middleware Signature
type HandlerFunc func(*gin.Context)

func LoggingMiddleware() HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        fmt.Printf("[%s] %s %s\n", start.Format(time.RFC3339), c.Request.Method, c.Request.URL.Path)

        c.Next()

        fmt.Printf("Request completed in %v\n", time.Since(start))
    }
}

func AuthMiddleware(authService AuthService) HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        user := authService.Validate(token)
        c.Set("userID", user.ID)
        c.Next()
    }
}

// Usage in Router
r := gin.Default()
r.Use(LoggingMiddleware())
r.Use(AuthMiddleware())
```

### Singleton Pattern

Ensure single instance of expensive resources.

```go
// ❌ Bad - Race condition
var db *Database

func GetDB() *Database {
    return db
}

// ✅ Good - Thread-safe using sync.Once
var (
    instance *Database
    once     sync.Once
)

func GetDB() *Database {
    once.Do(func() {
        instance = NewDatabaseConnection()
    })
    return instance
}
```

### Builder Pattern

Construct complex objects step-by-step.

```go
// ❌ Bad - Too many parameters
func NewServer(host string, port int, timeout time.Duration, maxConns int) *Server {
    return &Server{
        host:       host,
        port:       port,
        timeout:    timeout,
        maxConns:  maxConns,
    }
}

// ✅ Good - Fluent Builder
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

// Usage
s := NewServerBuilder().
    Host("localhost").
    Port(8080).
    Timeout(30 * time.Second).
    Build()
```

### Facade Pattern

Provide simplified interface to complex subsystems.

```go
type OrderSystem interface {
    ProcessPayment()
    CheckInventory()
    ShipItem()
}

type NotificationSystem interface {
    SendEmail(to string, msg string) error
    SendSMS(to string, msg string) error
}

// Facade simplifies interactions
type ShoppingFacade struct {
    orders  OrderSystem
    notify   NotificationSystem
}

func (s *ShoppingFacade) Checkout(orderID string) error {
    // Simplified interface
    if err := s.orders.CheckInventory(); err != nil {
        return err
    }
    if err := s.orders.ProcessPayment(); err != nil {
        return err
    }

    // Notify user
    s.notify.SendEmail("user@example.com", "Order confirmed")
    s.notify.SendSMS("+1234567890", "Order shipped")
    return nil
}
```

### Context Pattern

Pass request-scoped data down the call chain.

```go
// ❌ Bad - No context
func GetUser(id string) (*User, error) {
    return db.Where("id = ?", id).First(&user).Error
}

// ✅ Good - With Context
func GetUser(ctx context.Context, id string) (*User, error) {
    return db.WithContext(ctx).Where("id = ?", id).First(&user).Error
}

// Usage in Handlers
func (h *Handler) GetUser(c *gin.Context) {
    ctx := c.Request.Context()
    user, err := h.service.GetUser(ctx, c.Param("id"))
    // ...
}
```

## Clean Architecture Layers

```
┌─────────────────────────────────────┐
│           Presentation              │  Controllers, Routes
├─────────────────────────────────────┤
│           Application               │  Use Cases, Services
├─────────────────────────────────────┤
│             Domain                  │  Entities, Business Rules
├─────────────────────────────────────┤
│          Infrastructure             │  DB, External APIs
└─────────────────────────────────────┘
```

**Dependency Rule**: Inner layers DON'T know outer layers

## Pattern Selection Guide

| Problem                      | Pattern    |
| ---------------------------- | ---------- |
| Need to swap implementations | Strategy   |
| Complex object creation      | Factory    |
| Decouple data access         | Repository |
| Organize business logic      | Service    |
| React to events              | Observer   |
| Add behavior dynamically     | Decorator  |
| Add behavior dynamically     | Facade     |
| Simplify complex subsystems   | Facade     |

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
