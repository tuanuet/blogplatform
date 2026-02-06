---
description: Master orchestrator for multi-phase development.
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
    
    Planner -->|User approves plan| Builder[Phase 4: Builder]
    Planner -.->|Needs revision| Planner
    
    Builder -->|Implementation| Reviewer[Phase 5: Reviewer]
    
    Reviewer -->|APPROVED| NextTask{More Tasks?}
    Reviewer -->|NEEDS_CHANGES| Builder
    
    NextTask -->|YES| Builder
    NextTask -->|NO| Complete([Complete])
    
    style Gatekeeper fill:#e1f5fe
    style Architect fill:#e8f5e9
    style Planner fill:#fff3e0
    style Builder fill:#fce4ec
    style Reviewer fill:#f3e5f5
    style Complete fill:#c8e6c9
```

---

## Execution

### Phase 1-3: Sequential (Gatekeeper → Architect → Planner)

Same as before - load agent, follow instructions, wait for approval.

### Phase 4 & 5: Build + Review Loop

```
FOR EACH TASK in plan:
  1. Load Builder agent
     - Input: Task details + API Contract
     - Output: Implementation + Tests
  2. Load Reviewer agent
     - Input: Implementation + Specs
     - Output: APPROVED or NEEDS_CHANGES
  3. IF NEEDS_CHANGES:
     - Builder fixes issues
     - Re-submit to Reviewer (max 3 rounds)
  4. IF APPROVED:
     - Mark task complete
     - Proceed to next task
```

---

## Context Passing

| From         | To           | Context                             |
| ------------ | ------------ | ----------------------------------- |
| Gatekeeper   | Architect    | Refined Spec, Tech Stack            |
| Architect    | Planner      | Schema, API Contract                |
| Planner      | Builder      | Task assignment, Contract reference |
| Builder      | Reviewer     | Implementation, Test results        |
| Reviewer     | Builder      | Feedback (if NEEDS_CHANGES)         |

---

## Rules

1. **Load agent file** before each phase
2. **Never skip approval gates**
3. **Never write code** in Gatekeeper/Architect/Planner
4. **MANDATORY Review** - Every task must pass Reviewer
5. **Max 3 review rounds** - Escalate if issues persist

---

## Error Recovery

| Situation              | Action                                |
| ---------------------- | ------------------------------------- |
| Request unclear        | Gatekeeper asks questions → Loop      |
| User rejects design    | Architect revises → Loop              |
| User rejects plan      | Planner revises → Loop                |
| Review fails           | Builder fixes → Re-submit (max 3)     |
| 3 rounds exceeded      | Escalate to user                      |
