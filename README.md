# AI Dev Team

A portable multi-agent development team that ships features on your command. Works with any AI coding tool — opencode, Claude Code, Cursor, Copilot, Windsurf, or Codex CLI.

## The Team

| Agent | Role |
|-------|------|
| **PM** | Gathers requirements, communicates with you |
| **Architect** | Designs system architecture, writes ADRs |
| **Team Lead** | Breaks work into tickets, manages execution |
| **Developer** | Implements tickets, creates PRs |
| **Reviewer** | Reviews code quality and architecture |
| **QA** | Tests and validates acceptance criteria |
| **Doc Writer** | Updates documentation |

## How It Works

```
/ship "build a login system"
  → PM asks you clarifying questions (you're in the loop)
  → Architect designs the system
  → Team Lead creates tickets with dependency DAG
  → Developers implement in parallel (per DAG layers)
  → Reviewers review PRs (auto-loop for fixes)
  → QA validates (auto-loop for failures)
  → PM asks you to merge PRs
```

You only talk to the PM. Everything else is autonomous.

## Quick Start

```bash
# 1. Clone once
git clone <repo-url> ~/.ai-dev-team

# 2. Install globally
~/.ai-dev-team/ai-dev-team install

# 3. Restart your terminal, then in any project:
cd my-project
ai-dev-team activate

# 4. Ship something
/ship "your feature description"
```

## Commands

| Command | What it does |
|---------|-------------|
| `ai-dev-team install` | Installs globally, adds to PATH |
| `ai-dev-team activate` | Detects AI tools (or asks), creates `.ai/` symlinks, configures provider files |
| `ai-dev-team update` | Pulls latest team updates from git |

## Activation (smart detection)

When you run `ai-dev-team activate` in a project:

1. **Scans** for existing config files: `AGENTS.md`, `CLAUDE.md`, `.cursorrules`, `.windsurfrules`, etc.
2. **Detects** which AI tool you use — or asks if none found
3. **Creates** `.ai/` -> `~/.ai-dev-team/.ai` symlinks (agents, templates, ship.md)
4. **Configures** provider files with the `/ship` command
5. **Updates** `.gitignore` to exclude `.ai/`

Safe to run multiple times — detects existing state and skips what's already set up.

## Provider Detection

| File found | Tool detected |
|------------|-------------|
| `AGENTS.md` | opencode |
| `CLAUDE.md` | Claude Code |
| `.cursorrules` | Cursor |
| `.windsurfrules` | Windsurf |
| `.github/copilot-instructions.md` | GitHub Copilot |
| `CONVENTIONS.md` | Codex CLI |

If multiple are found, it asks which to configure. If none, it prompts you to choose.

## Updating

```bash
ai-dev-team update      # git pull at ~/.ai-dev-team/
```

Updates are instant across all projects since `.ai/` is a symlink pointing to `~/.ai-dev-team/.ai/`.

## Structure

```
~/.ai-dev-team/
├── ai-dev-team           # Main command (in PATH)
├── README.md             # This file
└── .ai/                  # Team resources (symlinked into projects)
    ├── ship.md           # Master workflow (6 phases)
    ├── init.sh           # Setup script
    ├── agents/           # Agent persona definitions
    │   ├── pm.md
    │   ├── architect.md
    │   ├── team-lead.md
    │   ├── developer.md
    │   ├── reviewer.md
    │   ├── qa.md
    │   └── doc-writer.md
    └── templates/        # Ticket and manifest templates
        ├── ticket.md
        └── manifest.md
```

## Provider Agnostic

Uses only markdown files, git commands, and sub-agents. No proprietary APIs or formats. Point any AI coding tool at `.ai/ship.md` and it can execute the workflow.
