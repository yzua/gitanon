package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/yzua/gitanon/internal/git"
)

var hookCmd = &cobra.Command{
	Use:   "hook <name>",
	Short: "Run a gitanon-aware git hook",
	Long: `Runs a git hook script that respects the mysystem.gitanon flag.

Usage in your hooks:
  #!/bin/sh
  gitanon hook pre-push

Available hooks:
  pre-commit  — placeholder (extend with your own checks)
  pre-push    — verifies GPG signatures unless anon mode is on`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRepo(); err != nil {
			return err
		}

		hookName := args[0]

		// Check if anon mode is active — if so, skip signing enforcement
		if git.WhoAmI().AnonMode {
			fmt.Fprintf(os.Stderr, "⚠ gitanon: anonymous mode, skipping %s signing checks\n", hookName)
			return nil
		}

		switch hookName {
		case "pre-commit":
			return runPreCommit()
		case "pre-push":
			return runPrePush()
		default:
			return fmt.Errorf("unknown hook %q (supported: pre-commit, pre-push)", hookName)
		}
	},
}

func runPreCommit() error {
	// Check for GPG signing
	signStr := git.Get("--local", "commit.gpgSign")
	if signStr == "" {
		signStr = git.Get("--global", "commit.gpgSign")
	}
	if signStr != "true" {
		fmt.Fprintln(os.Stderr, "⚠ gitanon: commit.gpgSign is not true")
	}
	return nil
}

func runPrePush() error {
	if err := exec.Command("git", "verify-commit", "HEAD").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "✗ gitanon: latest commit is NOT GPG-signed!")
		fmt.Fprintln(os.Stderr, "  Fix: git commit --amend -S --no-edit")
		return fmt.Errorf("unsigned commit detected")
	}

	fmt.Fprintln(os.Stderr, "✔ gitanon: commit has valid GPG signature")
	return nil
}

func init() {
	rootCmd.AddCommand(hookCmd)
}
