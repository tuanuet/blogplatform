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
│  ├─ Purpose: Break down into atomic tasks with dependencies         │
│  └─ Output: Todo List with Waves + Parallelization Plan             │
│                          │                                          │
│                          ▼                                          │
│             ┌────────────────────────┐                              │
│             │  APPROVAL GATE         │                              │
│             │  Present task list     │                              │
│             │  Show dependency graph │                              │
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
│  PHASE 3 & 4: BUILD + REVIEW (per wave)                             │
│                                                                     │
│  FOR EACH WAVE:                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  PARALLEL BUILDERS (max 3)                                    │  │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐                        │  │
│  │  │Builder 1│  │Builder 2│  │Builder 3│  ← Spawn in parallel   │  │
│  │  │ Task A  │  │ Task B  │  │ Task C  │                        │  │
│  │  └────┬────┘  └────┬────┘  └────┬────┘                        │  │
│  │       │            │            │                             │  │
│  │       ▼            ▼            ▼                             │  │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐                        │  │
│  │  │Reviewer │  │Reviewer │  │Reviewer │  ← MANDATORY per task  │  │
│  │  │ Task A  │  │ Task B  │  │ Task C  │                        │  │
│  │  └────┬────┘  └────┬────┘  └────┬────┘                        │  │
│  │       │            │            │                             │  │
│  │       ▼            ▼            ▼                             │  │
│  │  [APPROVED?]  [APPROVED?]  [APPROVED?]                        │  │
│  │       │            │            │                             │  │
│  │       └────────────┴────────────┘                             │  │
│  │                    │                                          │  │
│  │        ALL tasks in wave APPROVED?                            │  │
│  │                    │                                          │  │
│  │       ┌────────────┴────────────┐                             │  │
│  │       ▼                         ▼                             │  │
│  │   [NO: Loop]            [YES: Next Wave]                      │  │
│  │   Return to Builder            │                              │  │
│  │   with feedback                │                              │  │
│  └────────────────────────────────┼──────────────────────────────┘  │
│                                   │                                 │
│                                   ▼                                 │
│              [More waves?] ──Yes──► Process next wave               │
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
3. **Dependency Graph Analysis**:
   - Identify which tasks depend on others (e.g., Repository depends on Entity)
   - Group independent tasks that can run in parallel
   - Create execution waves (tasks in same wave can be parallelized)
4. **Determine Builder Concurrency**:
   - Count tasks per wave
   - Recommend number of parallel Builders (max 3 for safety)
   - Document: "Wave 1: Tasks A, B (parallel) → Wave 2: Task C (depends on A, B)"
5. Write tasks via `todowrite` tool with dependency metadata

**Task Format**:
```
[Priority] [Wave N] Task description
- Priority: High/Medium/Low
- Wave: Execution order (Wave 1 runs first, Wave 2 after Wave 1 completes)
- Depends On: [Task IDs] (empty if independent)
```

**Parallelization Strategy**:
```
┌─────────────────────────────────────────────────────────────┐
│  WAVE 1 (Independent - can parallelize)                    │
│  ├─ Task A: Entity definition                              │
│  ├─ Task B: DTO definition                                 │
│  └─ Task C: Migration script                               │
│           ↓ (all complete)                                 │
│  WAVE 2 (Depends on Wave 1)                                │
│  ├─ Task D: Repository (depends on A)                      │
│  └─ Task E: Service interface (depends on A, B)            │
│           ↓ (all complete)                                 │
│  WAVE 3 (Depends on Wave 2)                                │
│  └─ Task F: Handler (depends on D, E)                      │
└─────────────────────────────────────────────────────────────┘

Recommended Builders: min(tasks_in_wave, 3)
```

**APPROVAL GATE**:
```
Present to user:
  - Numbered task list with waves
  - Dependency graph (text or ASCII)
  - Recommended parallel builders per wave
  - Estimated scope

Ask: "Does this plan look correct? Approve to start building or adjust."

WAIT for explicit approval before Phase 3
```

---

### Phase 3: BUILDER

**Load**: `.agent/agents/builder/AGENT.md`

**Execution Strategy** (based on Planner output):
```
FOR each Wave in execution plan:
  1. Determine parallel_count = min(tasks_in_wave, 3)
  2. Spawn `parallel_count` Builder agents simultaneously
  3. Each Builder handles one task from the wave
  4. WAIT for ALL Builders in wave to complete
  5. Run Reviewer for EACH completed task (MANDATORY)
  6. Proceed to next wave only after ALL tasks in current wave are APPROVED
```

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

5. **MANDATORY Handoff to Reviewer** → Proceed to Phase 4
   - ⚠️ **NEVER skip this step**
   - ⚠️ **DO NOT mark task complete until Reviewer approves**

---

### Phase 4: REVIEWER (MANDATORY)

> ⚠️ **CRITICAL**: This phase is MANDATORY after EVERY Builder task.
> The Orchestrator MUST invoke the Reviewer agent after each Builder completes.
> Skipping this phase is a violation of the pipeline protocol.

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
  - Todo list with waves (ordered tasks with dependencies)
  - Parallelization plan (which tasks can run together)
  - All above context (Spec + Schema + Contract)

Builder → Reviewer (MANDATORY):
  - Implementation (changed files)
  - Test results
  - API Contract reference
  - Refined Spec reference (for acceptance criteria)
  - ⚠️ Orchestrator MUST invoke Reviewer after EVERY Builder task

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
6. **One task at a time per Builder** - Mark in_progress, complete, then next
7. **Tests before code** - Builder must follow TDD strictly
8. **MANDATORY Review for every task** - Reviewer MUST validate before marking complete. NO EXCEPTIONS.
9. **Max 3 review rounds** - Escalate if issues persist
10. **Respect task dependencies** - Never start a task before its dependencies are APPROVED
11. **Parallel execution by wave** - Only parallelize tasks within the same wave
12. **Max 3 concurrent Builders** - To prevent resource exhaustion and context confusion
