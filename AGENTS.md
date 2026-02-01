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
4. Delegate to **Builder** → Get Tests + Implementation
5. Return final result

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
│  [Builder] ──→ TDD Implementation                   │
│       │                                             │
│       ↓                                             │
│  Return to User                                     │
└─────────────────────────────────────────────────────┘
```
