---
name: ckb-code-scan
description: Use Code Knowledge Base (CKB) tools for semantic code scanning, understanding symbols, analyzing relationships, and impact analysis instead of basic grep/glob/read operations
---

# CKB Code Scanning

## Purpose

CKB (Code Knowledge Backend) transforms your codebase into a queryable knowledge base. Ask questions, understand impact, find owners, detect dead code—all through semantic analysis tools.

> Think of it as a senior engineer who knows every line of code, every decision, and every owner—available 24/7 to answer your questions.

## What CKB Can Do

| Question | Without CKB | With CKB |
|----------|-------------|----------|
| "What breaks if I change this?" | Grep and hope | Precise blast radius with risk score |
| "What tests should I run?" | Run everything (30 min) | Run affected tests only (2 min) |
| "How does this system work?" | Read code for hours | Query architecture instantly |
| "Is this code still used?" | Delete and see what breaks | Confidence-scored dead code detection |
| "Who owns this code?" | Search CODEOWNERS manually | Ownership with drift detection |

## When to Use

**USE CKB tools for:**
- Understanding what a function/class does (symbol explanation)
- Finding where a symbol is used (reference finding)
- Tracing call relationships (call graph)
- Analyzing code architecture and dependencies
- Impact analysis before changing code
- Exploring modules or directories comprehensively
- Finding related code across the codebase
- Detecting dead code
- Security scanning (secrets, credentials)

**USE basic tools for:**
- Simple file existence checks (`glob`)
- Reading entire file contents (`read`)
- Basic text pattern matching (`grep`) when semantic analysis is not needed

## Setup and Initialization

### Install CKB

```bash
# Install globally (recommended)
npm install -g @tastehub/ckb

# Or run directly with npx (no install needed)
npx @tastehub/ckb init
```

### Initialize in Project

```bash
# 1. Navigate to your project
cd /path/to/your/project

# 2. Initialize CKB
npx @tastehub/ckb init

# 3. Generate SCIP index (auto-detects language)
npx @tastehub/ckb index

# 4. Check status
npx @tastehub/ckb status
```

### Auto-refresh for AI Sessions

```bash
# MCP watch mode - auto-reindexes every 30s when stale
npx @tastehub/ckb mcp --watch
```

### MCP Integration Setup

**For Claude Code:**
```bash
npx @tastehub/ckb setup
```

Or manually add to `.mcp.json`:
```json
{
  "mcpServers": {
    "ckb": {
      "command": "npx",
      "args": ["@tastehub/ckb", "mcp"]
    }
  }
}
```

**For OpenCode:**
Add to `opencode.json`:
```json
{
  "mcp": {
    "ckb": {
      "type": "local",
      "command": ["npx", "@tastehub/ckb", "mcp"],
      "enabled": true
    }
  }
}
```

## Presets for Token Optimization

CKB exposes 80+ tools. Use presets to reduce token overhead by up to 83%:

```bash
# List all presets
npx @tastehub/ckb mcp --list-presets

# Core preset (14 essential tools) - DEFAULT
npx @tastehub/ckb mcp --preset=core

# Workflow-specific presets
npx @tastehub/ckb mcp --preset=review      # 19 tools - core + diff, ownership
npx @tastehub/ckb mcp --preset=refactor    # 19 tools - core + coupling, dead code
npx @tastehub/ckb mcp --preset=full        # 80+ tools - all tools
```

## CKB Tool Selection Guide

| Need                                    | Tool                          | Use Case Example |
| --------------------------------------- | ----------------------------- | ---------------- |
| Quick overview of file/directory        | `ckb_explore`                 | "What's in src/auth/?" |
| Deep dive into a specific function     | `ckb_understand`              | "How does validateToken work?" |
| Find all usages of a symbol            | `ckb_findReferences`          | "Where is createUser called?" |
| See call relationships (callers/callees) | `ckb_getCallGraph`            | "What functions does processData call?" |
| Analyze impact before change           | `ckb_prepareChange`            | "What breaks if I rename User?" |
| Get architecture/dependencies          | `ckb_getArchitecture`          | "Show me module dependencies" |
| Search for symbols by name             | `ckb_searchSymbols`            | "Find all database models" |
| Trace usage from entrypoints           | `ckb_traceUsage`              | "How is this config reached?" |
| Get explanation of symbol              | `ckb_explainSymbol`            | "What does this function do?" |
| Batch retrieve symbols                | `ckb_batchGet`                | "Get 50 symbols at once" |
| Batch search queries                  | `ckb_batchSearch`             | "Run 10 searches at once" |

## Tool Quick Reference

### `ckb_explore` - Area Exploration

**Best for:** Quick overview of structure, symbols, and hotspots

```bash
# Quick overview
ckb_explore target="src/auth/"

# Deep dive with focus on dependencies
ckb_explore target="src/auth/" depth="deep" focus="dependencies"

# Focus on recent changes
ckb_explore target="src/auth/" focus="changes"

# Shallow exploration
ckb_explore target="src/auth/" depth="shallow"
```

