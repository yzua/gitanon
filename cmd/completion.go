package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion <shell>",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for gitanon.

Supported shells: bash, zsh, fish, powershell

Installation:
  Bash:  gitanon completion bash > /etc/bash_completion.d/gitanon
  Zsh:   gitanon completion zsh > "${fpath[1]}/_gitanon"
  Fish:  gitanon completion fish > ~/.config/fish/completions/gitanon.fish`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			err = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return fmt.Errorf("unsupported shell %q", args[0])
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
