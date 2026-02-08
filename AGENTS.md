# Agent & Skill Registry

## Orchestration

### Development Workflow

**Location**: `.agent/workflows/develop.md`

**Trigger**: `/develop` command or auto-detected for multi-phase development

**Responsibility**: Coordinates the development pipeline by loading and following agent instructions in sequence.

**Flow**: Gatekeeper → Architect → Reviewer (Architecture) → Builder Phase 2 → Reviewer (Implementation) → Builder Phase 3 → Reviewer (Integration)

### Brainstorm Workflow

**Location**: `.agent/workflows/brainstorm.md`

**Trigger**: `/brainstorm` command

**Responsibility**: Collaborative feature discussion and requirement definition. Pre-cursor to `/develop`.

**Flow**: Discuss → Define → Prepare → Output Feature Spec for `/develop`

### Implementation Workflow

**Location**: `.agent/workflows/implementation.md`

**Trigger**: `/implementation` command

**Responsibility**: Implementation workflow for pre-defined specs. Skips Gatekeeper phase, starts from Architect.

**Flow**: Architect → Reviewer (Architecture) → Builder Phase 2 → Reviewer (Implementation) → Builder Phase 3 → Reviewer (Integration)

### Document Workflow

**Location**: `.agent/workflows/document.md`

**Trigger**: `/document` command

**Responsibility**: Generate and verify documentation on demand.

**Flow**: Documenter → Document-Reviewer

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

- Feature Specification (ready for `/develop`)

**Workflow**: DISCUSS (context) → DEFINE (specs) → PREPARE (output)

**Integration**: `/brainstorm` → Feature Spec → `/develop` (Gatekeeper)

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
- Phase-Based Implementation Plan

**⚠️ MANDATORY**: 
- Must ask user about design decisions and loop until ALL architectural choices are confirmed
- User must explicitly approve the design and task plan before proceeding
- **After approval → Handoff to Reviewer for Architecture Review**

**Constraint**: DO NOT write function bodies

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

**Input**: API Contract from Architect + Todo List from Architect

**Output**:

- Failing tests (RED)
- Implementation (GREEN)
- Refactored code (REFACTOR)

**Workflow**: TDD cycle - RED → GREEN → REFACTOR

**Handoff**: Pass to Reviewer for verification

---

### Reviewer Agent

**Role**: Quality Gatekeeper - Reviews at 3 stages

**Location**: `.agent/agents/reviewer/AGENT.md`

**Skills**:

- `code-review`
- `testing`
- `clean-code`
- `design-patterns`
- `ckb-code-scan`

**Input**: 
- Phase 1: Contracts and Plan from Architect
- Phase 2: Component implementations from Builder
- Phase 3: Complete feature from Builder

**Output**:

- **Phase 1 (Architecture)**: APPROVED/NEEDS_CHANGES on contracts
- **Phase 2 (Implementation)**: APPROVED/NEEDS_CHANGES on components
- **Phase 3 (Integration)**: APPROVED/NEEDS_CHANGES on complete feature

**Workflow**: 3 Review Gates → Feedback Loop → Until Approved

**Constraint**: Max 3 review rounds per phase, then escalate to user

---

### Documenter Agent

**Role**: Documentation Specialist

**Location**: `.agent/agents/documenter/AGENT.md`

**Skills**:

- `ckb-code-scan`
- `mermaid-diagram-specialist`
- `api-contract`
- `documentation`

**Input**: User request (Scope/Type) + Codebase

**Output**:

- Architecture Diagrams (C4)
- User Flow Diagrams (Sequence)
- API Documentation

**Workflow**: Analyze → Scan → Generate → Save

**Trigger**: On-demand (e.g. "Document auth flow")

---

### Document-Reviewer Agent

**Role**: Documentation QA

**Location**: `.agent/agents/document-reviewer/AGENT.md`

**Skills**:

- `code-review`
- `ckb-code-scan`
- `api-contract`
- `documentation`

**Input**: Draft Documentation + Codebase

**Output**:

- APPROVED (ready to save)
- NEEDS_CHANGES (feedback to Documenter)

**Workflow**: Verify Accuracy → Verify Quality → Report

---

## Pipeline Flow

```
┌─────────────────────────────────────────────────────┐
│  User Request                                       │
│       ↓                                             │
│  [/develop workflow] ──→ Load appropriate agent    │
│       │                                             │
│       ↓                                             │
│  [Gatekeeper] ──→ Refined Spec or Questions         │
│       │         ⚠️ MANDATORY LOOP:                  │
│       │         - Ask clarifying questions          │
│       │         - Wait for user response            │
│       │         - Loop until ALL clear              │
│       │         - User MUST approve spec            │
│       ↓                                             │
│  [Architect] ──→ Contracts + Phase-Based Plan       │
│       │         ⚠️ MANDATORY LOOP:                  │
│       │         - Ask design questions              │
│       │         - Create implementation plan            │
│       │         - Wait for user response            │
│       │         - Loop until ALL confirmed          │
│       │         - User MUST approve design + plan       │
│       ↓                                             │
│  [Reviewer] ──→ Architecture Review                 │
│       │         "Contracts OK? Patterns OK?"        │
│       │         ├─ NEEDS_CHANGES → Back to Architect│
│       │         └─ APPROVED → Continue              │
│       ↓                                             │
│  ┌─────────────────────────────────────────────┐    │
│  │  PHASE 2: CORE IMPLEMENTATION               │    │
│  │                                             │    │
│  │  [Builder] ──→ Implement components         │    │
│  │       │                                     │    │
│  │       ↓                                     │    │
│  │  [Reviewer] ──→ Implementation Review       │    │
│  │       │         "Components work? Quality?"   │    │
│  │       │         ├─ NEEDS_CHANGES → Back       │    │
│  │       │         └─ APPROVED → Phase 3         │    │
│  └─────────────────────────────────────────────┘    │
│       ↓                                             │
│  ┌─────────────────────────────────────────────┐    │
│  │  PHASE 3: INTEGRATION                       │    │
│  │                                             │    │
│  │  [Builder] ──→ Wire up + Tests              │    │
│  │       │                                     │    │
│  │       ↓                                     │    │
│  │  [Reviewer] ──→ Integration Review          │    │
│  │       │         "Feature complete? E2E OK?" │    │
│  │       │         ├─ NEEDS_CHANGES → Back       │    │
│  │       │         └─ APPROVED → Complete!       │    │
│  └─────────────────────────────────────────────┘    │
│       ↓                                             │
│  Return to User                                     │
└─────────────────────────────────────────────────────┘

Total Reviews: 3 (Architecture, Implementation, Integration)
```
