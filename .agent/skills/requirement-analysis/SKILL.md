---
name: requirement-analysis
description: Analyze and validate requirements, detect ambiguity, generate clarifying questions
---

# Requirement Analysis Skill

## Purpose

Analyze input requirements to detect ambiguity and ensure sufficient information before development.

## When to Use

- When receiving a new user request
- When validating requirement completeness
- When generating clarifying questions

## Ambiguity Detection Matrix

| Dimension | Question                       | Red Flags                  |
| --------- | ------------------------------ | -------------------------- |
| **WHO**   | Who will use this?             | "users", "everyone"        |
| **WHAT**  | What exactly needs to be done? | "fix", "improve", "update" |
| **WHY**   | Business value?                | No context provided        |
| **WHERE** | Which component/module?        | "somewhere", "in the app"  |
| **WHEN**  | Trigger conditions?            | "sometimes", "when needed" |
| **HOW**   | Technical constraints?         | No requirements specified  |

## Process

```
1. Parse user request
2. Run through 6W checklist
3. If missing info → Generate targeted questions
4. If complete → Mark as READY for next phase
```

## Question Templates

### Feature Request

```markdown
To implement this feature, I need clarification:

1. **User**: Who will use this feature? (guest/logged-in/admin)
2. **Scope**: Which screens/modules are affected?
3. **Edge cases**: What happens if [X fails/is empty/exceeds limit]?
4. **Priority**: Any performance/security requirements?
```

### Bug Report

```markdown
To fix this bug, I need to know:

1. **Reproduce**: Steps to reproduce the bug?
2. **Expected**: What should the correct behavior be?
3. **Actual**: What is currently happening?
4. **Environment**: Browser/device/version?
```

## Output Format

### When Incomplete

```json
{
  "status": "NEEDS_CLARIFICATION",
  "missing": ["WHO", "WHAT"],
  "questions": ["Question 1?", "Question 2?"]
}
```

### When Complete

```json
{
  "status": "READY",
  "summary": "Brief summary of requirement",
  "next": "GENERATE_REFINED_SPEC"
}
```
