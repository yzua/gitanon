package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yzua/gitanon/internal/git"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current repo identity",
	Long:  `Display the git identity (name, email, signing) for the current repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRepo(); err != nil {
			return err
		}

		user := git.WhoAmI()
		user.RepoName = git.RepoName()
		fmt.Print(user)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
