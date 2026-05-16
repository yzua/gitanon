// Package model defines shared types for gitanon.
package model

import "fmt"

// GitUser holds identity fields from git config.
type GitUser struct {
	Name     string
	Email    string
	Signing  bool
	SignKey  string
	AnonMode bool
	IsLocal  bool
	RepoName string
}

func (u GitUser) String() string {
	name := u.Name
	if name == "" {
		name = "(not set)"
	}
	email := u.Email
	if email == "" {
		email = "(none)"
	}
	key := u.SignKey
	if key == "" {
		key = "(none)"
	}

	signStr := "false"
	if u.Signing {
		signStr = "true"
	}

	scope := "global (inherited)"
	if u.IsLocal {
		scope = "local (override)"
	}

	anonMode := "off"
	if u.AnonMode {
		anonMode = "on"
	}

	return fmt.Sprintf("Repo:     %s\nName:     %s\nEmail:    %s\nSigning:  %s\nKey:      %s\nAnonMode: %s\nScope:    %s\n",
		u.RepoName, name, email, signStr, key, anonMode, scope)
}

// GitHubUser holds public fields from the GitHub API.
type GitHubUser struct {
	Login string `json:"login"`
	ID    int64  `json:"id"`
	Name  string `json:"name"`
}
