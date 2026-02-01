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
  } // Not UserService's job
  generateReport() {
    /* ... */
  } // Not UserService's job
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
class PaypalPayment implements PaymentMethod {
  /* ... */
}
```

### L - Liskov Substitution

> Subclass must be substitutable for parent class.

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
  private db = new PostgresDB(); // Concrete
}

// ✅ Good
class UserService {
  constructor(private db: IDatabase) {} // Abstraction
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
```

### Service Pattern

Business logic lives in services, not controllers.

```typescript
// Controller - only handles HTTP
class UserController {
  constructor(private userService: IUserService) {}

  async create(req, res) {
    const user = await this.userService.create(req.body);
    res.json(user);
  }
}

// Service - business logic
class UserService implements IUserService {
  async create(data: CreateUserInput) {
    // validation, business rules
  }
}
```

### Factory Pattern

Create complex objects.

### Strategy Pattern

Select algorithm at runtime.

### Observer Pattern

Event-driven communication.

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

Dependency Rule: Inner layers DON'T know outer layers
```

## Pattern Selection Guide

| Problem                      | Pattern    |
| ---------------------------- | ---------- |
| Need to swap implementations | Strategy   |
| Complex object creation      | Factory    |
| Decouple data access         | Repository |
| Organize business logic      | Service    |
| React to events              | Observer   |
| Add behavior dynamically     | Decorator  |
