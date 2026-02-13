---
description: Create or Update documentation on demand.
---

# Document workflow

Trigger: `/document $1`

Use this workflow to create or update documentation for the requested scope.

## Execution note

You can run this workflow as a single agent. If the task is large, you can
optionally spawn specialized subagents via `task(...)` for specific phases.

## Workflow diagram

This diagram shows the minimal flow for writing docs.

```mermaid
flowchart TD
  A[Start] --> B[Draft (documenter)]
  B --> C[Verify]
  C --> D[Done: docs match codebase]
```

## Steps

1. Draft (documenter)
    - Locate existing docs for the scope, or decide where new docs belong
    - Create or edit Markdown in place
    - Prefer `skill(name="docs-writer")` when writing `.md`

2. Verify
    - Re-read the doc for clarity and consistency
    - Verify details against the codebase (no invented APIs)
    - Fix issues you find

## Exit criteria

- Docs reflect actual behavior (no invented APIs)
- Examples and commands are runnable for this repo
- Links and references are consistent
