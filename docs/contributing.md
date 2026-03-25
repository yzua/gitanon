# Contributing

## Development Setup

```bash
git clone https://github.com/yzua/gitanon.git
cd gitanon
```

## Make Targets

| Target | Description |
|--------|-------------|
| `make all` | fmt + vet + lint + test + build |
| `make build` | Build the binary |
| `make test` | Run all tests |
| `make lint` | Run golangci-lint |
| `make fmt` | Format Go code |
| `make vet` | Run go vet |
| `make cover` | Generate test coverage report |
| `make install` | Install to `$GOPATH/bin` |
| `make clean` | Remove built binary |

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
└── docs/                # Documentation
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
make test
```

Tests use temporary git repos (created with `t.TempDir()`) and do not modify your real config.
