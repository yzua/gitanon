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

		name := user.Name
		if name == "" {
			name = "(not set)"
		}
		email := user.Email
		if email == "" {
			email = "(none)"
		}
		key := user.SignKey
		if key == "" {
			key = "(none)"
		}

		signStr := "false"
		if user.Signing {
			signStr = "true"
		}

		scope := "global (inherited)"
		if user.IsLocal {
			scope = "local (override)"
		}

		anonMode := "off"
		if user.AnonMode {
			anonMode = "on"
		}

		fmt.Printf("Repo:     %s\n", git.RepoName())
		fmt.Printf("Name:     %s\n", name)
		fmt.Printf("Email:    %s\n", email)
		fmt.Printf("Signing:  %s\n", signStr)
		fmt.Printf("Key:      %s\n", key)
		fmt.Printf("AnonMode: %s\n", anonMode)
		fmt.Printf("Scope:    %s\n", scope)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
