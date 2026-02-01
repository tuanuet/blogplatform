---
name: api-contract
description: Design API contracts with OpenAPI and interface definitions
---

# API Contract Skill

## Purpose

Design API contracts (endpoints, request/response schemas) before implementation.

## When to Use

- When designing REST/GraphQL APIs
- When defining service interfaces
- When documenting APIs for consumers

## Design Principles

### RESTful Guidelines

- Use nouns, not verbs (`/users` not `/getUsers`)
- Use HTTP methods correctly (GET, POST, PUT, PATCH, DELETE)
- Use proper status codes (200, 201, 400, 401, 404, 500)
- Version your API (`/api/v1/...`)

### Naming Conventions

- Endpoints: kebab-case (`/user-profiles`)
- Query params: camelCase (`?pageSize=10`)
- Request/Response: camelCase keys

## Contract Templates

### OpenAPI YAML

```yaml
openapi: 3.0.3
info:
  title: [API Name]
  version: 1.0.0

paths:
  /api/v1/[resource]:
    get:
      summary: List [resources]
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            default: 1
      responses:
        200:
          description: Success
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/[Resource]ListResponse"

    post:
      summary: Create [resource]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/Create[Resource]Request"
      responses:
        201:
          description: Created
        400:
          description: Validation error

  /api/v1/[resource]/{id}:
    get:
      summary: Get [resource] by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        200:
          description: Success
        404:
          description: Not found

components:
  schemas:
    [Resource]:
      type: object
      properties:
        id:
          type: string
          format: uuid
        # fields...
        createdAt:
          type: string
          format: date-time
```

### TypeScript Interface

```typescript
// DTOs
interface Create[Resource]Request {
  // input fields (no id, no timestamps)
}

interface Update[Resource]Request {
  // partial input fields
}

interface [Resource]Response {
  id: string;
  // output fields
  createdAt: Date;
  updatedAt: Date;
}

interface [Resource]ListResponse {
  data: [Resource]Response[];
  pagination: {
    page: number;
    pageSize: number;
    total: number;
  };
}

// Service Interface (Contract only, NO implementation)
interface I[Resource]Service {
  create(input: Create[Resource]Request): Promise<[Resource]Response>;
  getById(id: string): Promise<[Resource]Response | null>;
  list(params: ListParams): Promise<[Resource]ListResponse>;
  update(id: string, input: Update[Resource]Request): Promise<[Resource]Response>;
  delete(id: string): Promise<void>;
}

// Repository Interface
interface I[Resource]Repository {
  insert(data: [Resource]): Promise<[Resource]>;
  findById(id: string): Promise<[Resource] | null>;
  findMany(query: QueryParams): Promise<[Resource][]>;
  update(id: string, data: Partial<[Resource]>): Promise<[Resource]>;
  delete(id: string): Promise<void>;
}
```

## HTTP Status Codes

| Code               | When to Use                    |
| ------------------ | ------------------------------ |
| 200 OK             | GET success, PUT/PATCH success |
| 201 Created        | POST success                   |
| 204 No Content     | DELETE success                 |
| 400 Bad Request    | Validation error               |
| 401 Unauthorized   | Not authenticated              |
| 403 Forbidden      | Not authorized                 |
| 404 Not Found      | Resource not found             |
| 409 Conflict       | Duplicate/conflict             |
| 500 Internal Error | Server error                   |

## Validation Checklist

- [ ] All endpoints have clear request/response schemas
- [ ] Error responses are defined
- [ ] Authentication requirements are specified
- [ ] Pagination is consistent
- [ ] No implementation details (contracts only)
