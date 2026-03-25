// Package github provides GitHub API helpers.
package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yzua/gitanon/internal/model"
)

const apiBase = "https://api.github.com"

// LookupUser fetches public info for a GitHub username.
// No authentication required — this is a public endpoint.
// Rate limit: 60 requests/hour (unauthenticated).
func LookupUser(username string) (*model.GitHubUser, error) {
	url := fmt.Sprintf("%s/users/%s", apiBase, username)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "gitanon-cli")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		// ok
	case http.StatusNotFound:
		return nil, fmt.Errorf("GitHub user %q not found", username)
	case http.StatusForbidden:
		return nil, fmt.Errorf("GitHub API rate limit exceeded (60 req/hour unauthenticated)")
	default:
		return nil, fmt.Errorf("GitHub API error %d: %s", resp.StatusCode, string(body))
	}

	var user model.GitHubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &user, nil
}

// NoreplyEmail returns the GitHub noreply email for a user.
// Format: <id>+<login>@users.noreply.github.com
func NoreplyEmail(u *model.GitHubUser) string {
	return fmt.Sprintf("%d+%s@users.noreply.github.com", u.ID, u.Login)
}

// DisplayName returns the best display name for a user.
func DisplayName(u *model.GitHubUser) string {
	if u.Name != "" {
		return u.Name
	}
	return u.Login
}
