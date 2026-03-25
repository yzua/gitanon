package cmd

import "fmt"

// requireRepo checks we're inside a git repo.
func requireRepo() error {
	// This is a stub — actual check happens in each command via git.IsInsideRepo().
	// We keep this for future global validation.
	return nil
}

// repoErr wraps a repo check error.
func repoErr() error {
	return fmt.Errorf("not inside a git repository")
}
