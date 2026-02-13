---
description: Analyze the current session and extract any patterns worth saving as skills.
---

# Learn workflow

Trigger: `/learn`

Use this workflow after a non-trivial session to save one reusable pattern.

## Steps

1. Identify one pattern
   - Error resolution pattern, debugging technique, workaround, or convention
   - Avoid one-off incidents and trivial fixes

2. Draft a short note
   - Save to: `.opencode/skills/learned/<pattern-name>.md`
   - Include: context, problem, solution, example, when to use

3. Confirm and save
   - Ask for confirmation via `question` before writing the final file

## Exit criteria

- One focused note that is reusable in future sessions
