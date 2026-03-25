# Contributing

## Development Setup

```bash
git clone https://github.com/yzua/gitanon.git
cd gitanon
nix develop    # or install go + just manually
```

## Commands

| Command | Description |
|---------|-------------|
| `just all` | fmt + vet + lint + test + build |
| `just build` | Build the binary |
| `just test` | Run all tests |
| `just lint` | Run golangci-lint |
| `just fmt` | Format Go code |
| `just vet` | Run go vet |
| `just cover` | Generate test coverage report |
| `just install` | Install to `$GOPATH/bin` |
| `just clean` | Remove built binary |

## Project Structure

```
gitanon/
├── main.go              # Entry point
├── cmd/                 # Cobra subcommands
├── internal/
│   ├── git/             # git config operations
│   ├── github/          # GitHub API client
│   └── model/           # Shared types
├── testdata/            # Test fixtures
├── docs/                # Documentation
├── flake.nix            # Nix flake (nix run / nix build)
└── justfile             # Task runner
```

## Adding a New Subcommand

1. Create `cmd/<name>.go` with a `cobra.Command`
2. Register it in `init()` via `rootCmd.AddCommand(<name>Cmd)`
3. Use `RunE` (not `Run`) for error propagation
4. Add tests in `cmd/<name>_test.go` if applicable

## Code Style

- Standard Go formatting (`gofmt`)
- No comments unless needed for non-obvious logic
- Use `RunE` on all commands
- Errors returned from commands print to stderr via cobra

## Testing

```bash
just test
```

Tests use temporary git repos (created with `t.TempDir()`) and do not modify your real config.
