# Installation

## Using `go install`

```bash
go install github.com/yzua/gitanon@latest
```

This places `gitanon` in your `$GOPATH/bin` (usually `~/go/bin`).

## Nix

### Run directly (no install)

```bash
nix run github:yzua/gitanon
```

### Add to your flake

```nix
{
  inputs.gitanon.url = "github:yzua/gitanon";

  # Then add it to your packages
  environment.systemPackages = [ inputs.gitanon.packages.${pkgs.system}.default ];
}
```

### Install to profile

```bash
nix profile install github:yzua/gitanon
```

## From Source

```bash
git clone https://github.com/yzua/gitanon.git
cd gitanon
just build
# Binary: ./gitanon
```

Or install to `$GOPATH/bin`:

```bash
just install
```

## Requirements

- Git (any recent version)
- Go 1.23+ (only for building from source)
- [just](https://github.com/casey/just) (only for development)

## Verify Installation

```bash
gitanon --version
gitanon --help
```

## Setup Hook Integration

See [hooks.md](hooks.md) for wiring `gitanon hook` into `pre-commit` and `pre-push`, or for checking `mysystem.gitanon` in your own hooks.
