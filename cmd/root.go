// Package cmd contains all gitanon CLI subcommands.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var Version = "dev"

// rootCmd is the base command when called without subcommands.
var rootCmd = &cobra.Command{
	Use:   "gitanon",
	Short: "Anonymous git identity manager",
	Long: `gitanon lets you commit anonymously, impersonate other GitHub users,
and manage signing behavior per-repo — without touching your global git config.

  gitanon on          Anonymize current repo
  gitanon off         Restore global identity and re-enable signing
  gitanon as <user>   Commit as another GitHub user (fetches identity from API)
  gitanon whoami      Show current repo identity
  gitanon hook        Generate hook scripts for signing enforcement`,
	SilenceUsage: true,
	Version:      Version,
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
