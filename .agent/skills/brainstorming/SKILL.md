---
name: brainstorming
description: Structured brainstorming with proven creativity frameworks
---

# Brainstorming Skill

## Purpose

Apply structured brainstorming techniques to generate creative solutions systematically.

## When to Use

- Starting ideation phase for new features
- Exploring alternative approaches to a problem
- Breaking through creative blocks
- Team brainstorming facilitation

## Techniques

### SCAMPER Framework

Transform existing ideas using 7 lenses:

| Letter | Action     | Question                                  |
| ------ | ---------- | ----------------------------------------- |
| S      | Substitute | What can be replaced?                     |
| C      | Combine    | What can be merged or bundled?            |
| A      | Adapt      | What ideas from other domains apply?      |
| M      | Modify     | What can be enlarged, minimized, changed? |
| P      | Put to use | What other uses could this have?          |
| E      | Eliminate  | What can be removed or simplified?        |
| R      | Reverse    | What if we did the opposite?              |

```markdown
## SCAMPER Analysis for: [Feature/Problem]

- **Substitute**: Replace X with Y because...
- **Combine**: Merge A and B to create...
- **Adapt**: Borrow from [industry] the concept of...
- **Modify**: Scale up/down the...
- **Put to use**: Repurpose for...
- **Eliminate**: Remove unnecessary...
- **Reverse**: Instead of X→Y, try Y→X...
```

### Six Thinking Hats

Explore problem from 6 perspectives:

| Hat       | Focus      | Questions                     |
| --------- | ---------- | ----------------------------- |
| ⚪ White  | Facts      | What data do we have?         |
| 🔴 Red    | Emotions   | How do users feel about this? |
| ⚫ Black  | Risks      | What could go wrong?          |
| 🟡 Yellow | Benefits   | What are the advantages?      |
| 🟢 Green  | Creativity | What new ideas emerge?        |
| 🔵 Blue   | Process    | What's our next step?         |

```markdown
## Six Hats Analysis for: [Topic]

### ⚪ White Hat (Facts)

- Current metrics: ...
- User feedback data: ...

### 🔴 Red Hat (Emotions)

- Users feel frustrated when...
- This creates excitement because...

### ⚫ Black Hat (Risks)

- Risk 1: ...
- Risk 2: ...

### 🟡 Yellow Hat (Benefits)

- Benefit 1: ...
- Benefit 2: ...

### 🟢 Green Hat (Creativity)

- Wild idea 1: ...
- Alternative approach: ...

### 🔵 Blue Hat (Process)

- Next step: ...
- Decision needed: ...
```

### 5 Whys

Dig deeper to find root cause:

```markdown
## 5 Whys Analysis

**Problem**: [Initial problem statement]

1. **Why?** → [First answer]
2. **Why?** → [Deeper answer]
3. **Why?** → [Even deeper]
4. **Why?** → [Getting closer]
5. **Why?** → [Root cause identified]

**Root Cause**: [Summary]
**Solution Direction**: [Based on root cause]
```

### Mind Mapping

Structure ideas visually:

```
                    ┌─── Sub-idea 1.1
          ┌─ Idea 1 ┤
          │         └─── Sub-idea 1.2
          │
CENTRAL ──┼─ Idea 2 ─── Sub-idea 2.1
TOPIC     │
          │         ┌─── Sub-idea 3.1
          └─ Idea 3 ┤
                    └─── Sub-idea 3.2
```

## Output Format

```json
{
  "technique_used": "SCAMPER",
  "problem_statement": "How to improve X",
  "ideas": [
    { "id": 1, "description": "...", "category": "Substitute" },
    { "id": 2, "description": "...", "category": "Combine" }
  ],
  "top_3": [1, 2, 5],
  "next_step": "Proceed to solution-evaluation skill"
}
```

## Best Practices

1. **No judgment during ideation** - Evaluate later
2. **Build on ideas** - "Yes, and..." not "No, but..."
3. **Quantity first** - Aim for 10+ ideas before filtering
4. **Time-box** - Set 15-30 min limits per technique
5. **Document everything** - Capture all ideas, even "crazy" ones
