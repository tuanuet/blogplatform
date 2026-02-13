---
description: Create or Update documentation on demand.
---

# Document workflow

Trigger: `/document $1`

Use this workflow to create or update documentation for the requested scope.

## Steps

1. Draft (documenter)
   - Locate existing docs for the scope, or decide where new docs belong
   - Create or edit Markdown in place
   - Prefer `skill(name="docs-writer")` when writing `.md`

2. Review (document-reviewer)
   - Verify the doc against the codebase
   - Return `APPROVED` or `NEEDS_CHANGES`

3. Loop
   - If `NEEDS_CHANGES`, apply feedback and re-submit
   - Stop after 3 loops and escalate unresolved issues

## Exit criteria

- Docs reflect actual behavior (no invented APIs)
- Examples and commands are runnable for this repo
- Links and references are consistent
