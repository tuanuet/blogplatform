---
description: Master orchestrator for multi-phase development (Token-Optimized).
---

# Develop workflow

Input: $1

Use this workflow for new features when requirements need refinement.

## Roles (subagents)

- `gatekeeper`: clarify requirements and produce a user-approved spec
- `architect`: design contracts and a phased implementation plan
- `builder`: implement with tests (core, then integration)
- `reviewer`: diff-only quality gate (loads `consolidated-review` only)

## How to run subagents (minimal)

- Start: `task(subagent_type="gatekeeper", description="Spec", prompt="...")`
- Continue: `task(task_id="<id>", subagent_type="gatekeeper", description="Spec", prompt="...")`
- Load a skill: `skill(name="<skill>")`
- Ask one blocking question: `question(...)`

## Flow

1. Spec (gatekeeper)
   - Exit: user confirms acceptance criteria and non-goals

2. Design (architect)
   - Exit: contracts + phased plan + test plan

3. Build core (builder)
   - Exit: unit tests green, core behavior implemented

4. Implementation review (reviewer)
   - Exit: `APPROVED` or `NEEDS_CHANGES` (max 3 loops)

5. Integrate (builder)
   - Exit: wiring complete, integration or e2e coverage added as needed

6. Final review (reviewer)
   - Exit: acceptance criteria met, test suite green

## Token rules

- Reviews are diff-only (`git diff` + targeted reads)
- Reviewer loads only `skill(name="consolidated-review")`
- Reuse `task_id` for review loops; avoid re-sending full context
