package github_test

import (
	"testing"

	"github.com/yzua/gitanon/internal/github"
	"github.com/yzua/gitanon/internal/model"
)

func TestLookupUser(t *testing.T) {
	// Test helpers with mock data (no network call needed)
	t.Run("DisplayName with name", func(t *testing.T) {
		user := &model.GitHubUser{Login: "octocat", ID: 583231, Name: "The Octocat"}
		got := github.DisplayName(user)
		if got != "The Octocat" {
			t.Errorf("DisplayName() = %q, want %q", got, "The Octocat")
		}
	})

	t.Run("DisplayName without name", func(t *testing.T) {
		user := &model.GitHubUser{Login: "octocat", ID: 583231, Name: ""}
		got := github.DisplayName(user)
		if got != "octocat" {
			t.Errorf("DisplayName() = %q, want %q", got, "octocat")
		}
	})

	t.Run("NoreplyEmail", func(t *testing.T) {
		user := &model.GitHubUser{Login: "octocat", ID: 583231}
		got := github.NoreplyEmail(user)
		want := "583231+octocat@users.noreply.github.com"
		if got != want {
			t.Errorf("NoreplyEmail() = %q, want %q", got, want)
		}
	})
}
