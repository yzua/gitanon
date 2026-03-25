package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yzua/gitanon/internal/git"
)

var onCmd = &cobra.Command{
	Use:     "on",
	Short:   "Anonymize the current repo",
	Aliases: []string{"anon"},
	Long: `Set local git config to anonymous mode: no name, no email, no signing.
Sets the mysystem.gitanon=true flag so hooks can detect anonymous commits.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRepo(); err != nil {
			return err
		}

		if err := git.Anonymize("user", ""); err != nil {
			return fmt.Errorf("setting anonymous identity: %w", err)
		}

		fmt.Printf("✔ Anonymous mode in %s\n", git.RepoName())
		fmt.Println("  Name:    user")
		fmt.Println("  Email:   (none)")
		fmt.Println("  Signing: disabled")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(onCmd)
}
