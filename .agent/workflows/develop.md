---
description: Master orchestrator for multi-phase development with parallel execution support.
---

# Development Workflow

User Request: $1

> **Core Principle**: Delegate, don't do. Load agent files and follow their instructions.

## Pipeline Flow

```mermaid
flowchart TB
    Start([User Request]) --> Gatekeeper[Phase 1: Gatekeeper]
    
    Gatekeeper -->|User approves spec| Architect[Phase 2: Architect]
    Gatekeeper -.->|Needs clarification| Gatekeeper
    
    Architect -->|User approves design| Planner[Phase 3: Planner]
    Architect -.->|Needs revision| Architect
    
    Planner -->|User approves plan| Wave{For Each Wave}
    Planner -.->|Needs revision| Planner
    
    Wave -->|Tasks in wave| ParallelBuild[Parallel Build Phase]
    
    subgraph ParallelBuild[Parallel Execution]
        direction TB
        Builder1[Builder: Task 1] --> Review1[Reviewer]
        Builder2[Builder: Task 2] --> Review2[Reviewer]
        Builder3[Builder: Task N] --> Review3[Reviewer]
    end
    
    Review1 & Review2 & Review3 --> AllApproved{All Tasks\nApproved?}
    
    AllApproved -->|NO| FixFailed[Fix Failed Tasks]
    FixFailed --> ParallelBuild
    
    AllApproved -->|YES| MoreWaves{More Waves?}
    
    MoreWaves -->|YES| Wave
    MoreWaves -->|NO| Complete([Complete])
    
    style Gatekeeper fill:#e1f5fe
    style Architect fill:#e8f5e9
    style Planner fill:#fff3e0
    style ParallelBuild fill:#fce4ec
    style Complete fill:#c8e6c9
```

---

## Execution

### Phase 1-3: Sequential (Gatekeeper → Architect → Planner)

Same as before - load agent, follow instructions, wait for approval.

### Phase 4 & 5: Parallel Build + Review

```
FOR EACH WAVE in plan:
  1. Get tasks for this wave
  2. Calculate builder_count = min(tasks_in_wave, 2)
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

| From         | To           | Context                             |
| ------------ | ------------ | ----------------------------------- |
| Gatekeeper   | Architect    | Refined Spec, Tech Stack            |
| Architect    | Planner      | Schema, API Contract                |
| Planner      | Orchestrator | Todo List with Waves, Builder count |
| Orchestrator | Builders     | Task assignment, Contract reference |
| Builder      | Reviewer     | Implementation, Test results        |
| Reviewer     | Builder      | Feedback (if NEEDS_CHANGES)         |

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

| Situation              | Action                                |
| ---------------------- | ------------------------------------- |
| Request unclear        | Gatekeeper asks questions → Loop      |
| User rejects design    | Architect revises → Loop              |
| User rejects plan      | Planner revises → Loop                |
| One task fails in wave | Other tasks continue, fix failed task |
| Review fails           | Builder fixes → Re-submit (max 3)     |
| 3 rounds exceeded      | Escalate to user                      |
