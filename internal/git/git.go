// Package git provides helpers for reading and writing git config.
package git

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/yzua/gitanon/internal/model"
)

var ErrNotInRepo = errors.New("not inside a git repository")

// Run executes git with the given arguments.
func Run(args ...string) error {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), strings.TrimSpace(string(out)), err)
	}
	return nil
}

// Get reads a git config value. Returns empty string if unset.
func Get(scope, key string) string {
	args := []string{"config"}
	if scope != "" {
		args = append(args, scope)
	}
	args = append(args, "--get", key)
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// Set writes a local git config value.
func Set(key, value string) error {
	return Run("config", "--local", key, value)
}

// Unset removes a local git config key (no error if missing).
func Unset(key string) {
	_ = Run("config", "--local", "--unset", key)
}

// IsInsideRepo checks if we're in a git work tree.
func IsInsideRepo() bool {
	err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Run()
	return err == nil
}

// RepoName returns the basename of the repository root.
func RepoName() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "(unknown)"
	}
	return filepath.Base(strings.TrimSpace(string(out)))
}

// WhoAmI reads the current git identity.
func WhoAmI() model.GitUser {
	user := model.GitUser{
		Name:     Get("--local", "user.name"),
		Email:    Get("--local", "user.email"),
		SignKey:  Get("--local", "user.signingKey"),
		IsLocal:  Get("--local", "user.name") != "",
		AnonMode: Get("--local", "mysystem.gitanon") == "true",
	}
	signStr := Get("--local", "commit.gpgSign")
	user.Signing = signStr == "true"

	// Fall back to global if no local override
	if user.Name == "" && !user.IsLocal {
		user.Name = Get("--global", "user.name")
		user.Email = Get("--global", "user.email")
		user.SignKey = Get("--global", "user.signingKey")
		signStr = Get("--global", "commit.gpgSign")
		user.Signing = signStr == "true"
		user.IsLocal = false
	}
	return user
}

// Anonymize sets the repo to anonymous mode (no identity, no signing).
func Anonymize(name, email string) error {
	if err := Set("user.name", name); err != nil {
		return err
	}
	if err := Set("user.email", email); err != nil {
		return err
	}
	if err := Set("commit.gpgSign", "false"); err != nil {
		return err
	}
	if err := Set("tag.gpgSign", "false"); err != nil {
		return err
	}
	if err := Set("push.gpgSign", "false"); err != nil {
		return err
	}
	if err := Set("user.signingKey", ""); err != nil {
		return err
	}
	return Set("mysystem.gitanon", "true")
}

// Restore removes anon overrides and re-enables signing.
func Restore() error {
	Unset("user.name")
	Unset("user.email")
	Unset("user.signingKey")
	Unset("mysystem.gitanon")
	if err := Set("commit.gpgSign", "true"); err != nil {
		return err
	}
	if err := Set("tag.gpgSign", "true"); err != nil {
		return err
	}
	return Set("push.gpgSign", "true")
}
