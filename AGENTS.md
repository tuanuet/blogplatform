# Agent Registry

## Orchestrator (Lead Agent)

**Role**: Pipeline Controller

**Location**: `.agent/agents/orchestrator/AGENT.md`

**Responsibility**: Orchestrates the entire 3-Phase Pipeline, delegating tasks to sub-agents in proper sequence.

**Workflow**:

1. Receive user request
2. Delegate to **Gatekeeper** → Get Refined Spec
   - If ambiguous → Return questions → Loop
3. Delegate to **Architect** → Get Schema + API Contract
4. Delegate to **Planner** → Get Todo List
5. Delegate to **Builder** → Get Tests + Implementation
6. Return final result

---

## Sub-Agents

### 🚪 Gatekeeper Agent

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

### 📐 Architect Agent

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

### 📋 Planner Agent

**Role**: Technical Lead

**Location**: `.agent/agents/planner/AGENT.md`

**Skills**:

- `todowrite`
- `ckb-code-scan`
- `requirement-analysis`

**Input**: Architect's Design OR Bug Report

**Output**:

- Atomic, sequential Todo List

**Constraint**: Tasks must be implementable in one TDD cycle

---

### 🔨 Builder Agent

**Role**: Senior Developer

**Location**: `.agent/agents/builder/AGENT.md`

**Skills**:

- `tdd-workflow`
- `clean-code`
- `testing`
- `refactoring`
- `code-review`
- `ckb-code-scan`
- `documentation`

**Input**: API Contract from Architect

**Output**:

- Failing tests (RED)
- Implementation (GREEN)
- Refactored code (REFACTOR)

**Workflow**: TDD cycle - RED → GREEN → REFACTOR

---

## Delegation Rules

```
┌─────────────────────────────────────────────────────┐
│  User Request                                       │
│       ↓                                             │
│  [Orchestrator] ──→ Is request clear?               │
│       │                   │                         │
│       │ No                │ Yes                     │
│       ↓                   ↓                         │
│  [Gatekeeper] ←──── Ask questions                   │
│       │                                             │
│       │ Refined Spec ready                          │
│       ↓                                             │
│  [Architect] ──→ Schema + API Contract              │
│       │                                             │
│       ↓                                             │
│  [Planner] ──→ Implementation Plan (Todo List)      │
│       │                                             │
│       ↓                                             │
│  [Builder] ──→ TDD Implementation                   │
│       │                                             │
│       ↓                                             │
│  Return to User                                     │
└─────────────────────────────────────────────────────┘
```
