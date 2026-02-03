---
description: Master orchestrator for multi-phase development. Coordinates Gatekeeper → Architect → Planner → Builder → Reviewer with approval gates and feedback loops. Use for features, complex changes, or any task requiring structured development flow.
---

# Development Pipeline Orchestrator

> **Core Principle**: Delegate, don't do. Load agent files and follow their instructions.
> Never write code directly in orchestration phases.

## Quick Reference

```
/pipeline  → Full 5-phase flow (Gatekeeper → Architect → Planner → Builder → Reviewer)
/bug-fix   → Skip Architect (Gatekeeper → Planner → Builder → Reviewer)
/refactor  → Skip requirements (Planner → Builder → Reviewer)
```

---

## Pipeline Flow

```
USER REQUEST
     │
     ▼
┌─────────────────────────────────────────────────────────────────────┐
│  PHASE 1: GATEKEEPER                                                │
│  ├─ Load: .agent/agents/gatekeeper/AGENT.md                         │
│  ├─ Purpose: Validate requirements, detect ambiguity                │
│  └─ Output: Refined Spec OR Clarifying Questions                    │
│                          │                                          │
│         ┌────────────────┴────────────────┐                         │
│         ▼                                 ▼                         │
│    [Questions?]                    [Refined Spec]                   │
│         │                                 │                         │
│    Return to User ◀───────────────────────┘                         │
│    Wait for answers                       │                         │
│    Re-run Phase 1                         ▼                         │
└─────────────────────────────────────────────────────────────────────┘
                                            │
                                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│  PHASE 2: ARCHITECT                                                 │
│  ├─ Load: .agent/agents/architect/AGENT.md                          │
│  ├─ Input: Refined Spec                                             │
│  ├─ Purpose: Design schema and API contracts                        │
│  └─ Output: Schema + API Contract (NO code)                         │
│                          │                                          │
│                          ▼                                          │
│             ┌────────────────────────┐                              │
│             │  APPROVAL GATE         │                              │
│             │  Present design        │                              │
│             │  WAIT for user confirm │                              │
│             └────────────────────────┘                              │
│                          │                                          │
│         ┌────────────────┴────────────────┐                         │
│         ▼                                 ▼                         │
│    [Rejected]                       [Approved]                      │
│    Revise design                          │                         │
└─────────────────────────────────────────────────────────────────────┘
                                            │
                                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│  PHASE 2.5: PLANNER                                                 │
│  ├─ Load: .agent/agents/planner/AGENT.md                            │
│  ├─ Input: Schema + API Contract                                    │
│  ├─ Purpose: Break down into atomic tasks                           │
│  └─ Output: Todo List (via todowrite)                               │
│                          │                                          │
│                          ▼                                          │
│             ┌────────────────────────┐                              │
│             │  APPROVAL GATE         │                              │
│             │  Present task list     │                              │
│             │  WAIT for user confirm │                              │
│             └────────────────────────┘                              │
│                          │                                          │
│         ┌────────────────┴────────────────┐                         │
│         ▼                                 ▼                         │
│    [Rejected]                       [Approved]                      │
│    Revise plan                            │                         │
└─────────────────────────────────────────────────────────────────────┘
                                            │
                                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│  PHASE 3: BUILDER (per task)                                        │
│  ├─ Load: .agent/agents/builder/AGENT.md                            │
│  ├─ Input: Current Task + API Contract                              │
│  ├─ Purpose: TDD implementation (RED → GREEN → REFACTOR)            │
│  └─ Output: Tests + Implementation                                  │
│                          │                                          │
│                          ▼                                          │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  PHASE 4: REVIEWER (per task)                                 │  │
│  │  ├─ Load: .agent/agents/reviewer/AGENT.md                     │  │
│  │  ├─ Input: Builder's implementation                           │  │
│  │  ├─ Purpose: Verify quality, run tests, check criteria        │  │
│  │  └─ Output: APPROVED or NEEDS_CHANGES                         │  │
│  │                       │                                       │  │
│  │      ┌────────────────┴────────────────┐                      │  │
│  │      ▼                                 ▼                      │  │
│  │ [NEEDS_CHANGES]                   [APPROVED]                  │  │
│  │      │                                 │                      │  │
│  │      ▼                                 ▼                      │  │
│  │ Return to Builder              Mark task complete             │  │
│  │ with feedback                         │                       │  │
│  │      │                                │                       │  │
│  │      └──────► (loop max 3x) ◄─────────┘                       │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                          │                                          │
│                          ▼                                          │
│              [More tasks?] ──Yes──► Next task (Phase 3)             │
│                          │                                          │
│                          No                                         │
│                          ▼                                          │
│                     [COMPLETE]                                      │
│                  Return to User                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Execution Instructions

### Phase 1: GATEKEEPER

**Load**: `.agent/agents/gatekeeper/AGENT.md`

**Execute**:
1. Run ambiguity check using `requirement-analysis` skill
2. Scan codebase with `ckb-code-scan` skill:
   - `ckb_explore` for project overview
   - `ckb_getArchitecture` for module structure
3. Detect tech stack using `tech-stack-detect` skill

**Decision Point**:
- **Missing info?** → Generate questions → Return to user → Wait → Re-run
- **Complete?** → Output Refined Spec → Proceed to Phase 2

**Refined Spec must include**:
- User Stories with acceptance criteria
- Edge cases identified
- Tech stack context
- Affected modules/files

---

### Phase 2: ARCHITECT

**Load**: `.agent/agents/architect/AGENT.md`

**Execute**:
1. Analyze existing patterns with `ckb-code-scan`:
   - `ckb_searchSymbols` for similar entities
   - `ckb_understand` for existing patterns
2. Design schema using `schema-design` skill (auto-detect format)
3. Design API contracts using `api-contract` skill
4. Apply `design-patterns` skill (SOLID, Repository, Service)

**Output**:
- Database Schema (Prisma/Drizzle/SQL based on codebase)
- API Contract (TypeScript interfaces/OpenAPI)
- NO implementation code

**APPROVAL GATE**:
```
Present to user:
  - Schema design
  - API contracts
  - Design decisions rationale

