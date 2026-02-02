## Skill: task-breakdown

**Base directory**: .agent/skills/task-breakdown

# Task Breakdown Skill

## Purpose

Decompose high-level requirements into atomic, technical implementation tasks based on the existing system architecture. This skill ensures that development work is strictly aligned with the architectural design and schema.

## Inputs

1.  **Requirements**: Feature specification, user stories, or PRD.
2.  **Architecture Context**: Database schema, API contracts, and system design documents.

## Process

### 1. Context Loading
Before generating tasks, you MUST load and analyze the architecture:
-   **Read Schema**: Understand existing data models.
-   **Read API Contracts**: Understand endpoints and interfaces.
-   **Read Existing Code**: Use `glob`/`read` to see where new code fits.

### 2. Dependency Analysis
Map requirements to architectural components:
-   **Data Layer**: Migrations, Models, Repositories.
-   **Service Layer**: Business Logic, External Integrations.
-   **Interface Layer**: API Controllers, UI Components.
-   **Validation**: Tests (Unit/Integration).

### 3. Execution Strategy (The "Sandwich" Method)
Order tasks to minimize blocking dependencies:
1.  **Foundation**: Database migrations & Domain models.
2.  **Core**: Service logic & Unit tests.
3.  **Exposure**: API endpoints / UI & Integration tests.

## Output Format

Generate a hierarchical task list. Use the `todowrite` tool to persist these tasks if working in an active session.

```markdown
## Implementation Plan: [Feature Name]

### Phase 1: Data & Domain
- [ ] **DB**: Create migration for table `[table_name]`
- [ ] **Model**: Implement `[ModelName]` struct/class with validation
- [ ] **Test**: Write unit tests for `[ModelName]`

### Phase 2: Core Logic
- [ ] **Service**: Implement `[ServiceName].method()`
- [ ] **Refactor**: Update `[ExistingComponent]` to support new flow
- [ ] **Test**: Add service-level tests

### Phase 3: Interface & Integration
- [ ] **API**: Create endpoint `GET /path/to/resource`
- [ ] **Docs**: Update OpenAPI/Swagger spec
- [ ] **E2E**: Verify flow with integration test
```

## Rules
-   **No Assumptions**: Verify file paths and existing class names before creating tasks.
-   **Atomic**: One task = one logical commit (approx).
-   **Test-Driven**: Always include testing tasks *before* or *with* implementation tasks.
