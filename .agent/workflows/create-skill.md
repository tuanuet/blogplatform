---
description: Create a new skill using the skill-creator
---

# Create-skill workflow

Input: $1

Use this workflow to add a reusable, loadable skill to the agent system.

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
