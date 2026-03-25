package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yzua/gitanon/internal/git"
	gh "github.com/yzua/gitanon/internal/github"
)

var asCmd = &cobra.Command{
	Use:   "as <username>",
	Short: "Commit as another GitHub user",
	Long: `Fetch a GitHub user's public profile and set the local git config
to commit as them. Uses their display name and GitHub noreply email
(<id>+<username>@users.noreply.github.com).

No GPG signing. Sets mysystem.gitanon=true so hooks skip signing enforcement.

Examples:
  gitanon as octocat
  gitanon as torvalds`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRepo(); err != nil {
			return err
		}

		username := args[0]

		user, err := gh.LookupUser(username)
		if err != nil {
			return err
		}

		name := gh.DisplayName(user)
		email := gh.NoreplyEmail(user)

		if err := git.Anonymize(name, email); err != nil {
			return fmt.Errorf("setting identity: %w", err)
		}

		fmt.Printf("✔ Committing as %s <%s> in %s\n", name, email, git.RepoName())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(asCmd)
}
