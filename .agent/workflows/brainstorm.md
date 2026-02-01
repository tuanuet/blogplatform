---
description: Brainstorming and ideation before development
---

# Brainstorming Workflow

Structured idea exploration and solution design **BEFORE** entering the development pipeline.

## When to Use

- When exploring a new product idea or feature concept
- When solving a complex problem with multiple possible solutions
- When needing to evaluate trade-offs between approaches
- When the requirement is vague and needs creative exploration

## Phases

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   IDEATE    │ ──▶ │   ANALYZE   │ ──▶ │   REFINE    │
│  (Explore)  │     │  (Evaluate) │     │  (Document) │
└─────────────┘     └─────────────┘     └─────────────┘
```

### Phase 1: IDEATE

Generate and explore ideas without judgment.

```
1. Define the problem statement clearly
2. Use brainstorming techniques (SCAMPER, Six Hats, etc.)
3. Generate at least 3-5 alternative approaches
4. No filtering at this stage - quantity over quality
```

- Skills: `brainstorming`, `ideation`
- Output: List of raw ideas and approaches

### Phase 2: ANALYZE

Evaluate and compare the generated ideas.

```
1. Create decision matrix with weighted criteria
2. Assess pros/cons for top candidates
3. Identify risks and dependencies
4. Check technical feasibility
```

- Skills: `solution-evaluation`
- Output: Ranked options with analysis

### Phase 3: REFINE

Select and document the chosen solution.

```
1. Choose the best option based on analysis
2. Document the design decision (ADR format)
3. Create high-level design diagram
4. Identify open questions for next phase
```

- Output: Design document ready for `/pipeline`

## Output Checklist

- [ ] Problem statement defined
- [ ] Multiple alternatives explored
- [ ] Decision matrix completed
- [ ] Chosen solution documented
- [ ] Design diagram created
- [ ] Open questions listed

## Integration with Other Workflows

```
/brainstorm  ──▶  /pipeline  ──▶  /new-feature
                            ──▶  /bug-fix
                            ──▶  /refactor
```

Use `/brainstorm` when you need to **explore** before you **build**.

## Example

```
Problem: "How should we implement user notifications?"

IDEATE output:
1. Push notifications (Firebase)
2. Email notifications (SendGrid)
3. In-app notifications (WebSocket)
4. SMS notifications (Twilio)

ANALYZE output:
| Option    | Cost | Speed | Reliability | Score |
|-----------|------|-------|-------------|-------|
| In-app    | 5    | 5     | 4           | 14    |
| Push      | 4    | 4     | 4           | 12    |
| Email     | 3    | 3     | 5           | 11    |

REFINE output:
→ Design doc: "ADR-001: In-app notifications with WebSocket"
→ Diagram: Client ↔ WebSocket Server ↔ Redis Pub/Sub
→ Ready for /pipeline
```
