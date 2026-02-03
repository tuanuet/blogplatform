# Agent & Skill Registry

## Orchestration

**Location**: `.agent/workflows/pipeline.md`

**Trigger**: `/pipeline` command or auto-detected for multi-phase development

**Responsibility**: Coordinates the 5-Phase Pipeline by loading and following agent instructions in sequence.

**Flow**: Gatekeeper → Architect → Planner → Builder ⇄ Reviewer

---

## Agents

### Gatekeeper Agent

**Role**: Technical Product Manager

**Location**: `.agent/agents/gatekeeper/AGENT.md`

**Skills**:

- `requirement-analysis`
- `tech-stack-detect`
- `ckb-code-scan`
- `documentation`

**Input**: Raw user request

**Output**:

- Refined Spec (User Stories + Edge Cases)
- OR Clarifying Questions (if ambiguous)

**Stop Condition**: DO NOT proceed if request is vague

---

### Architect Agent

**Role**: System Architect

**Location**: `.agent/agents/architect/AGENT.md`

**Skills**:

- `schema-design`
- `api-contract`
- `design-patterns`
- `ckb-code-scan`
- `documentation`

**Input**: Refined Spec from Gatekeeper

**Output**:

- Database Schema (auto-detect format from codebase)
- API Contract (OpenAPI/Interface)

**Constraint**: DO NOT write function bodies

---

### Planner Agent

**Role**: Technical Lead

**Location**: `.agent/agents/planner/AGENT.md`

**Skills**:

- `task-breakdown`
- `ckb-code-scan`
- `requirement-analysis`

**Input**: Architect's Design OR Bug Report

**Output**:

- Atomic, sequential Todo List

**Constraint**: Tasks must be implementable in one TDD cycle

---

### Builder Agent

**Role**: Senior Developer

**Location**: `.agent/agents/builder/AGENT.md`

**Skills**:

- `tdd-workflow`
- `clean-code`
- `testing`
- `mock-testing`
- `refactoring`
- `code-review`
- `ckb-code-scan`
- `documentation`

**Input**: API Contract from Architect + Todo List from Planner

**Output**:

- Failing tests (RED)
- Implementation (GREEN)
- Refactored code (REFACTOR)

**Workflow**: TDD cycle - RED → GREEN → REFACTOR

**Handoff**: Pass to Reviewer for verification

---

### Reviewer Agent

**Role**: Code Reviewer

**Location**: `.agent/agents/reviewer/AGENT.md`

**Skills**:

- `code-review`
- `testing`
- `clean-code`
- `design-patterns`
- `ckb-code-scan`

**Input**: Builder's implementation + API Contract + Refined Spec

**Output**:

- APPROVED (task complete)
- NEEDS_CHANGES (feedback to Builder)

**Workflow**: Review → Feedback Loop → Until Approved

**Constraint**: Max 3 review rounds, then escalate to user

---

## Pipeline Flow

```
┌─────────────────────────────────────────────────────┐
│  User Request                                       │
│       ↓                                             │
│  [/pipeline workflow] ──→ Load appropriate agent    │
│       │                                             │
│       ↓                                             │
│  [Gatekeeper] ──→ Refined Spec or Questions         │
│       │                                             │
│       ↓                                             │
│  [Architect] ──→ Schema + API Contract              │
│       │         (STOP: wait for approval)           │
│       ↓                                             │
│  [Planner] ──→ Implementation Plan (Todo List)      │
│       │         (STOP: wait for approval)           │
│       ↓                                             │
│  ┌─────────────────────────────────────────────┐    │
│  │ FOR EACH TASK:                              │    │
│  │                                             │    │
│  │  [Builder] ──→ TDD Implementation           │    │
│  │       │                                     │    │
│  │       ↓                                     │    │
│  │  [Reviewer] ──→ APPROVED or NEEDS_CHANGES   │    │
│  │       │              │                      │    │
│  │       │         NEEDS_CHANGES               │    │
│  │       │              ↓                      │    │
│  │       │         Back to Builder (loop)      │    │
│  │       │                                     │    │
│  │       ↓ APPROVED                            │    │
│  │  Mark task complete                         │    │
│  │       ↓                                     │    │
│  │  Next task...                               │    │
│  └─────────────────────────────────────────────┘    │
│       │                                             │
│       ↓                                             │
│  Return to User                                     │
└─────────────────────────────────────────────────────┘
```
