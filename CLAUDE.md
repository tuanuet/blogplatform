# AI Agent Development Pipeline

> **3-Phase Pipeline**: Every development task MUST go through 3 phases.

## Core Principle

```
GATEKEEPER → ARCHITECT → BUILDER
(Refinement)   (Design)    (TDD)
```

## Quick Commands

| Command        | Description             |
| -------------- | ----------------------- |
| `/pipeline`    | Full 3-phase workflow   |
| `/new-feature` | New feature development |
| `/bug-fix`     | Bug fix workflow        |

## Agents

See [AGENTS.md](./AGENTS.md) for details on:

- **Orchestrator**: Pipeline coordinator
- **Gatekeeper**: Validate & refine requirements
- **Architect**: Design schema & API contracts
- **Builder**: TDD implementation

## Skills

Skills are located in `.agent/skills/`. Each agent uses appropriate skills:

- `tech-stack-detect` - Auto-detect from codebase
- `design-patterns` - SOLID, DDD, Clean Architecture
- `clean-code` - Code quality standards
- `testing` - Testing strategies
