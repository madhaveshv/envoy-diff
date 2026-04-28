package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envoy-diff/internal/envfile"
)

var lintCmd = &cobra.Command{
	Use:   "lint <file>",
	Short: "Lint an env file for style and structural issues",
	Args:  cobra.ExactArgs(1),
	RunE:  runLint,
}

func init() {
	envfile.RegisterLintFlags(lintCmd)
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	env, err := envfile.Load(args[0])
	if err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	opts := envfile.LintOptionsFromFlags(cmd)
	issues := envfile.Lint(env, opts)

	if len(issues) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No lint issues found.")
		return nil
	}

	// Sort for deterministic output.
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Key < issues[j].Key
	})

	for _, issue := range issues {
		fmt.Fprintln(cmd.OutOrStdout(), issue.String())
	}

	fmt.Fprintf(cmd.OutOrStdout(), "\n%d issue(s) found.\n", len(issues))
	os.Exit(1)
	return nil
}
