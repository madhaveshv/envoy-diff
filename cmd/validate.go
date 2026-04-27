package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-diff/internal/envfile"
)

var validateFlags envfile.ValidateFlags

var validateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate an env file against configurable rules",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := envfile.Load(args[0])
		if err != nil {
			return fmt.Errorf("loading env file: %w", err)
		}

		opts := envfile.ValidateOptionsFromFlags(validateFlags)
		issues := envfile.Validate(env, opts)

		if len(issues) == 0 {
			fmt.Println("✔  No validation issues found.")
			return nil
		}

		fmt.Fprintf(os.Stderr, "Found %d validation issue(s):\n\n", len(issues))
		hasError := false
		for _, issue := range issues {
			fmt.Fprintln(os.Stderr, " ", issue.String())
			if issue.Severity == "error" {
				hasError = true
			}
		}
		fmt.Fprintln(os.Stderr)

		if hasError {
			os.Exit(2)
		}
		return nil
	},
}

func init() {
	envfile.RegisterValidateFlags(validateCmd, &validateFlags)
	rootCmd.AddCommand(validateCmd)
}
