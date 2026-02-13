name: commit-message
description: Generate Git commit messages for this repository using its existing Conventional Commits style (feat/fix/refactor/test/docs/chore, optional scope). Use when the user asks for a commit message, asks you to create a commit, or wants guidance on how to name commits.
---

# Commit Message (Repo Style)

## Goal

Produce a clear, consistent commit message that matches this repo's history.

## Format

Use Conventional Commits:

```
type(scope): short imperative summary

1-2 sentence body explaining why (optional)
```

- **type**: `feat` | `fix` | `refactor` | `test` | `docs` | `chore`
- **scope** (optional but preferred): the affected module/domain, e.g. `auth`, `blog`, `ranking`, `persistence`, `router`, `payment`
- **summary**: imperative verb, present tense; prefer no trailing period; keep it short (~50-72 chars)
- **body**: explain rationale/impact; mention migrations, compatibility, or risk notes when relevant

## How to Choose `type`

- `feat`: new behavior or endpoint
- `fix`: bug fix
- `refactor`: behavior-preserving restructuring
- `test`: adds/fixes tests only
- `docs`: documentation-only changes (including generated swagger files if that is the only change)
- `chore`: tooling, deps, CI, scripts, formatting-only changes

## How to Choose `scope`

Use the most specific stable boundary:

- API area: `auth`, `blog`, `series`, `subscription`, `notification`, `payment`
- Infra: `persistence`, `cache`, `router`, `middleware`, `config`

If changes span multiple areas, either omit scope (`feat:`) or pick the primary user-facing domain.

## Workflow (When Generating a Message)

1. Inspect what will be committed:
   - `git status`
   - `git diff --staged` (or `git diff` if nothing is staged)
   - `git log -n 20 --oneline` to match tone
2. Identify the user-visible outcome and the primary domain (scope).
3. Write the subject line.
4. Add a body only when it helps explain intent, risk, or rollout notes.

## Examples (Seen In This Repo)

```
feat(blog): add slug lookup endpoint

Wire author-scoped slug route and handler; update swagger docs
```

```
fix(persistence): resolve soft delete and upsert issues in tag tier mapping repository
```

```
docs: reorganize blog and ranking system documentation files
```

## Guardrails

- Do not include secrets or environment values in commit messages.
- Prefer describing the intent/outcome over listing every file touched.
- If there are unrelated changes in the working tree, do not mention them unless they are part of the commit.