Ask: "Does this design look correct? Approve to proceed or provide feedback."

WAIT for explicit approval before Phase 2.5
```

---

### Phase 2.5: PLANNER

**Load**: `.agent/agents/planner/AGENT.md`

**Execute**:
1. Analyze dependencies with `ckb-code-scan`
2. Decompose using `task-breakdown` skill:
   - 1 task = 1 TDD cycle (RED-GREEN-REFACTOR)
   - Order: Model → Repository → Service → Controller → UI
3. Write tasks via `todowrite` tool

**Task Format**:
```
[Priority] Task description
- High: Core functionality, blocking others
- Medium: Important but non-blocking
- Low: Nice-to-have, polish
```

**APPROVAL GATE**:
```
Present to user:
  - Numbered task list
  - Dependencies noted
  - Estimated scope

Ask: "Does this plan look correct? Approve to start building or adjust."

WAIT for explicit approval before Phase 3
```

---

### Phase 3: BUILDER

**Load**: `.agent/agents/builder/AGENT.md`

**For each task in Todo List**:

1. **PRE-ANALYSIS**:
   - `ckb_prepareChange` for impact analysis
   - `ckb_understand` for existing code patterns

2. **RED** (Write failing test):
   - Create test file based on API contract
   - Run test → Must fail
   - Commit: "test: add [feature] test"

3. **GREEN** (Make it pass):
   - Write minimal implementation
   - Run test → Must pass
   - Commit: "feat: implement [feature]"

4. **REFACTOR** (Clean up):
   - Apply `clean-code` skill
   - Run test → Must still pass
   - Commit: "refactor: clean [feature]"

5. **Handoff to Reviewer** → Proceed to Phase 4

---

### Phase 4: REVIEWER

**Load**: `.agent/agents/reviewer/AGENT.md`

**For each task after Builder completes**:

1. **RUN TESTS**:
   - Execute full test suite
   - Check: All tests pass?
   - Check: Coverage acceptable?

2. **CODE REVIEW**:
   - Security check (injection, XSS, auth)
   - Performance check (N+1, loops, indexing)
   - Clean code check (naming, functions, duplication)
   - Architecture check (patterns, dependencies)

3. **ACCEPTANCE CRITERIA**:
   - Verify against Refined Spec
   - Check all criteria implemented
   - Check edge cases handled

**Decision Point**:
```
┌─────────────────────────────────────────┐
│ All checks pass?                        │
│                                         │
│   YES → APPROVED                        │
│         • Mark task complete            │
│         • Proceed to next task          │
│                                         │
│   NO  → NEEDS_CHANGES                   │
│         • Generate structured feedback  │
│         • Return to Builder             │
│         • Builder fixes issues          │
│         • Re-submit for review          │
│         • (max 3 rounds)                │
└─────────────────────────────────────────┘
```

**Feedback Format** (when NEEDS_CHANGES):
```markdown
## Review Result: NEEDS_CHANGES

