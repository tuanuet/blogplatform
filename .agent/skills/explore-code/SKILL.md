---
name: explore-code
description: Core capability for navigating and understanding codebases using Serena MCP. Use this skill when you need to find symbols, explore directories, understand code logic, or analyze impact. It replaces raw file searching with semantic code intelligence.
---

# Explore Code

This skill provides semantic code intelligence capabilities through the Serena MCP server. It allows you to navigate the codebase like an IDE, finding definitions, references, and understanding code structure without relying on brittle text searches.

## When to Use

Use this skill for:
- **Navigation**: Finding where classes, functions, or variables are defined.
- **Exploration**: Understanding the structure of a module or directory.
- **Impact Analysis**: Finding all usages of a symbol before making changes.
- **Code Understanding**: Getting high-level summaries or deep-dives into specific symbols.

## Available Tools (Serena MCP)

This skill utilizes the following Serena MCP tools. **Prefer these over `grep` or `find` for code tasks.**

### 1. File & Directory Overview

- **`serena_list_dir`**: List files and directories.
  - *Usage*: `serena_list_dir relative_path="src/components" recursive=false`
  - *Best for*: Getting a high-level view of project structure.

- **`serena_get_symbols_overview`**: List all symbols (classes, functions, variables) in a file.
  - *Usage*: `serena_get_symbols_overview relative_path="src/utils/helpers.ts"`
  - *Best for*: Understanding what a file contains without reading the entire content.

### 2. Finding Symbols (Definitions)

- **`serena_find_symbol`**: Find symbols by name or pattern.
  - *Usage*: `serena_find_symbol name_path_pattern="AuthService" include_body=true`
  - *Options*:
    - `include_body=true`: Returns the source code of the symbol.
    - `include_info=true`: Returns docstrings and signatures.
    - `substring_matching=true`: Matches partial names (e.g., "Auth" matches "AuthService").
  - *Best for*: Locating specific classes, functions, or constants.

### 3. Finding References (Usages)

- **`serena_find_referencing_symbols`**: Find all usages of a symbol.
  - *Usage*: `serena_find_referencing_symbols name_path="User" relative_path="src/models/User.ts"`
  - *Best for*: Impact analysis—seeing what will break if you change a symbol.

### 4. Code Search

- **`serena_search_for_pattern`**: Semantic/Regex search across the codebase.
  - *Usage*: `serena_search_for_pattern substring_pattern="TODO|FIXME" paths_include_glob="src/**/*.ts"`
  - *Best for*: Finding specific code patterns or comments when symbol search isn't enough.

## Workflow Examples

### Workflow 1: Understanding a Feature

1.  **Explore the directory**:
    ```
    serena_list_dir relative_path="src/features/auth"
    ```
2.  **Get symbols in key files**:
    ```
    serena_get_symbols_overview relative_path="src/features/auth/AuthService.ts"
    ```
3.  **Read the main logic**:
    ```
    serena_find_symbol name_path_pattern="AuthService" include_body=true
    ```

### Workflow 2: Impact Analysis (Before Changing Code)

1.  **Locate the symbol**:
    ```
    serena_find_symbol name_path_pattern="calculateTotal"
    ```
2.  **Find all references**:
    ```
    serena_find_referencing_symbols name_path="calculateTotal" relative_path="src/utils/calc.ts"
    ```
3.  **Review usages**: Check if the change will break any call sites.

### Workflow 3: Finding Related Code

1.  **Search for similar names**:
    ```
    serena_find_symbol name_path_pattern="*Controller" substring_matching=true
    ```
