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
class Car {
  getSpeed() {
    return this.baseSpeed * 1.2;
  }
}
class Bike {
  getSpeed() {
    return this.baseSpeed * 0.8;
  }
}
```

### Replace Magic Number with Constant

```typescript
// Before
if (password.length < 8) {
  /* ... */
}

// After
const MIN_PASSWORD_LENGTH = 8;
if (password.length < MIN_PASSWORD_LENGTH) {
  /* ... */
}
```

## When to Refactor

| Signal                | Action                        |
| --------------------- | ----------------------------- |
| Before adding feature | Clean the area first          |
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