### Issues to Fix

#### CRITICAL (must fix)
1. [Security/Bug issue] at `file:line`
   - Problem: ...
   - Fix: ...

#### IMPORTANT (should fix)
2. [Performance/Code quality] at `file:line`
   - Problem: ...
   - Fix: ...

### Next Steps
1. Fix CRITICAL first
2. Address IMPORTANT
3. Re-submit for review
```

**Loop until APPROVED** (max 3 rounds, then escalate to user)

---

## Workflow Variants

| Command        | Phases                  | Skip                      | Use Case                   |
| -------------- | ----------------------- | ------------------------- | -------------------------- |
| `/pipeline`    | 1 → 2 → 2.5 → 3 → 4     | -                         | New features, full process |
| `/new-feature` | 1 → 2 → 2.5 → 3 → 4     | -                         | Explicit new feature       |
| `/bug-fix`     | 1 → 2.5 → 3 → 4         | Architect                 | Bug with no schema change  |
| `/refactor`    | 2.5 → 3 → 4             | Gatekeeper, Architect     | Code improvement only      |

---

## Error Handling

| Situation                       | Action                                       |
| ------------------------------- | -------------------------------------------- |
| Gatekeeper: Request unclear     | Generate questions → Return to user → Loop   |
| Architect: Missing requirements | Return to Gatekeeper for clarification       |
| Architect: User rejects design  | Revise design based on feedback → Re-present |
| Planner: User rejects plan      | Revise tasks based on feedback → Re-present  |
| Builder: Test fails             | Debug in Builder → Do NOT proceed until green|
| Builder: Blocked by dependency  | Check if earlier task missed → Add task      |
| Reviewer: Issues found          | Return to Builder with feedback → Loop       |
| Reviewer: 3 rounds exceeded     | Escalate to user for decision                |

---

## Context Passing

Each phase passes context to the next:

```
Gatekeeper → Architect:
  - Refined Spec (user stories, acceptance criteria, edge cases)
  - Tech Stack info
  - Codebase context

Architect → Planner:
  - Schema design
  - API contracts
  - Design patterns applied

Planner → Builder:
  - Todo list (ordered tasks)
  - All above context (Spec + Schema + Contract)

Builder → Reviewer:
  - Implementation (changed files)
  - Test results
  - API Contract reference
  - Refined Spec reference (for acceptance criteria)

Reviewer → Builder (if issues):
  - Structured feedback
  - Priority order
  - File:line references
  - Fix suggestions
```

---

## Rules

1. **Always start with Gatekeeper** (unless `/refactor`)
2. **Never skip approval gates** - User must confirm design and plan
3. **Never write code** in Gatekeeper/Architect/Planner phases
4. **Loop back if blocked** - Return to previous phase for clarification
5. **Complete all tasks** before returning final result
6. **One task at a time** - Mark in_progress, complete, then next
7. **Tests before code** - Builder must follow TDD strictly
8. **Review every task** - Reviewer validates before marking complete
9. **Max 3 review rounds** - Escalate if issues persist
