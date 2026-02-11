---
name: consolidated-review
description: High-density, multi-phase quality gate checklist (token-optimized)
---

# Consolidated Review Checklist

## Context Analysis
- Use `serena_find_referencing_symbols` to verify cross-component impact.
- Use `serena_search_for_pattern` to check for similar existing patterns.
- Focus on `git diff` or changed files only to minimize token load.

## Phase 1: Architecture Review
- **Contracts**: Interfaces are minimal, follow Go idioms, and avoid leaked implementation details.
- **SOLID**: Proper abstraction; No circular dependencies between `domain`, `application`, and `infrastructure`.
- **Plan**: Logical task order, covers all requirements, includes verification steps.

## Phase 2: Implementation Review (Code Quality)
- **TDD**: Verification that tests were written and fail before implementation.
- **Clean Code**: No magic numbers, small functions (<30 lines), meaningful names (Go: concise for locals).
- **Security**: Parameterized queries (GORM: `?`), input validation, no secrets in logs.
- **Performance**: No N+1 queries (GORM: `Preload`), appropriate indexing, efficient loops.
- **Errors**: Explicit handling, wrapping with `%w`, context propagation (`ctx`).

## Phase 3: Integration Review
- **Wiring**: Components are correctly injected; Middleware is applied; Routers are registered.
- **Coverage**: Integration tests cover happy path + major edge cases + error scenarios.
- **Spec**: Matches the original requirements from `/brainstorm` or `/develop`.

## Verdict Format
- **APPROVED** [Phase X] - [Optional brief praise/suggestion]
- **NEEDS_CHANGES** [Phase X] - [Numbered list of critical issues with File:Line]
