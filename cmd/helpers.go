package cmd

import (
	"fmt"

	"github.com/yzua/gitanon/internal/git"
)

// requireRepo checks we're inside a git repo.
func requireRepo() error {
	if git.IsInsideRepo() {
		return nil
	}
	return fmt.Errorf("not inside a git repository")
}
