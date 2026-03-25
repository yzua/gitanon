# gitanon

Anonymous git identity manager — commit anonymously, impersonate GitHub users, and manage signing behavior per-repo.

```
  gitanon on          Anonymize current repo
  gitanon off         Restore global identity and re-enable signing
  gitanon as <user>   Commit as another GitHub user
  gitanon whoami      Show current repo identity
  gitanon hook        Run gitanon-aware git hooks
```

## What it does

`gitanon` overrides your local git config (never global) so you can:

- **Go anonymous** — erase name/email, disable GPG signing, mark the repo with a config flag
- **Impersonate anyone** — fetch a GitHub user's public profile and commit as them (name + noreply email)
- **Restore cleanly** — undo all overrides and re-enable signing in one command
- **Hook integration** — ships a `mysystem.gitanon` config flag that your hooks can check to skip signing enforcement

## Install

See the full [Installation Guide](docs/installation.md) for all methods including Nix, `go install`, and building from source.

```bash
# Quick install via go
go install github.com/yzua/gitanon@latest

# Or via nix
nix run github:yzua/gitanon

# Or add to your flake inputs
# inputs.gitanon.url = "github:yzua/gitanon";
```

## Quick Start

```bash
# Anonymize current repo
cd my-repo
gitanon on
# ✔ Anonymous mode in my-repo

# Check identity
gitanon whoami
# Repo:     my-repo
# Name:     user
# Email:    (none)
# Signing:  false
# AnonMode: on

# Commit anonymously (no hooks blocking you)
git add .
git commit -m "chore: anonymous work"

# Impersonate another GitHub user
gitanon as octocat
# ✔ Committing as The Octocat <583231+octocat@users.noreply.github.com> in my-repo

# Restore your real identity and re-enable signing
gitanon off
# ✔ Restored global identity in my-repo
```

## Commands

| Command | Aliases | Description |
|---------|---------|-------------|
| `gitanon on` | `anon` | Set anonymous identity (no name, no email, no signing) |
| `gitanon off` | `back`, `undo` | Restore global identity and re-enable GPG signing |
| `gitanon as <user>` | — | Fetch GitHub user and commit as them |
| `gitanon whoami` | — | Show current repo identity |
| `gitanon hook <name>` | — | Run a hook (pre-commit, pre-push) that respects anon mode |
| `gitanon completion <shell>` | — | Generate shell completions (bash, zsh, fish, powershell) |
| `gitanon version` | — | Print version |

## How `gitanon as` works

```
gitanon as octocat
```

1. Calls `GET https://api.github.com/users/octocat` (unauthenticated, public endpoint)
2. Fetches: login, numeric ID, display name
3. Sets local git config:
   - `user.name` = "The Octocat"
   - `user.email` = "583231+octocat@users.noreply.github.com"
   - `commit.gpgSign` = false
   - `mysystem.gitanon` = true

The noreply email format is `<id>+<login>@users.noreply.github.com`, which is GitHub's native anonymous email format.

## Hook Integration

`gitanon` sets a config flag that your hooks can detect:

```bash
# Check if anon mode is active
if [ "$(git config --bool --get mysystem.gitanon || echo false)" = "true" ]; then
  echo "Anonymous mode — skipping signing enforcement"
  exit 0
fi
```

Or use `gitanon hook` directly in your hooks:

```bash
#!/bin/sh
# .git/hooks/pre-push (or global hooksPath)
gitanon hook pre-push
```

See [docs/hooks.md](docs/hooks.md) for detailed setup.

## Config Flag

`gitanon` uses `mysystem.gitanon` in local git config as its marker:

| Value | Meaning |
|-------|---------|
| `true` | Anonymous mode active — hooks should skip signing enforcement |
| absent | Normal mode — hooks enforce signing as usual |

This flag is only set in local config (`git config --local`), never global.

## Development

```bash
just all       # fmt + vet + lint + test + build
just test      # Run tests
just lint      # Run golangci-lint
just cover     # Generate coverage report

nix develop    # Enter dev shell (go, golangci-lint, just, git)
```

## Why not just `git config --local`?

You could do all this manually, but `gitanon`:
- Batches 7 config changes into one command
- Fetches GitHub user IDs for noreply email construction
- Provides a standard flag for hook integration
- Restores signing correctly (the thing you always forget)
- Works the same across shell environments

## License

[WTFPL](LICENSE)
