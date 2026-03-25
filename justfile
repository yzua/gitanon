# gitanon justfile

version := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`
ldflags := "-s -w -X github.com/yzua/gitanon/cmd.Version=" + version
binary := "gitanon"

# Default: show available recipes
default:
    @just --list

# Build the binary
build:
    go build -trimpath -ldflags '{{ ldflags }}' -o {{ binary }} .

# Run all tests
test:
    go test ./... -v -count=1

# Format Go code
fmt:
    gofmt -w .

# Run go vet
vet:
    go vet ./...

# Run golangci-lint
lint:
    golangci-lint run ./...

# Generate test coverage
cover:
    go test ./... -coverprofile=coverage.out
    go tool cover -html=coverage.out -o coverage.html

# Install to $GOPATH/bin
install:
    go install -ldflags '{{ ldflags }}' .

# Clean build artifacts
clean:
    rm -f {{ binary }} coverage.out coverage.html

# Full pipeline: fmt -> vet -> lint -> test -> build
all: fmt vet lint test build

# Quick build + smoke test
check: build
    ./{{ binary }} --help
    ./{{ binary }} --version

# Bump a minor version tag
bump msg:
    git tag -a $(semver next minor $(git tag -l 'v*')) -m "{{ msg }}"
    git push --tags
