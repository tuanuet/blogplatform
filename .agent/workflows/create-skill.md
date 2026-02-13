---
description: Create a new skill using the skill-creator
---

# Create-skill workflow

Input: $1

Use this workflow to add a reusable, loadable skill to the agent system.

## Execution note

You can run this workflow as a single agent. If the task is large, you can
optionally spawn specialized subagents via `task(...)` for specific phases.

## Workflow diagram

This diagram shows the end-to-end flow for adding a new skill.

```mermaid
flowchart TD
  A[Start] --> B[Define the skill]
  B --> C[Load guidelines: skill(name="skill-creator")]
  C --> D[Scaffold: create SKILL.md]
  D --> E[Draft SKILL.md]
  E --> F[Verify in fresh session]
  F --> G[Done: skill loads + examples work]
```

## Steps

1. Define the skill
   - Problem it solves
   - Trigger conditions (when to load it)
   - Inputs and outputs
   - Constraints (tools, safety rules, boundaries)

2. Load the guidelines
   - `skill(name="skill-creator")`

3. Scaffold
   - Create: `.opencode/skills/<skill-name>/SKILL.md`

4. Draft `SKILL.md`
   - Role and goal
   - Rules and constraints
   - Tool workflow
   - One or two realistic examples

5. Verify
   - Load the new skill in a fresh session
   - Run a tiny dry-run prompt to confirm behavior

## Exit criteria

- The skill loads successfully
- The examples produce the intended behavior
- The skill is narrow and reusable (not project-specific trivia)
