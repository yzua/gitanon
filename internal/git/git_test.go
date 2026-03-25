package git_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/yzua/gitanon/internal/git"
)

// setupTempRepo creates a temporary git repo for testing.
func setupTempRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.name", "Test User"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "commit.gpgSign", "true"},
		{"git", "config", "user.signingKey", "ABC123"},
	}

	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%s failed: %s: %v", args, string(out), err)
		}
	}
	return dir
}

func TestAnonymizeAndRestore(t *testing.T) {
	dir := setupTempRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Verify initial state
	user := git.WhoAmI()
	if user.Name != "Test User" {
		t.Errorf("initial name = %q, want %q", user.Name, "Test User")
	}
	if !user.Signing {
		t.Error("initial signing should be true")
	}
	if user.AnonMode {
		t.Error("initial anon mode should be false")
	}

	// Anonymize
	if err := git.Anonymize("anon-user", "anon@example.com"); err != nil {
		t.Fatalf("Anonymize() error: %v", err)
	}

	user = git.WhoAmI()
	if user.Name != "anon-user" {
		t.Errorf("anon name = %q, want %q", user.Name, "anon-user")
	}
	if user.Email != "anon@example.com" {
		t.Errorf("anon email = %q, want %q", user.Email, "anon@example.com")
	}
	if user.Signing {
		t.Error("anon signing should be false")
	}
	if !user.AnonMode {
		t.Error("anon mode should be true after Anonymize()")
	}

	// Restore
	if err := git.Restore(); err != nil {
		t.Fatalf("Restore() error: %v", err)
	}

	// Verify anon marker is removed
	anonFlag := git.Get("--local", "mysystem.gitanon")
	if anonFlag != "" {
		t.Errorf("after restore, mysystem.gitanon should be unset, got %q", anonFlag)
	}

	// Verify signing is re-enabled
	signVal := git.Get("--local", "commit.gpgSign")
	if signVal != "true" {
		t.Errorf("after restore, commit.gpgSign = %q, want true", signVal)
	}

	// Local name should be unset (falls back to global if available)
	localName := git.Get("--local", "user.name")
	if localName != "" {
		t.Errorf("after restore, local user.name should be unset, got %q", localName)
	}
}

func TestIsInsideRepo(t *testing.T) {
	// Should be true when we're in a repo
	dir := setupTempRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	if !git.IsInsideRepo() {
		t.Error("IsInsideRepo() = false in a git repo")
	}

	// Should be false in a temp dir without git
	noRepo := t.TempDir()
	os.Chdir(noRepo)
	if git.IsInsideRepo() {
		t.Error("IsInsideRepo() = true outside a git repo")
	}
}

func TestRepoName(t *testing.T) {
	dir := setupTempRepo(t)
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	name := git.RepoName()
	if name != filepath.Base(dir) {
		t.Errorf("RepoName() = %q, want %q", name, filepath.Base(dir))
	}
}
