---
description: Full 3-phase development pipeline (generic)
---

# Development Pipeline Workflow

The **master workflow** for any development task. Other workflows extend this.

## When to Use

- When you need the full process from requirements to implementation
- When you want explicit human approval at key design and planning stages

## The 4 Phases

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ GATEKEEPER  │ ──▶ │  ARCHITECT  │ ──▶ │   PLANNER   │ ──▶ │   BUILDER   │
│ (Refine)    │     │  (Design)   │     │ (Breakdown) │     │   (TDD)     │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
                               │                   │
                       [Wait for Confirm]  [Wait for Confirm]
                               │                   │
                               ▼                   ▼
```

### Phase 1: GATEKEEPER

- Load: `.agent/agents/gatekeeper/AGENT.md`
- Skills: `requirement-analysis`, `tech-stack-detect`
- Output: Refined Spec OR clarifying questions

### Phase 2: ARCHITECT

- Load: `.agent/agents/architect/AGENT.md`
- Skills: `schema-design`, `api-contract`, `design-patterns`
- Output: Schema + API Contract (NO code)
- **STOP CONDITION**: Present Design to User. **WAIT** for approval before proceeding to Planning.

### Phase 2.5: PLANNER

- Skill: `task-breakdown`
- Action: Break down the approved Architect design into atomic tasks using `todowrite`.
- **STOP CONDITION**: Present Task List to User. **WAIT** for approval before proceeding to Building.

### Phase 3: BUILDER

- Load: `.agent/agents/builder/AGENT.md`
- Skills: `tdd-workflow`, `clean-code`, `testing`
- Output: Tests + Implementation
- Input: Takes output of all previous phase

## Workflow Variants

| Workflow       | Phases         | Use Case                             |
| -------------- | -------------- | ------------------------------------ |
| `/pipeline`    | All 4          | Generic, full process                |
| `/new-feature` | All 4 + extras | New features with migration/rollback |
| `/bug-fix`     | 1, 2.5, 3      | Skip Architect if no schema changes  |
| `/refactor`    | 2.5, 3         | Skip Gatekeeper & Architect          |
