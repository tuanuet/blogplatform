# Agent & Skill Registry

## Orchestration

### Pipeline Workflow

**Location**: `.agent/workflows/pipeline.md`

**Trigger**: `/pipeline` command or auto-detected for multi-phase development

**Responsibility**: Coordinates the 5-Phase Pipeline by loading and following agent instructions in sequence.

**Flow**: Gatekeeper → Architect → Planner → Builder ⇄ Reviewer

### Brainstorm Workflow

**Location**: `.agent/workflows/brainstorm.md`

**Trigger**: `/brainstorm` command

**Responsibility**: Collaborative feature discussion and requirement definition. Pre-cursor to `/pipeline`.

**Flow**: Discuss → Define → Prepare → Output Feature Spec for `/pipeline`

---

## Agents

### Brainstormer Agent

**Role**: Creative Facilitator

**Location**: `.agent/agents/brainstormer/AGENT.md`

**Skills**:

- `brainstorming`
- `ideation`
- `requirement-analysis`
- `ckb-code-scan`

**Input**: Raw feature idea (may be vague)

**Output**:

- Feature Specification (ready for `/pipeline`)

**Workflow**: DISCUSS (context) → DEFINE (specs) → PREPARE (output)

**Integration**: `/brainstorm` → Feature Spec → `/pipeline` (Gatekeeper)

---

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

**⚠️ MANDATORY**: Must ask user clarifying questions and loop until ALL requirements are clear. User must explicitly approve the Refined Spec before proceeding.

**Stop Condition**: DO NOT proceed if request is vague or user hasn't confirmed spec

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

**⚠️ MANDATORY**: Must ask user about design decisions and loop until ALL architectural choices are confirmed. User must explicitly approve the design before proceeding.

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
│       │         ⚠️ MANDATORY LOOP:                  │
│       │         - Ask clarifying questions          │
│       │         - Wait for user response            │
│       │         - Loop until ALL clear              │
│       │         - User MUST approve spec            │
│       ↓                                             │
│  [Architect] ──→ Schema + API Contract              │
│       │         ⚠️ MANDATORY LOOP:                  │
│       │         - Ask design questions              │
│       │         - Wait for user response            │
│       │         - Loop until ALL confirmed          │
│       │         - User MUST approve design          │
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
