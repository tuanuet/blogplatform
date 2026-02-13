---
description: Implementation workflow for pre-defined specs (Token-Optimized).
---

# Implementation workflow

Input: $1 (feature spec)

Use this workflow when the spec is already approved. Skip requirement discovery
and go straight to design and implementation.

## Flow

1. Design (architect)
   - Contracts and phased plan
   - Test plan and risks

2. Build core (builder)
   - Unit tests and minimal implementation

3. Implementation review (reviewer)
   - Diff-only review
   - Reviewer loads only `skill(name="consolidated-review")`

4. Integrate (builder)
   - Wire components (routes, DB, cache, config)
   - Add integration or e2e tests when needed

5. Final review (reviewer)
   - Verify acceptance criteria and test suite

## Token rules

- Reviewer reads only changed symbols or lines
- Do not paste full files into reviews
- Reuse `task_id` for loops (max 3)
