---
description: Master orchestrator for multi-phase development with parallel execution support.
---

# Development Pipeline Orchestrator

User Request: $1

> **Core Principle**: Delegate, don't do. Load agent files and follow their instructions.

## Pipeline Flow

```
USER REQUEST
     │
     ▼
┌─────────────────────────────────────────────────────────────────────┐
│  PHASE 1: GATEKEEPER                                                │
│  Load: .agent/agents/gatekeeper/AGENT.md                            │
│  Output: Refined Spec                                               │
│  Gate: User must approve spec                                       │
└─────────────────────────────────────────────────────────────────────┘
                                            │
                                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│  PHASE 2: ARCHITECT                                                 │
│  Load: .agent/agents/architect/AGENT.md                             │
│  Output: Schema + API Contract                                      │
│  Gate: User must approve design                                     │
└─────────────────────────────────────────────────────────────────────┘
                                            │
                                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│  PHASE 3: PLANNER                                                   │
│  Load: .agent/agents/planner/AGENT.md                               │
│  Output: Todo List with Waves + Builder count per wave              │
│  Gate: User must approve plan                                       │
└─────────────────────────────────────────────────────────────────────┘
                                            │
                                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│  PHASE 4 & 5: BUILD + REVIEW (Parallel per Wave)                    │
│                                                                     │
│  FOR EACH WAVE:                                                     │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  Spawn N Builders in PARALLEL (N = min(tasks, 5))             │  │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐                        │  │
│  │  │Builder 1│  │Builder 2│  │Builder 3│  ...up to 5            │  │
│  │  │ Task A  │  │ Task B  │  │ Task C  │                        │  │
│  │  └────┬────┘  └────┬────┘  └────┬────┘                        │  │
│  │       │            │            │                             │  │
│  │       ▼            ▼            ▼                             │  │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐                        │  │
│  │  │Reviewer │  │Reviewer │  │Reviewer │  ← Each task reviewed  │  │
│  │  └────┬────┘  └────┬────┘  └────┬────┘                        │  │
│  │       │            │            │                             │  │
│  │       └────────────┴────────────┘                             │  │
│  │                    │                                          │  │
│  │        ALL tasks in wave APPROVED?                            │  │
│  │                    │                                          │  │
│  │       ┌────────────┴────────────┐                             │  │
│  │       ▼                         ▼                             │  │
│  │   [NO: Loop]            [YES: Next Wave]                      │  │
│  │   Fix failed tasks              │                             │  │
│  └─────────────────────────────────┼─────────────────────────────┘  │
│                                    │                                │
│                                    ▼                                │
│              [More waves?] ──Yes──► Process next wave               │
│                          │                                          │
│                          No                                         │
│                          ▼                                          │
│                     [COMPLETE]                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Execution

### Phase 1-3: Sequential (Gatekeeper → Architect → Planner)

Same as before - load agent, follow instructions, wait for approval.

### Phase 4 & 5: Parallel Build + Review

```
FOR EACH WAVE in plan:
  1. Get tasks for this wave
  2. Calculate builder_count = min(tasks_in_wave, 5)
  3. Spawn builder_count Builder agents IN PARALLEL
     - Each Builder handles one task
     - Use Task tool with parallel invocations
  4. WAIT for ALL Builders to complete
  5. FOR EACH completed task:
     - Load Reviewer agent
     - Review the task
     - If NEEDS_CHANGES → Builder fixes → Re-review (max 3 rounds)
  6. ALL tasks in wave APPROVED?
     - YES → Proceed to next wave
     - NO → Fix remaining tasks → Loop
  7. After all waves complete → Return to user
```

### Parallel Execution Example

```
Wave 1 has 3 tasks:
  → Spawn 3 Builders in parallel (single message, 3 Task tool calls)
  → Wait for all 3 to complete
  → Review each task
  → All approved → Move to Wave 2

Wave 2 has 2 tasks:
  → Spawn 2 Builders in parallel
  → Wait for both to complete
  → Review each task
  → All approved → Move to Wave 3
```

---

## Context Passing

| From | To | Context |
|------|----|---------|
| Gatekeeper | Architect | Refined Spec, Tech Stack |
| Architect | Planner | Schema, API Contract |
| Planner | Orchestrator | Todo List with Waves, Builder count |
| Orchestrator | Builders | Task assignment, Contract reference |
| Builder | Reviewer | Implementation, Test results |
| Reviewer | Builder | Feedback (if NEEDS_CHANGES) |

---

## Rules

1. **Load agent file** before each phase
2. **Never skip approval gates**
3. **Never write code** in Gatekeeper/Architect/Planner
4. **Parallel within waves** - Tasks in same wave run in parallel
5. **Sequential between waves** - Wait for wave N before starting wave N+1
6. **Max 5 Builders per wave** - Prevent resource exhaustion
7. **MANDATORY Review** - Every task must pass Reviewer
8. **Max 3 review rounds** - Escalate if issues persist

---

## Error Recovery

| Situation | Action |
|-----------|--------|
| Request unclear | Gatekeeper asks questions → Loop |
| User rejects design | Architect revises → Loop |
| User rejects plan | Planner revises → Loop |
| One task fails in wave | Other tasks continue, fix failed task |
| Review fails | Builder fixes → Re-submit (max 3) |
| 3 rounds exceeded | Escalate to user |
