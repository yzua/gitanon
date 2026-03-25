# AGENTS

Repository-specific guidance for agentic coding tools operating in this project.
Use this file as the default policy unless direct user/system instructions override it.

## Project Snapshot

- Language: Go 1.23 (`go.mod`)
- CLI framework: Cobra (`github.com/spf13/cobra`)
- Binary: `gitanon`
- Entry point: `main.go`
- Task runner: `just`
- Optional dev environment: Nix (`flake.nix`)

## Directory Intent

- `cmd/`: CLI command definitions, argument validation, user-facing output
- `internal/git/`: git config/process operations
- `internal/github/`: GitHub API client logic
- `internal/model/`: shared structs for git/github data
- `docs/`: user and contributor docs

Design rule: keep `cmd/` thin. Business behavior belongs in `internal/*` packages.

## Build, Lint, Test Commands

Canonical commands from `justfile`:

```bash
just build      # build binary with ldflags version metadata
just test       # go test ./... -v -count=1
just fmt        # gofmt -w .
just vet        # go vet ./...
just lint       # golangci-lint run ./...
just cover      # generate coverage.out and coverage.html
just check      # build + smoke run --help and --version
just all        # fmt -> vet -> lint -> test -> build
```

Fallback when local Go env config causes trouble:

```bash
GOENV=off go test ./...
GOENV=off go vet ./...
GOENV=off golangci-lint run ./...
GOENV=off go build ./...
```

## Running a Single Test (Important)

There is no `just` recipe for one test; use `go test` with `-run`.

Run one package:

```bash
go test ./internal/git -count=1 -v
```

Run one top-level test:

```bash
go test ./internal/git -run '^TestAnonymizeAndRestore$' -count=1 -v
```

Run one subtest:

```bash
go test ./internal/github -run 'TestLookupUser/DisplayName with name' -count=1 -v
```

Run a pattern across all packages:

```bash
go test ./... -run 'TestRepoName|TestIsInsideRepo' -count=1 -v
```

If needed, prefix any command with `GOENV=off`.

## Validation Order After Edits

Use the narrowest checks first, then full validation:

```bash
gofmt -w <changed-files>
go test ./<changed-package> -count=1
go test ./... -count=1
go vet ./...
golangci-lint run ./...
go build ./...
```

For a full local CI-like pass, run `just all`.

## Code Style and Conventions

### Formatting

- Always `gofmt` changed Go files.
- Keep code/doc edits ASCII by default; only use Unicode when already established.
- Do not add comments unless logic is not obvious.

### Imports

- Use standard Go import layout as produced by `gofmt`.
- Keep stdlib imports first, then a blank line, then external/internal imports.
- Alias imports only when clarity improves (e.g. `gh` for `internal/github`).

### Types and Structs

- Put shared cross-package data in `internal/model`.
- Keep structs focused and small (`GitUser`, `GitHubUser` pattern).
- Use JSON tags only where data maps to external API payloads.
- Keep package-private helper structs local to package files.

### Naming

- Exported names: PascalCase (`LookupUser`, `RepoName`).
- Unexported names: lowerCamelCase (`readUser`, `setLocal`).
- Package-private constants: lowerCamelCase (`keyUserName`, `configTrue`).
- Test functions: `TestXxx`; use `t.Run("case", ...)` for subtests.
- Prefer clear words over novel abbreviations (except idiomatic `err`, `cmd`, `ID`).

### Error Handling

- Return errors from library code; avoid printing in `internal/*`.
- In Cobra commands, prefer `RunE` for any path that can fail.
- Wrap errors with context using `%w` (`fmt.Errorf("setting identity: %w", err)`).
- Use sentinel errors (`errors.New`) only when callers branch on them.
- For defer cleanup where failure is non-actionable, ignore safely (`_ = resp.Body.Close()`).

### CLI Command Patterns

- Define each subcommand in `cmd/<name>.go` as `var <name>Cmd = &cobra.Command{...}`.
- Register in `init()` with `rootCmd.AddCommand(...)`.
- Validate args with Cobra helpers (`cobra.ExactArgs(1)`, etc.).
- Use `requireRepo()` before repo-dependent operations.
- Keep user-visible output concise via `fmt.Printf`/`fmt.Println`.

### Testing Practices

- Prefer black-box package tests (`package git_test`, not `package git`).
- Use `t.TempDir()` for git repo fixtures and filesystem isolation.
- Mark helper funcs with `t.Helper()`.
- Use `t.Fatalf` for setup failures; `t.Errorf`/`t.Error` for assertions.
- Keep tests deterministic and independent of host global git config.

## Cursor and Copilot Rules

Checked repository-specific rule files:

- `.cursor/rules/`: not present
- `.cursorrules`: not present
- `.github/copilot-instructions.md`: not present

If these files are added later, incorporate their guidance here and follow them as project policy.
