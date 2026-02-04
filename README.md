# Gitflow

Gitflow is a professional command line tool that turns everyday Git workflows into single intentional commands.

Remove friction from your daily development work.

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

## Key features

- **Status and diagnostics**: Repository health checks, including `gitflow status` and `gitflow doctor`.
- **Branch management**: Start, sync, commit, clean up, and inspect branches safely.
- **Pull request workflows**: Create, list, and view PRs via GitHub or GitLab.
- **Release system**: Deterministic versioning, changelogs, tags, and publishing.
- **CI-friendly output**: JSON, env, no-color, and non-interactive modes.

---

## Installation

Clone the repository

```bash
git clone https://github.com/C-NASIR/gitflow
cd gitflow
```

Build the binary

```bash
go build -o gitflow .
```

Or install globally

```bash
go install .
```

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

## Commands

### General

- `gitflow version` prints the gitflow build version.
- `gitflow status` shows the repository status summary.
- `gitflow doctor` runs diagnostics without mutating the repo.
- `gitflow init` writes a starter `.gitflow.yml` file.
- `gitflow config show` prints the resolved configuration.
- `gitflow config validate` validates configuration and reports errors.

### Branches and sync

- `gitflow start <name>` starts a new branch using conventions.
- `gitflow sync` syncs the current branch with the base branch.
- `gitflow commit` creates a commit using conventions or prompts.
- `gitflow cleanup` deletes merged or stale branches safely.
- `gitflow branch list` lists local branches with age and ahead/behind.

### Pull requests

- `gitflow pr create` creates a pull request for the current branch.
- `gitflow pr list` lists pull requests from the provider.
- `gitflow pr view <number>` shows a pull request by number.

### Providers

- `gitflow provider check` validates provider credentials and access.

### Releases

- `gitflow release preview` previews the next version and changelog.
- `gitflow release create` creates an annotated tag with changelog.
- `gitflow release changelog` outputs the changelog since last release.
- `gitflow release version` prints the next release version.
- `gitflow release publish` publishes release notes to the provider.

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

## License

MIT License

Use it.
Learn from it.
Extend it.
Ship better tools.
