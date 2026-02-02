---
description: Create a new skill using the skill-creator
---

# Create Skill Workflow

Uses the `skill-creator` skill to guide the creation of high-quality, effective new skills for the agent system.

## When to Use

- You want to teach the agent a new capability or domain knowledge.
- You need to standardize a complex workflow that isn't covered by existing skills.
- You want to package a set of instructions and best practices for reuse.

## Phases

### Phase 1: Preparation & Design

1. **Load Skill**: `skill-creator`
   - Use the command: `Skill(name="skill-creator")`
2. **Define Context**:
   - What problem does this skill solve?
   - What are the inputs and outputs?
   - What tools will this skill use?

### Phase 2: Creation

1. **Scaffold Directory**:
   - Create a new directory: `.agent/skills/<skill-name>/`
2. **Draft Content**:
   - Use the `skill-creator` guidelines to write the `SKILL.md` file.
   - Include sections for: Role, Description, Instructions, Examples, and Constraints.
3. **Review**:
   - Ensure the skill follows the standard format and best practices.

### Phase 3: Verification

1. **Test the Skill**:
   - Try to load the new skill in a new session to ensure it works.
   - Verify that the agent understands and follows the new instructions.

## Output Checklist

- [ ] New directory created in `.agent/skills/`
- [ ] `SKILL.md` file created and populated
- [ ] Skill functionality verified
