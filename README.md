# Gitflow

Gitflow is a professional command line tool that turns everyday Git workflows into single intentional commands.

This project exists to remove friction from daily development work while teaching how to build a real production quality CLI in Go.

---

## Why this project exists

Modern developers repeat the same Git actions constantly.

```
- Switch branch
- Pull latest changes
- Create a feature branch
- Push and set upstream
- Open a pull request
- Clean up old branches
- Prepare releases
- Generate changelogs
```

Each task is simple in isolation, but together they create constant cognitive overhead.

Gitflow solves this by encoding team conventions and best practices into a single tool that understands intent instead of raw Git commands.

---

## What Gitflow does

Gitflow is a smart workflow manager built on top of Git.

It provides

1. High level workflow commands instead of low level Git invocations
2. Safe automation with strong defaults and guardrails
3. A consistent user interface across all commands
4. Deterministic release and versioning logic
5. Optional GitHub and GitLab integration
6. CI friendly behavior for automation and pipelines

---

## Core principles

### Intent over mechanics

You say what you want to do.  
Gitflow figures out how to do it safely.

Example

Instead of manually running several Git commands to start a feature branch, you run

```
gitflow start user auth
```

Gitflow handles base branch selection, cleanliness checks, branch naming, and pushing.

### Safety first

Gitflow never rewrites history silently.  
It refuses dangerous operations by default.  
Destructive actions require explicit confirmation.

### Deterministic behavior

Given the same repository state, Gitflow always produces the same result.  
This is critical for releases and CI usage.

### Cohesive UX

Every command uses the same UI system.  
Tables look the same.  
Headers look the same.  
Errors are consistent and actionable.

---

## Project structure

The project is organized around clear boundaries.

cmd contains CLI commands only  
internal contains implementation details  
pkg contains shared domain types

High level architecture

CLI commands call workflows  
Workflows orchestrate Git, config, providers, and UI  
Git helpers wrap Git commands safely  
Providers abstract GitHub and GitLab APIs

This separation makes the project testable, extensible, and easy to reason about.

---

## Key features

### Status and diagnostics

```
gitflow status
```

Shows branch, working tree state, and configuration source.

```
gitflow doctor
```

Diagnoses common issues without mutating the repository.

Checks include

1. Is this a Git repository?
2. Is the working tree clean?
3. Is configuration present?
4. Is configuration valid?
5. Are provider credentials available?

---

### Branch management

```
gitflow start
```

Creates a new branch using configured naming conventions.

```
gitflow branch list
```

Lists local branches with age and ahead behind analytics.

```
gitflow cleanup
```

Safely deletes merged or stale branches with interactive selection.

---

### Pull request workflows

```
gitflow pr create
```

Creates pull requests via GitHub or GitLab APIs.

```
gitflow pr list
```

Lists pull requests with consistent table output.

```
gitflow pr view
```

Shows detailed pull request information.

All provider interactions are optional and validated.

---

### Release system

Gitflow includes a complete deterministic release workflow.

```
gitflow release preview
```

Computes the next version and changelog without side effects.

```
gitflow release create
```

Creates annotated Git tags safely.

```
gitflow release changelog
```

Generates changelog output only.

Release logic is based entirely on commit history and conventional commits.

---

### CI and automation support

Gitflow is safe to run in CI environments.

```
gitflow release version
```

Outputs the next version as a single line.

Commands support

1. JSON output
2. Environment variable output
3. No color mode
4. Non interactive operation

This allows Gitflow to be used in pipelines for versioning, artifacts, and releases.

---

## Configuration

Gitflow uses a YAML configuration file named `.gitflow.yml`.

Configuration is optional.  
Defaults are provided for everything.

Example configuration

```yaml
provider:
  type: github
  token_env: GITHUB_TOKEN
  owner: myorg
  repo: myrepo

branches:
  feature_prefix: feature/
  main_branch: main

workflows:
  start:
    base_branch: main
    auto_push: true

  cleanup:
    merged_only: true
    age_threshold_days: 30
    protected_branches:
      - main

ui:
  color: true
  emoji: false
  verbose: false
```

You can generate a starter config using

```
gitflow init
```

---

## UI customization

UI behavior can be controlled via config, environment variables, or CLI flags.

Precedence order

1. CLI flags
2. Environment variables
3. Config file
4. Defaults

Supported controls include

Color output
Emoji usage
Verbose logging

This makes Gitflow usable both locally and in CI.

---

## Testing philosophy

This project is heavily tested.

Tests include

1. Unit tests for Git helpers
2. Workflow level tests using temporary Git repositories
3. Provider tests using mocked HTTP servers
4. Deterministic release computation tests

All tests run with

```
go test ./...
```

No network access is required.

---

## Installation

Clone the repository

```bash
git clone https://github.com/C-NASIR/gitflow
cd gitflow
```

Build the binary

```bash
go build -o gitflow cmd/gitflow/main.go
```

Or install globally

```bash
go install run main.go
```

---

## Daily usage examples

```bash
gitflow status
gitflow start user auth
gitflow pr create
gitflow branch list
gitflow cleanup
gitflow release preview
gitflow release create
```

---

## License

MIT License

Use it.
Learn from it.
Extend it.
Ship better tools.