### `ckb_understand` - Symbol Deep Dive

**Best for:** Understanding a specific function, class, or module

```bash
# Understand a function completely
ckb_understand query="validateToken"

# With references and call graph
ckb_understand query="AuthService.login" includeReferences=true includeCallGraph=true

# Limit reference count
ckb_understand query="UserService" includeReferences=true maxReferences=50
```

### `ckb_searchSymbols` - Semantic Search

**Best for:** Finding symbols by name or type (more accurate than grep)

```bash
# Search for functions
ckb_searchSymbols query="createUser" kinds=["function"]

# Search for database models
ckb_searchSymbols query="User" kinds=["class"]

# Search within a module
ckb_searchSymbols query="validate" scope="src/auth"

# Limit results
ckb_searchSymbols query="handler" limit=20
```

### `ckb_findReferences` - Find Usages

**Best for:** Finding where a symbol is used across codebase

```bash
# Find all references
ckb_findReferences symbolId="ckb:<repo>:sym:<fingerprint>"

# Limited results
ckb_findReferences symbolId="..." limit=50

# With merge strategy
ckb_findReferences symbolId="..." merge="union"
```

### `ckb_getCallGraph` - Call Relationships

**Best for:** Understanding call flow and dependencies

```bash
# Both directions (default)
ckb_getCallGraph symbolId="..." direction="both" depth=1

# Only callers
ckb_getCallGraph symbolId="..." direction="callers" depth=2

# Only callees (what this calls)
ckb_getCallGraph symbolId="..." direction="callees" depth=3

# Deep traversal
ckb_getCallGraph symbolId="..." direction="both" depth=4
```

### `ckb_prepareChange` - Impact Analysis

**Best for:** Before modifying, renaming, or deleting code

```bash
# Before modifying
ckb_prepareChange target="AuthService.login" changeType="modify"

# Before renaming
ckb_prepareChange target="User" changeType="rename"

# Before deleting
ckb_prepareChange target="deprecated_function" changeType="delete"

# Before extracting
ckb_prepareChange target="LegacyCode" changeType="extract"
```

### `ckb_getArchitecture` - Module Dependencies

**Best for:** High-level architecture view

```bash
# Module level
ckb_getArchitecture granularity="module"

# Directory level with metrics
ckb_getArchitecture granularity="directory" includeMetrics=true depth=3

# Focus on specific path
ckb_getArchitecture targetPath="src/api" granularity="file"

# Include external dependencies
ckb_getArchitecture granularity="module" includeExternalDeps=true
```

### `ckb_traceUsage` - Usage Tracing

**Best for:** Understanding how code is reached from entrypoints

```bash
# Trace usage paths
ckb_traceUsage symbolId="..." maxDepth=5 maxPaths=10

# Shallow trace
ckb_traceUsage symbolId="..." maxDepth=3 maxPaths=5
```

### `ckb_explainSymbol` - Symbol Explanation

**Best for:** Getting AI-friendly explanation including usage, history, and summary

```bash
# Get explanation
ckb_explainSymbol symbolId="ckb:<repo>:sym:<fingerprint>"
```

### `ckb_batchGet` - Batch Symbol Retrieval

**Best for:** Retrieving multiple symbols at once (max 50)

```bash
# Get multiple symbols
ckb_batchGet symbolIds=["ckb:<repo>:sym:<fp1>", "ckb:<repo>:sym:<fp2>", ...]
```

### `ckb_batchSearch` - Batch Search

**Best for:** Running multiple searches at once (max 10)

```bash
# Run multiple searches
ckb_batchSearch queries=[{"query": "User", "kinds": ["class"]}, {"query": "create", "kinds": ["function"]}]
```

## Language Support

CKB classifies languages into **quality tiers** based on indexer maturity:

| Tier | Quality | Languages |
|------|---------|-----------|
| **Tier 1** | Full support, all features | Go |
| **Tier 2** | Full support, minor edge cases | TypeScript, JavaScript, Python |
| **Tier 3** | Basic support, call graph may be incomplete | Rust, Java, Kotlin, C++, Ruby, Dart |
| **Tier 4** | Experimental | C#, PHP |

**Key limitations:**
- **Incremental indexing** is Go-only. Other languages require full reindex.
- **TypeScript monorepos** may need `--infer-tsconfig` flag
- **C/C++** requires `compile_commands.json`
- **Python** works best with activated virtual environment

Run `ckb doctor --tier standard` to check if your language tools are properly installed.

## Common Workflows

### Workflow 1: Understand Existing Feature

```
1. ckb_explore target="src/feature/" - Get overview
2. ckb_understand query="MainFunction" - Deep dive
3. ckb_getCallGraph symbolId="..." direction="both" - See relationships
4. ckb_findReferences symbolId="..." - Check all usages
```

### Workflow 2: Impact Analysis Before Change

