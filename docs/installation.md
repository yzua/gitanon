# Installation

## From Source

```bash
git clone https://github.com/yzua/gitanon.git
cd gitanon
make build
# Binary: ./gitanon
```

Or install to `$GOPATH/bin`:

```bash
make install
```

## Using `go install`

```bash
go install github.com/yzua/gitanon@latest
```

This places `gitanon` in your `$GOPATH/bin` (usually `~/go/bin`).

## Pre-built Binaries

Download from [Releases](https://github.com/yzua/gitanon/releases) and place in your `$PATH`.

## Requirements

- Git (any recent version)
- Go 1.23+ (only for building from source)

## Verify Installation

```bash
gitanon --version
gitanon --help
```

## Setup Hook Integration

See [docs/hooks.md](hooks.md) for adding the `mysystem.gitanon` flag check to your git hooks.
