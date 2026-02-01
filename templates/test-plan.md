# Test Plan: [Feature Name]

## Overview

[Brief description of testing strategy]

---

## Unit Tests

### [Service/Function Name]

| Test Case                | Input            | Expected Output       | Priority |
| ------------------------ | ---------------- | --------------------- | -------- |
| Happy path - create      | Valid data       | Success response      | High     |
| Validation - empty field | Missing required | Throw ValidationError | High     |
| Edge case - max length   | 256 char string  | Throw ValidationError | Medium   |
| Not found                | Invalid ID       | Return null           | High     |

### Test Code Template (TypeScript/Vitest)

```typescript
import { describe, it, expect, beforeEach } from 'vitest';
import { [Service] } from './[service]';

describe('[Service]', () => {
  let service: [Service];

  beforeEach(() => {
    service = new [Service]();
  });

  describe('create', () => {
    it('should create with valid input', async () => {
      // Arrange
      const input = { field1: 'value' };

      // Act
      const result = await service.create(input);

      // Assert
      expect(result.id).toBeDefined();
      expect(result.field1).toBe('value');
    });

    it('should throw when field1 is empty', async () => {
      // Arrange
      const input = { field1: '' };

      // Act & Assert
      await expect(service.create(input))
        .rejects.toThrow('field1 is required');
    });
  });

  describe('getById', () => {
    it('should return resource when exists', async () => {
      // Arrange
      const created = await service.create({ field1: 'test' });

      // Act
      const result = await service.getById(created.id);

      // Assert
      expect(result).not.toBeNull();
      expect(result?.id).toBe(created.id);
    });

    it('should return null when not found', async () => {
      // Act
      const result = await service.getById('non-existent-id');

      // Assert
      expect(result).toBeNull();
    });
  });
});
```

---

## Integration Tests

| Test Case           | Components   | Setup     | Verification             |
| ------------------- | ------------ | --------- | ------------------------ |
| Create and retrieve | Service + DB | Clean DB  | Data persisted correctly |
| Update existing     | Service + DB | Seed data | Changes reflected        |

---

## E2E Tests (if applicable)

| User Flow      | Steps                     | Expected            |
| -------------- | ------------------------- | ------------------- |
| Create via API | POST /api/v1/resource     | 201 + response body |
| Get via API    | GET /api/v1/resource/{id} | 200 + correct data  |

---

## Test Coverage Goals

| Type      | Target |
| --------- | ------ |
| Lines     | 80%+   |
| Branches  | 75%+   |
| Functions | 90%+   |

---

## Notes

- [Special testing considerations]
- [Mocking requirements]
- [Database seeding needs]
