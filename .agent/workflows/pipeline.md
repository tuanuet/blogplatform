---
description: Full 3-phase development pipeline (generic)
---

# Development Pipeline Workflow

The **master workflow** for any development task. Other workflows extend this.

## When to Use

- When you need the full 3-phase process
- As a reference for other workflows

## The 3 Phases

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ GATEKEEPER  │ ──▶ │  ARCHITECT  │ ──▶ │   BUILDER   │
│ (Refine)    │     │  (Design)   │     │   (TDD)     │
└─────────────┘     └─────────────┘     └─────────────┘
```

### Phase 1: GATEKEEPER

- Load: `.agent/agents/gatekeeper/AGENT.md`
- Skills: `requirement-analysis`, `tech-stack-detect`
- Output: Refined Spec OR clarifying questions

### Phase 2: ARCHITECT

- Load: `.agent/agents/architect/AGENT.md`
- Skills: `schema-design`, `api-contract`, `design-patterns`
- Output: Schema + API Contract (NO code)

### Phase 3: BUILDER

- Load: `.agent/agents/builder/AGENT.md`
- Skills: `tdd-workflow`, `clean-code`, `testing`
- Output: Tests + Implementation

## Workflow Variants

| Workflow       | Phases         | Use Case                             |
| -------------- | -------------- | ------------------------------------ |
| `/pipeline`    | All 3          | Generic, full process                |
| `/new-feature` | All 3 + extras | New features with migration/rollback |
| `/bug-fix`     | 1 + 3 only     | Skip Architect, focus on regression  |
| `/refactor`    | 3 only         | Skip Gatekeeper & Architect          |
