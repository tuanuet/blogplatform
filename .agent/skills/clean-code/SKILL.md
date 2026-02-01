---
name: clean-code
description: Writing readable, maintainable code following clean code principles
---

# Clean Code Skill

## Purpose

Write clean, readable, maintainable code following industry standards.

## When to Use

- During REFACTOR phase of TDD
- When reviewing code
- When refactoring legacy code

## Naming Conventions

### Variables

```typescript
// ❌ Bad
const d = new Date();
const arr = getUsers();
const flag = true;

// ✅ Good
const currentDate = new Date();
const activeUsers = getUsers();
const isAuthenticated = true;
```

### Functions

```typescript
// ❌ Bad
function process() {}
function handle() {}
function doIt() {}

// ✅ Good
function calculateTotalPrice() {}
function validateUserInput() {}
function sendVerificationEmail() {}
```

### Booleans

```typescript
// ❌ Bad
const active = true;
const login = false;

// ✅ Good
const isActive = true;
const hasLoggedIn = false;
const canEdit = true;
const shouldRefresh = false;
```

## Functions

### Single Responsibility

```typescript
// ❌ Bad - does too many things
function processOrder(order) {
  validateOrder(order);
  calculateTotal(order);
  saveToDatabase(order);
  sendEmail(order);
}

// ✅ Good - one thing
function validateOrder(order) {
  /* ... */
}
function calculateTotal(order) {
  /* ... */
}
function saveOrder(order) {
  /* ... */
}
function notifyOrderCreated(order) {
  /* ... */
}
```

### Keep Functions Small

- Max 20-30 lines
- Max 3 parameters
- Max 2-3 levels of nesting

### Avoid Side Effects

```typescript
// ❌ Bad - hidden side effect
function getUser(id) {
  logAccess(id); // Side effect!
  return users.find((u) => u.id === id);
}

// ✅ Good - no side effects
function getUser(id) {
  return users.find((u) => u.id === id);
}
```

## Comments

### When to Comment

- Explain WHY, not WHAT
- Document public APIs
- Explain complex algorithms
- Mark TODOs with context

### When NOT to Comment

```typescript
// ❌ Bad - comment explains WHAT
// Increment i by 1
i++;

// ✅ Good - self-explanatory code needs no comment
i++;
```

## Code Smells

| Smell               | Fix                         |
| ------------------- | --------------------------- |
| Long function       | Extract smaller functions   |
| Long parameter list | Use object parameter        |
| Magic numbers       | Use named constants         |
| Duplicate code      | Extract to shared function  |
| Deep nesting        | Early return, guard clauses |
| Dead code           | Delete it                   |
| Large class         | Split into smaller classes  |

## Early Return Pattern

```typescript
// ❌ Bad - deep nesting
function process(user) {
  if (user) {
    if (user.isActive) {
      if (user.hasPermission) {
        // do something
      }
    }
  }
}

// ✅ Good - early returns
function process(user) {
  if (!user) return;
  if (!user.isActive) return;
  if (!user.hasPermission) return;

  // do something
}
```

## Constants

```typescript
// ❌ Bad
if (status === 1) {
  /* ... */
}
if (age >= 18) {
  /* ... */
}

// ✅ Good
const STATUS_ACTIVE = 1;
const MINIMUM_AGE = 18;

if (status === STATUS_ACTIVE) {
  /* ... */
}
if (age >= MINIMUM_AGE) {
  /* ... */
}
```

## Checklist

- [ ] Names reveal intent
- [ ] Functions do one thing
- [ ] No magic numbers
- [ ] No deep nesting (max 2-3 levels)
- [ ] No duplicate code
- [ ] No dead code
- [ ] Comments explain WHY
