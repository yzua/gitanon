// Package model defines shared types for gitanon.
package model

// GitUser holds identity fields from git config.
type GitUser struct {
	Name     string
	Email    string
	Signing  bool
	SignKey  string
	AnonMode bool
	IsLocal  bool
}

// GitHubUser holds public fields from the GitHub API.
type GitHubUser struct {
	Login string `json:"login"`
	ID    int64  `json:"id"`
	Name  string `json:"name"`
}
