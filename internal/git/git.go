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

const (
	localScope  = "--local"
	globalScope = "--global"

	keyUserName        = "user.name"
	keyUserEmail       = "user.email"
	keySigningKey      = "user.signingKey"
	keyCommitSign      = "commit.gpgSign"
	keyTagSign         = "tag.gpgSign"
	keyPushSign        = "push.gpgSign"
	keyAnonMode        = "mysystem.gitanon"
	keyPriorCommitSign = "mysystem.gitanon.priorcommitsign"
	keyPriorTagSign    = "mysystem.gitanon.priortagsign"
	keyPriorPushSign   = "mysystem.gitanon.priorpushsign"
	configTrue         = "true"
	configFalse        = "false"
	anonDefaultName    = "user"
	unknownRepo        = "(unknown)"
)

type configEntry struct {
	key   string
	value string
}

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
	return Run("config", localScope, key, value)
}

// Unset removes a local git config key (no error if missing).
func Unset(key string) {
	_ = Run("config", localScope, "--unset", key)
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
		return unknownRepo
	}
	return filepath.Base(strings.TrimSpace(string(out)))
}

// WhoAmI reads the current git identity.
func WhoAmI() model.GitUser {
	user := readUser(localScope)
	user.IsLocal = user.Name != ""
	user.AnonMode = Get(localScope, keyAnonMode) == configTrue

	// Fall back to global if no local override
	if !user.IsLocal {
		globalUser := readUser(globalScope)
		user.Name = globalUser.Name
		user.Email = globalUser.Email
		user.SignKey = globalUser.SignKey
		user.Signing = globalUser.Signing
	}
	return user
}

// Anonymize sets the repo to anonymous mode with a custom name and email.
// It saves the current signing state so Restore can recover it.
// Skips saving prior state if already in anon mode (preserves the original).
func Anonymize(name, email string) error {
	if Get(localScope, keyAnonMode) != configTrue {
		if err := savePriorSigning(); err != nil {
			return err
		}
	}
	return setLocal(
		configEntry{key: keyUserName, value: name},
		configEntry{key: keyUserEmail, value: email},
		configEntry{key: keyCommitSign, value: configFalse},
		configEntry{key: keyTagSign, value: configFalse},
		configEntry{key: keyPushSign, value: configFalse},
		configEntry{key: keySigningKey, value: ""},
		configEntry{key: keyAnonMode, value: configTrue},
	)
}

// AnonymizeDefault sets the repo to anonymous mode with the default identity.
func AnonymizeDefault() error {
	return Anonymize(anonDefaultName, "")
}

// Restore removes anon overrides and restores prior signing state.
func Restore() error {
	commitSign := priorOrDefault(keyPriorCommitSign, configTrue)
	tagSign := priorOrDefault(keyPriorTagSign, configTrue)
	pushSign := priorOrDefault(keyPriorPushSign, configTrue)

	unsetLocal(keyUserName, keyUserEmail, keySigningKey, keyAnonMode,
		keyPriorCommitSign, keyPriorTagSign, keyPriorPushSign)

	return setLocal(
		configEntry{key: keyCommitSign, value: commitSign},
		configEntry{key: keyTagSign, value: tagSign},
		configEntry{key: keyPushSign, value: pushSign},
	)
}

// VerifyCommitSigning checks that the HEAD commit has a valid GPG signature.
func VerifyCommitSigning() error {
	return Run("verify-commit", "HEAD")
}

// SigningEnabled reports whether commit signing is active
// (local or global config).
func SigningEnabled() bool {
	val := Get(localScope, keyCommitSign)
	if val == "" {
		val = Get(globalScope, keyCommitSign)
	}
	return val == configTrue
}

func readUser(scope string) model.GitUser {
	return model.GitUser{
		Name:    Get(scope, keyUserName),
		Email:   Get(scope, keyUserEmail),
		SignKey: Get(scope, keySigningKey),
		Signing: Get(scope, keyCommitSign) == configTrue,
	}
}

func setLocal(entries ...configEntry) error {
	for _, entry := range entries {
		if err := Set(entry.key, entry.value); err != nil {
			return err
		}
	}
	return nil
}

func unsetLocal(keys ...string) {
	for _, key := range keys {
		Unset(key)
	}
}

func savePriorSigning() error {
	prior := []configEntry{
		{key: keyPriorCommitSign, value: signingValue(keyCommitSign)},
		{key: keyPriorTagSign, value: signingValue(keyTagSign)},
		{key: keyPriorPushSign, value: signingValue(keyPushSign)},
	}
	return setLocal(prior...)
}

// signingValue returns "true" or "false" for a signing key,
// checking local then global scope.
func signingValue(key string) string {
	val := Get(localScope, key)
	if val == "" {
		val = Get(globalScope, key)
	}
	if val == configTrue {
		return configTrue
	}
	return configFalse
}

// priorOrDefault reads a saved prior value; falls back to default if unset.
func priorOrDefault(priorKey, fallback string) string {
	val := Get(localScope, priorKey)
	if val == "" {
		return fallback
	}
	return val
}
