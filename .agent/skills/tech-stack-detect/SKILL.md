---
name: tech-stack-detect
description: Auto-detect project tech stack from codebase files
---

# Tech Stack Detection Skill

## Purpose

Automatically detect project tech stack from files in the codebase.

## When to Use

- When starting to analyze a new project
- When choosing the correct template/format for output
- When determining which testing framework to use

## Detection Rules

### Language Detection

| File                                  | Language              |
| ------------------------------------- | --------------------- |
| `package.json`                        | JavaScript/TypeScript |
| `go.mod`                              | Go                    |
| `requirements.txt` / `pyproject.toml` | Python                |
| `Cargo.toml`                          | Rust                  |
| `pom.xml` / `build.gradle`            | Java                  |

### Framework Detection

| Indicator                     | Framework  |
| ----------------------------- | ---------- |
| `next` in package.json        | Next.js    |
| `express` in package.json     | Express.js |
| `fastify` in package.json     | Fastify    |
| `hono` in package.json        | Hono       |
| `gin-gonic` in go.mod         | Gin        |
| `echo` in go.mod              | Echo       |
| `fastapi` in requirements.txt | FastAPI    |
| `django` in requirements.txt  | Django     |

### Database Detection

| Indicator                            | Database/ORM |
| ------------------------------------ | ------------ |
| `prisma/schema.prisma`               | Prisma       |
| `drizzle.config.ts`                  | Drizzle      |
| `typeorm` in package.json            | TypeORM      |
| `gorm` in go.mod                     | GORM         |
| `sqlalchemy` in requirements.txt     | SQLAlchemy   |
| `docker-compose.yml` with `postgres` | PostgreSQL   |
| `docker-compose.yml` with `mysql`    | MySQL        |

### Testing Framework Detection

| Indicator                    | Testing Framework |
| ---------------------------- | ----------------- |
| `vitest` in package.json     | Vitest            |
| `jest` in package.json       | Jest              |
| `mocha` in package.json      | Mocha             |
| `*_test.go` files            | Go testing        |
| `pytest` in requirements.txt | Pytest            |
| `unittest` imports           | Python unittest   |

## Process

```
1. Scan root directory for indicator files
2. Parse package managers (package.json, go.mod, etc.)
3. Check for config files (prisma, drizzle, etc.)
4. Check docker-compose for infrastructure
5. Return detected stack
```

## Output Format

```json
{
  "language": "TypeScript",
  "runtime": "Node.js 20",
  "framework": "Express",
  "database": {
    "type": "PostgreSQL",
    "orm": "Prisma"
  },
  "testing": {
    "unit": "Vitest",
    "e2e": "Playwright"
  },
  "infrastructure": {
    "docker": true,
    "ci": "GitHub Actions"
  }
}
```

## Fallback

If unable to detect → Ask the user about tech stack.
