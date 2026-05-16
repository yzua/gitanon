package model_test

import (
	"strings"
	"testing"

	"github.com/yzua/gitanon/internal/model"
)

func TestGitUserString(t *testing.T) {
	t.Run("full identity", func(t *testing.T) {
		u := model.GitUser{
			Name:     "Alice",
			Email:    "alice@example.com",
			Signing:  true,
			SignKey:  "KEY123",
			AnonMode: false,
			IsLocal:  true,
			RepoName: "myproject",
		}
		s := u.String()
		for _, want := range []string{
			"Repo:     myproject",
			"Name:     Alice",
			"Email:    alice@example.com",
			"Signing:  true",
			"Key:      KEY123",
			"AnonMode: off",
			"Scope:    local (override)",
		} {
			if !strings.Contains(s, want) {
				t.Errorf("String() missing %q\n got:\n%s", want, s)
			}
		}
	})

	t.Run("empty fields show defaults", func(t *testing.T) {
		u := model.GitUser{
			RepoName: "test",
		}
		s := u.String()
		for _, want := range []string{
			"Name:     (not set)",
			"Email:    (none)",
			"Signing:  false",
			"Key:      (none)",
			"AnonMode: off",
			"Scope:    global (inherited)",
		} {
			if !strings.Contains(s, want) {
				t.Errorf("String() missing %q\n got:\n%s", want, s)
			}
		}
	})

	t.Run("anon mode", func(t *testing.T) {
		u := model.GitUser{
			Name:     "user",
			AnonMode: true,
			RepoName: "secret",
		}
		s := u.String()
		if !strings.Contains(s, "AnonMode: on") {
			t.Errorf("String() missing 'AnonMode: on'\n got:\n%s", s)
		}
	})
}
