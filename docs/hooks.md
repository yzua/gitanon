# Hook Integration

`gitanon` sets a `mysystem.gitanon=true` flag in local git config when anonymous mode is active. Your hooks can check this flag to conditionally skip signing enforcement.

## The Config Flag

```bash
# Check if anon mode is active
git config --bool --get mysystem.gitanon
# Returns: "true" or empty (not set = false)
```

## Option 1: Use `gitanon hook` directly

Wire `gitanon hook` into your git hooks:

### Global hooks (recommended)

Set a global hooks directory:

```bash
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

## Full Example: Pre-Push with GPG Enforcement

```bash
#!/bin/sh
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
```