```
1. ckb_prepareChange target="TargetSymbol" changeType="modify"
2. Review affected tests, coupled files, risk score
3. Use ckb_findReferences for complete usage list
4. Plan changes accordingly
```

### Workflow 3: Find Related Code

```
1. ckb_searchSymbols query="pattern" kinds=["function","class"]
2. Select relevant symbols
3. ckb_batchGet symbolIds=[...] - Get details in batch
4. ckb_getCallGraph to see relationships
```

### Workflow 4: Architecture Review

```
1. ckb_getArchitecture granularity="module" - High-level view
2. ckb_getArchitecture granularity="directory" includeMetrics=true - Detailed view
3. ckb_explore target="complex-module/" - Deep dive into problem areas
```

### Workflow 5: Pre-Implementation Analysis

```
1. ckb_searchSymbols query="similarFeature" - Find existing patterns
2. ckb_understand query="similarFeature" - Understand implementation
3. ckb_getArchitecture - Ensure consistency
4. Use patterns for new implementation
```

## Comparison: CKB vs Basic Tools

| Task               | Basic Tools (grep/glob/read)              | CKB Tools                              |
| ------------------ | ----------------------------------------- | -------------------------------------- |
| Find function calls| `grep -r "functionName"` - text matching  | `ckb_findReferences` - semantic links |
| Understand code    | Read file + manual analysis               | `ckb_understand` - explanation + usage|
| See relationships  | Manual tracing                            | `ckb_getCallGraph` - automatic graph   |
| Impact analysis    | Manual search + assumptions               | `ckb_prepareChange` - blast radius    |
| Explore codebase   | Multiple `ls` and `read` calls             | `ckb_explore` - comprehensive view    |
| Find symbols       | `grep` - prone to false positives         | `ckb_searchSymbols` - semantic search |
| Architecture view  | Manual diagramming                        | `ckb_getArchitecture` - auto-generated |
| What breaks if...  | Guesswork                                  | Precise blast radius with risk score   |

## Best Practices

1. **Start with exploration** - Use `ckb_explore` for overview before diving deep
2. **Batch operations** - Use `ckb_batchGet` and `ckb_batchSearch` for efficiency
3. **Always impact analyze** - Use `ckb_prepareChange` before modifications
4. **Leverage call graphs** - `ckb_getCallGraph` reveals hidden dependencies
5. **Trace from entrypoints** - `ckb_traceUsage` shows execution paths
6. **Use compound operations** - `ckb_understand`, `ckb_explore`, `ckb_prepareChange` reduce tool calls by 60-70%
7. **Keep index fresh** - Use `--watch` mode during development
8. **Check status** - Use `ckb status` to verify index is up to date

## Symbol IDs

Most CKB tools require a stable symbol ID with format: `ckb:<repo>:sym:<fingerprint>`

Get symbol IDs from:
- `ckb_searchSymbols` results
- `ckb_explore` results
- `ckb_understand` results
- `ckb_batchGet` results

## Integration in Agent Workflows

### Gatekeeper Agent

**Use CKB for:**
- Understanding project structure during requirement analysis
- Finding relevant code for feature context
- Identifying existing patterns

```
Workflow:
1. ckb_explore target="src/" - Get project overview
2. ckb_getArchitecture granularity="module" - Understand module structure
3. ckb_searchSymbols query="relatedFeature" - Find similar implementations
4. Use context to generate accurate Refined Spec
```

### Architect Agent

**Use CKB for:**
- Understanding existing patterns before designing
- Analyzing module dependencies
- Ensuring consistency with existing architecture

```
Workflow:
1. ckb_getArchitecture - Understand module dependencies
2. ckb_searchSymbols query="Model" kinds=["class"] - Find existing models
3. ckb_understand query="ExistingModel" - Understand schema patterns
4. Design new schema/API following existing conventions
```

### Builder Agent

**Use CKB for:**
- Impact analysis before any modifications
- Understanding existing code before writing tests
- Locating all test files related to a feature
- Understanding execution flow

```
Workflow:
1. ckb_prepareChange target="TargetSymbol" changeType="modify" - Impact analysis
2. ckb_understand query="FunctionToTest" - Understand implementation
3. ckb_findReferences symbolId="..." - Locate test files
4. ckb_traceUsage - Understand execution flow
5. Write TDD tests based on understanding
```

## Additional Resources

- [Features Guide](https://github.com/SimplyLiz/CodeMCP/wiki/Features) - Complete feature list
- [Prompt Cookbook](https://github.com/SimplyLiz/CodeMCP/wiki/Prompt-Cookbook) - Real prompts for real problems
- [Integration Guide](https://github.com/SimplyLiz/CodeMCP/wiki/Integration-Guide) - Use CKB in your own tools
- [Impact Analysis](https://github.com/SimplyLiz/CodeMCP/wiki/Impact-Analysis) - Blast radius, affected tests
- [Practical Limits](https://github.com/SimplyLiz/CodeMCP/wiki/Practical-Limits) - Accuracy notes, blind spots
