# API Contract: [Feature Name]

## Overview

[Brief description of the API endpoints]

---

## Endpoints

### Create [Resource]

```
POST /api/v1/[resource]
```

**Request Body:**

```json
{
  "field1": "string (required)",
  "field2": "number (optional)"
}
```

**Response (201 Created):**

```json
{
  "id": "uuid",
  "field1": "string",
  "field2": "number",
  "createdAt": "2026-01-31T00:00:00Z",
  "updatedAt": "2026-01-31T00:00:00Z"
}
```

**Errors:**
| Status | Error Code | Description |
|--------|------------|-------------|
| 400 | VALIDATION_ERROR | Invalid input |
| 401 | UNAUTHORIZED | Not authenticated |

---

### Get [Resource] by ID

```
GET /api/v1/[resource]/{id}
```

**Response (200 OK):**

```json
{
  "id": "uuid",
  "field1": "string",
  "createdAt": "2026-01-31T00:00:00Z"
}
```

**Errors:**
| Status | Error Code | Description |
|--------|------------|-------------|
| 404 | NOT_FOUND | Resource not found |

---

### List [Resources]

```
GET /api/v1/[resource]?page=1&pageSize=10
```

**Response (200 OK):**

```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "pageSize": 10,
    "total": 100
  }
}
```

---

## TypeScript Interfaces

```typescript
// DTOs
interface Create[Resource]Request {
  field1: string;
  field2?: number;
}

interface Update[Resource]Request {
  field1?: string;
  field2?: number;
}

interface [Resource]Response {
  id: string;
  field1: string;
  field2: number | null;
  createdAt: Date;
  updatedAt: Date;
}

// Service Interface (NO implementation)
interface I[Resource]Service {
  create(input: Create[Resource]Request): Promise<[Resource]Response>;
  getById(id: string): Promise<[Resource]Response | null>;
  update(id: string, input: Update[Resource]Request): Promise<[Resource]Response>;
  delete(id: string): Promise<void>;
}

// Repository Interface
interface I[Resource]Repository {
  insert(data: [Resource]): Promise<[Resource]>;
  findById(id: string): Promise<[Resource] | null>;
  update(id: string, data: Partial<[Resource]>): Promise<[Resource]>;
  delete(id: string): Promise<void>;
}
```
