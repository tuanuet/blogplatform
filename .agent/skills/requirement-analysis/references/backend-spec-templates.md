# Backend Spec Templates (Go API Repo)

Use these as defaults when the request impacts HTTP APIs, data models, or migrations.

## Feature Spec (Brainstorm Output)

```markdown
# Feature: [Feature Name]

## Objective
[1-2 sentences: what we're building and why]

## Problem Statement
[Who has the problem and what they experience today]

## Requirements

### Functional
- [ ] ...

### Non-Functional
- [ ] Performance: [latency/QPS targets if known]
- [ ] Security: [authn/authz rules]
- [ ] Observability: [logs/traces/metrics expectations]
- [ ] Reliability: [idempotency/retries/timeouts]
- [ ] Migration safety: [backfill/rollout notes if data changes]

## API Sketch

### Endpoint: [Name]
- Method: `GET|POST|PUT|PATCH|DELETE`
- Path: `/api/v1/...`
- Auth: `public` | `session` | `rbac:role_name`
- Idempotency: `none` | `Idempotency-Key` (scope: user/global)

Request Example:
```json
{}
```

Response Example:
```json
{}
```

Errors:
- `400 VALIDATION_ERROR` when ...
- `401 UNAUTHENTICATED` when ...
- `403 FORBIDDEN` when ...
- `404 NOT_FOUND` when ...
- `409 CONFLICT` when ...

## User Stories

### Story 1: [Name]
As a [role], I want to [action] so that [benefit].

Acceptance Criteria:
- [ ] Given [context], When [action], Then [result]

## User Flows

### Happy Path
1. ...

### Edge/Error Cases
1. ...

## Technical Context
- Impacted modules: ...
- Data changes: ...
- Dependencies: ...

## Out of Scope
- ...

## Open Questions
- ...
```

## Refined Spec (Gatekeeper Output)

```markdown
# Refined Spec: [Feature Name]

## User Story
As a [role], I want to [action] so that [benefit].

## Acceptance Criteria
- [ ] Given [context], When [action], Then [result]

## Edge Cases
1. ...

## Tech Stack (confirm from codebase)
- Language: ...
- Framework: ...
- Database: ...

## Affected Modules
- [module/path] - why affected

## API (if applicable)

### Endpoint: [Name]
- Method: `...`
- Path: `...`
- Auth: `public|session|rbac:role_name`
- Idempotency: `none|Idempotency-Key`

Request Example:
```json
{}
```

Response Example:
```json
{}
```

Errors:
- `400 VALIDATION_ERROR` when ...
- `401 UNAUTHENTICATED` when ...
- `403 FORBIDDEN` when ...
- `404 NOT_FOUND` when ...
- `409 CONFLICT` when ...

## Data (if applicable)
- Schema changes: ...
- Migration/backfill notes: ...

## Out of Scope
- ...

## Open Questions
- ...
```

## Definition of Ready (Backend)

- Endpoints: method + path + auth requirements + request/response examples.
- Data changes: tables/columns, constraints, indexes, migration/backfill notes.
- Error model: stable error codes and when each applies.
- Validation: required fields, formats, bounds.
- Non-functional: performance expectations, rate limits (if any), observability expectations.
- Testability: what can be unit-tested vs needs integration tests.
