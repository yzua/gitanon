# Hook Integration

`gitanon` sets `mysystem.gitanon=true` in local git config when anonymous mode is active. Your hooks can check this flag to skip signing enforcement.

## The Config Flag

```bash
# Check if anon mode is active
git config --bool --get mysystem.gitanon
# Returns: "true" or empty (not set = false)
```

## Option 1: Use `gitanon hook` directly

Wire `gitanon hook` into your git hooks:

- `gitanon hook pre-commit` warns to stderr when `commit.gpgSign` is not `true`, then exits successfully.
- `gitanon hook pre-push` runs `git verify-commit HEAD`; it fails when the current `HEAD` commit is not signed.
- Both commands exit successfully without checks when `mysystem.gitanon=true`.

### Global hooks (recommended)

Set a global hooks directory:

```bash
mkdir -p ~/.config/git/hooks
git config --global core.hooksPath ~/.config/git/hooks
```

Create `~/.config/git/hooks/pre-push`:

```bash
#!/bin/sh
gitanon hook pre-push
```

Create `~/.config/git/hooks/pre-commit`:

```bash
#!/bin/sh
gitanon hook pre-commit
```

Make them executable:

```bash
chmod +x ~/.config/git/hooks/pre-push ~/.config/git/hooks/pre-commit
```

### Per-repo hooks

Same thing, but place scripts in `.git/hooks/` inside your repo.

## Option 2: Check the flag yourself

If you already have hooks, just add a guard at the top:

```bash
#!/bin/sh
# Your existing pre-push hook

# Skip signing enforcement when anon mode is active
if [ "$(git config --bool --get mysystem.gitanon || echo false)" = "true" ]; then
  exit 0
fi

# ... your existing hook logic ...
```

## Custom Full Example: Pre-Push with GPG Enforcement

If you want to verify every pushed commit instead of just `HEAD`, use your own hook logic instead of `gitanon hook pre-push`:

```bash
#!/usr/bin/env bash
set -euo pipefail

# Skip if anon mode
if [ "$(git config --bool --get mysystem.gitanon || echo false)" = "true" ]; then
  echo "⚠ gitanon: anonymous mode, skipping signing checks"
  exit 0
fi

# Verify all pushed commits are signed
zero_sha="0000000000000000000000000000000000000000"
failed=0
input=$(cat)

while IFS= read -r line; do
  [ -z "$line" ] && continue
  read -r local_ref local_sha remote_ref remote_sha <<< "$line"

  [ "$local_sha" = "$zero_sha" ] && continue

  if [ "$remote_sha" = "$zero_sha" ]; then
    range="$local_sha"
  else
    range="${remote_sha}..${local_sha}"
  fi

  for commit in $(git rev-list "$range" 2>/dev/null); do
    if ! git verify-commit "$commit" >/dev/null 2>&1; then
      echo "✗ Unsigned commit: ${commit:0:12} (${local_ref} -> ${remote_ref})"
      failed=1
    fi
  done
done <<< "$input"

if [ "$failed" -ne 0 ]; then
  echo ""
  echo "Push rejected. All commits must have valid GPG signatures."
  echo "Fix: git rebase --exec 'git commit --amend -S --no-edit' <base>"
  exit 1
fi
```

## Shell Completions

Generate completions for your shell:

```bash
# Bash
gitanon completion bash > /etc/bash_completion.d/gitanon

# Zsh
gitanon completion zsh > "${fpath[1]}/_gitanon"

# Fish
gitanon completion fish > ~/.config/fish/completions/gitanon.fish

# PowerShell
gitanon completion powershell > gitanon.ps1
```
