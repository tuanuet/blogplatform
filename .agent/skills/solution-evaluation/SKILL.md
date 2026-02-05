---
name: solution-evaluation
description: Evaluate and compare solutions using decision frameworks
---

# Solution Evaluation Skill

## Purpose

Objectively evaluate and compare solution options to make informed decisions.

## When to Use

- After brainstorming generates multiple options
- When choosing between architectural approaches
- When making build vs. buy decisions
- When prioritizing features or solutions

## Techniques

### Decision Matrix

Score options against weighted criteria:

```markdown
## Decision Matrix

**Options under evaluation**:

1. Option A: [Description]
2. Option B: [Description]
3. Option C: [Description]

**Criteria** (weight 1-5):
| Criteria | Weight | Option A | Option B | Option C |
| -------------- | ------ | -------- | -------- | -------- |
| Cost | 4 | 3 | 5 | 2 |
| Time to build | 3 | 4 | 3 | 5 |
| Maintainability| 5 | 5 | 3 | 4 |
| Scalability | 4 | 4 | 4 | 3 |
| Team expertise | 3 | 5 | 2 | 4 |

**Weighted scores**:

- Option A: (4×3)+(3×4)+(5×5)+(4×4)+(3×5) = **76**
- Option B: (4×5)+(3×3)+(5×3)+(4×4)+(3×2) = **66**
- Option C: (4×2)+(3×5)+(5×4)+(4×3)+(3×4) = **67**

**Winner**: Option A (score: 76)
```

### Pros/Cons Analysis

Structured advantage/disadvantage listing:

```markdown
## Pros/Cons Analysis

### Option A: [Name]

**Pros** ✅

1. [Advantage 1]
2. [Advantage 2]
3. [Advantage 3]

**Cons** ❌

1. [Disadvantage 1]
2. [Disadvantage 2]

**Mitigations for cons**:

- Con 1 → Mitigation: ...
- Con 2 → Mitigation: ...

---

### Option B: [Name]

**Pros** ✅

1. ...

**Cons** ❌

1. ...

---

**Summary**:
| Aspect | Option A | Option B |
| ---------- | -------- | -------- |
| Total Pros | 3 | 4 |
| Total Cons | 2 | 3 |
| Mitigable | 2/2 | 1/3 |
```

### Risk Assessment

Evaluate potential risks:

```markdown
## Risk Assessment Matrix

| Risk     | Probability | Impact     | Risk Score | Mitigation |
| -------- | ----------- | ---------- | ---------- | ---------- |
| [Risk 1] | High (3)    | High (3)   | 9 🔴       | [Strategy] |
| [Risk 2] | Medium (2)  | High (3)   | 6 🟡       | [Strategy] |
| [Risk 3] | Low (1)     | Medium (2) | 2 🟢       | [Strategy] |

**Risk Score**: Probability × Impact (1-3 scale)

- 🔴 High (7-9): Needs active mitigation
- 🟡 Medium (4-6): Monitor and prepare
- 🟢 Low (1-3): Accept or monitor

**Top risks requiring attention**:

1. Risk 1 - Mitigation plan: ...
2. Risk 2 - Mitigation plan: ...
```

### Feasibility Check

Assess implementation viability:

```markdown
## Feasibility Assessment

### Technical Feasibility

- [ ] Team has required skills
- [ ] Technology is proven/mature
- [ ] Integration with existing systems possible
- [ ] Performance requirements achievable

**Score**: X/4 ✅

### Time Feasibility

- Estimated effort: [X person-weeks]
- Available time: [Y weeks]
- Buffer needed: [Z weeks]
- **Feasible?**: Yes/No

### Resource Feasibility

- Required: [N developers, $X budget]
- Available: [M developers, $Y budget]
- Gap: [describe if any]
- **Feasible?**: Yes/No

### Operational Feasibility

- [ ] Users will accept the change
- [ ] Training requirements manageable
- [ ] Support team can handle
- [ ] Rollback plan exists

**Overall Feasibility**: ✅ Go / ⚠️ Conditional / ❌ No-Go
```

### Trade-off Analysis

Compare competing priorities:

```markdown
## Trade-off Analysis

**Classic trade-offs**:

| Choose     | Over        | Reason                    |
| ---------- | ----------- | ------------------------- |
| Speed      | Features    | MVP first, iterate        |
| Quality    | Speed       | Long-term maintainability |
| Simplicity | Flexibility | YAGNI principle           |

**For this decision**:

We choose **[Option A]** because:

1. Trade-off 1 aligns with business priority
2. Trade-off 2 is acceptable given timeline
3. Trade-off 3 can be revisited later

**What we're giving up**:

- [Feature/Quality X] - acceptable because...
- [Capability Y] - can add in v2
```

## Output Format

```json
{
  "evaluation_method": "Decision Matrix",
  "options_evaluated": 3,
  "winner": {
    "name": "Option A",
    "score": 76,
    "confidence": "High"
  },
  "key_trade_offs": ["Speed over features", "Simplicity over flexibility"],
  "risks_identified": 3,
  "high_risks": 1,
  "feasibility": "Go",
  "recommendation": "Proceed with Option A, with mitigation plan for Risk 1",
  "next_step": "Create design document and proceed to /develop"
}
```

## Decision Documentation (ADR Format)

```markdown
# ADR-XXX: [Decision Title]

## Status

Proposed | Accepted | Deprecated | Superseded

## Context

[Why this decision was needed]

## Decision

[What was decided]

## Options Considered

1. Option A - [summary]
2. Option B - [summary]
3. Option C - [summary]

## Consequences

### Positive

- ...

### Negative

- ...

### Risks

- ...

## References

- [Link to analysis]
- [Link to discussion]
```

## Best Practices

1. **Use multiple techniques** - No single method captures everything
2. **Involve stakeholders** - Different perspectives improve evaluation
3. **Document reasoning** - Future you will thank present you
4. **Revisit periodically** - Decisions may need updating
5. **Accept uncertainty** - Perfect information rarely exists
