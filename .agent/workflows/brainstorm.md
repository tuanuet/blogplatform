---
description: Discuss and define a feature before development
---

# Brainstorming Workflow

Collaborative feature discussion and requirement definition. This workflow serves as the **pre-cursor** to the development pipeline.

## When to Use

- When you have a feature idea but need to flesh out the details
- When you need to discuss requirements, edge cases, and user flows
- When you want to prepare a clear feature specification for development
- **Goal**: Turn a rough idea into a structured input for `/pipeline`

## Phases

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   DISCUSS   │ ──▶ │   DEFINE    │ ──▶ │   PREPARE   │
│  (Context)  │     │  (Specs)    │     │  (Output)   │
└─────────────┘     └─────────────┘     └─────────────┘
```

### Phase 1: DISCUSS (Context & Scope)

Understand the "Why" and "What".

```
1. Discuss the feature's goal and user value
2. Identify target users and use cases
3. Explore potential technical approaches (high-level)
4. Discuss constraints and non-functional requirements
```

- Skills: `brainstorming`, `ideation`, `requirement-analysis`
- Output: Clear understanding of feature scope

### Phase 2: DEFINE (Specifics)

Flesh out the "How".

```
1. Define user stories and acceptance criteria
2. Outline key user flows (happy path & edge cases)
3. Identify necessary data models/fields
4. List required API endpoints or interface changes
```

- Output: Detailed feature points

### Phase 3: PREPARE (Pipeline Input)

Structure the output for the development pipeline.

```
1. Summarize the feature into a structured format
2. Create a "Feature Request" document
3. Ensure all inputs required by /pipeline (Gatekeeper) are present
```

- Output: A structured **Feature Specification** ready for `/pipeline`

## Output Format (for /pipeline)

The final output should be a clear prompt/document containing:

```markdown
# Feature: [Feature Name]

## Objective

[Brief description of what we are building and why]

## Requirements

- [ ] Requirement 1
- [ ] Requirement 2

## Technical Context

- Impacted Services: [List]
- New Data/Fields: [Brief list]

## Acceptance Criteria

1. User can...
2. System should...
```

## Integration

```
/brainstorm (This Workflow) ──▶ Output: Feature Spec ──▶ /pipeline (Gatekeeper Phase)
```

Use `/brainstorm` to **talk** about the feature.
Use `/pipeline` to **build** the feature.
