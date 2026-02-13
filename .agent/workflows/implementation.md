---
description: Implementation workflow for pre-defined specs (Token-Optimized).
---

# Implementation workflow

Input: $1 (feature spec)

Use this workflow when the spec is already approved. Skip requirement discovery
and go straight to design and implementation.

## Flow

Follow these steps to go from an approved spec to production-ready code with a
single review pass.

1. Design (architect)
   - Define contracts and a phased plan.
   - Define the test plan and risks.

2. Build end-to-end (builder)
   - Implement the core logic with unit tests.
   - Wire integration points (routes, DB, cache, config).
   - Add integration or end-to-end tests when needed.
   - Run the full test suite and any local checks used in this repo.

3. Single review (reviewer)
   - Review the implementation code with enough surrounding context to judge
     correctness and maintainability.
   - Verify acceptance criteria and that the test suite is adequate.
   - Validate Clean Architecture boundaries and error handling.

4. Apply review feedback (builder)
   - Address feedback and rerun tests.
   - If changes are high risk or large, ask for a targeted follow-up review.

## Token rules

- Reviewer reads only what is necessary to make a correct call.
- Do not paste full files into reviews.
- Reuse `task_id` for loops (max 3)
