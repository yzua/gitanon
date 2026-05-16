package github_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yzua/gitanon/internal/github"
	"github.com/yzua/gitanon/internal/model"
)

func TestLookupUser(t *testing.T) {
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

func TestLookupUserHTTP(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		payload, _ := json.Marshal(model.GitHubUser{Login: "octocat", ID: 583231, Name: "The Octocat"})
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/users/octocat" {
				t.Errorf("request path = %q, want /users/octocat", r.URL.Path)
			}
			if r.Header.Get("Accept") != "application/vnd.github+json" {
				t.Errorf("Accept header = %q, want application/vnd.github+json", r.Header.Get("Accept"))
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(payload)
		}))
		defer srv.Close()

		user, err := github.LookupUser("octocat", github.HTTPConfig{
			Client:  srv.Client(),
			BaseURL: srv.URL,
		})
		if err != nil {
			t.Fatalf("LookupUser() error: %v", err)
		}
		if user.Login != "octocat" {
			t.Errorf("Login = %q, want %q", user.Login, "octocat")
		}
		if user.ID != 583231 {
			t.Errorf("ID = %d, want %d", user.ID, 583231)
		}
	})

	t.Run("not found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"Not Found"}`))
		}))
		defer srv.Close()

		_, err := github.LookupUser("nonexistent", github.HTTPConfig{
			Client:  srv.Client(),
			BaseURL: srv.URL,
		})
		if err == nil {
			t.Fatal("expected error for 404")
		}
		want := `GitHub user "nonexistent" not found`
		if err.Error() != want {
			t.Errorf("error = %q, want %q", err.Error(), want)
		}
	})

	t.Run("rate limited", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"message":"API rate limit exceeded"}`))
		}))
		defer srv.Close()

		_, err := github.LookupUser("octocat", github.HTTPConfig{
			Client:  srv.Client(),
			BaseURL: srv.URL,
		})
		if err == nil {
			t.Fatal("expected error for 403")
		}
		if err.Error() != "GitHub API rate limit exceeded (60 req/hour unauthenticated)" {
			t.Errorf("error = %q, want rate limit message", err.Error())
		}
	})

	t.Run("malformed response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`not json`))
		}))
		defer srv.Close()

		_, err := github.LookupUser("octocat", github.HTTPConfig{
			Client:  srv.Client(),
			BaseURL: srv.URL,
		})
		if err == nil {
			t.Fatal("expected error for malformed JSON")
		}
	})

	t.Run("network error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		}))
		defer srv.Close()

		client := &http.Client{}
		client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("connection refused")
		})

		_, err := github.LookupUser("octocat", github.HTTPConfig{
			Client:  client,
			BaseURL: srv.URL,
		})
		if err == nil {
			t.Fatal("expected error for network failure")
		}
	})
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
