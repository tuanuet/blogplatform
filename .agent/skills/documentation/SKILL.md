---
name: documentation
description: Writing clear documentation for code, APIs, and processes
---

# Documentation Skill

## Purpose

Write clear, useful documentation for both humans and AI agents.

## When to Use

- Creating Refined Spec (Gatekeeper)
- Documenting API contracts (Architect)
- Code comments and READMEs (Builder)
- Changelogs and release notes

## Documentation Types

### 1. Refined Spec

Output of Gatekeeper Agent.

```markdown
# Refined Spec: [Feature Name]

## User Story

As a [role], I want to [action] so that [benefit].

## Acceptance Criteria

- [ ] Given [context], When [action], Then [result]
- [ ] Given [context], When [action], Then [result]

## Edge Cases

1. [What happens when X]
2. [What happens when Y]

## Out of Scope

- [What this does NOT include]

## Tech Stack

- Language: [detected]
- Database: [detected]
- Testing: [detected]
```

### 2. API Documentation

Output of Architect Agent.

```markdown
# API: [Endpoint Name]

## Endpoint

`POST /api/v1/users`

## Description

Creates a new user account.

## Request

\`\`\`json
{
"email": "string (required)",
"password": "string (required, min 8 chars)",
"name": "string (optional)"
}
\`\`\`

## Response

### Success (201)

\`\`\`json
{
"id": "uuid",
"email": "string",
"name": "string",
"createdAt": "datetime"
}
\`\`\`

### Error (400)

\`\`\`json
{
"error": "VALIDATION_ERROR",
"message": "Email is required"
}
\`\`\`
```

### 3. Code Comments

```typescript
/**
 * Calculates the total order price with applicable discounts.
 *
 * Discount rules:
 * - 10% off for orders > $100
 * - 20% off for VIP customers
 * - Discounts are NOT stackable
 *
 * @param order - The order to calculate
 * @param customer - The customer placing the order
 * @returns Total price after discounts
 */
function calculateTotalPrice(order: Order, customer: Customer): number {
  // ...
}
```

### 4. README

```markdown
# Project Name

Brief description.

## Quick Start

\`\`\`bash
npm install
npm run dev
\`\`\`

## Project Structure

\`\`\`
src/
├── modules/ # Feature modules
├── shared/ # Shared utilities
└── config/ # Configuration
\`\`\`

## Development

### Prerequisites

- Node.js 20+
- PostgreSQL 15+

### Running Tests

\`\`\`bash
npm test
\`\`\`

## API Reference

See [API docs](./docs/api.md)
```

### 5. Changelog

```markdown
# Changelog

## [1.2.0] - 2026-01-31

### Added

- User password reset feature (#123)

### Changed

- Improved login performance (#124)

### Fixed

- Email validation bug (#125)

### Security

- Updated dependencies with vulnerabilities
```

## Writing Guidelines

### Be Concise

```markdown
# ❌ Too verbose

This function is responsible for taking an input parameter
which represents the user's email address and then performing
validation on that email address to ensure it conforms to
the standard email format as defined by RFC 5322.

# ✅ Concise

Validates email format per RFC 5322.
```

### Use Examples

```markdown
# ❌ Abstract

The function accepts configuration options.

# ✅ With example

\`\`\`typescript
createUser({
email: 'user@example.com',
role: 'admin'
});
\`\`\`
```

### Keep Updated

- Update docs when code changes
- Mark outdated sections clearly
- Review docs in PRs

## Checklist

- [ ] Purpose is clear
- [ ] Examples are provided
- [ ] Edge cases documented
- [ ] Up to date with code
- [ ] Formatted consistently
