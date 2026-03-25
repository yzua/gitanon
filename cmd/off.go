package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yzua/gitanon/internal/git"
)

var offCmd = &cobra.Command{
	Use:     "off",
	Short:   "Restore global identity and re-enable signing",
	Aliases: []string{"back", "undo"},
	Long: `Remove anonymous overrides from local git config and re-enable GPG signing.
After this, the repo falls back to your global git identity.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !git.IsInsideRepo() {
			return fmt.Errorf("not inside a git repository")
		}

		if err := git.Restore(); err != nil {
			return fmt.Errorf("restoring identity: %w", err)
		}

		fmt.Printf("✔ Restored global identity in %s\n", git.RepoName())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(offCmd)
}
