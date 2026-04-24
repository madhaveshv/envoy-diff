package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-diff/internal/diff"
	"github.com/user/envoy-diff/internal/envfile"
)

var (
	formatFlag string
	outputFlag string
)

var rootCmd = &cobra.Command{
	Use:   "envoy-diff <from-file> <to-file>",
	Short: "Diff and audit environment variable changes across deployment stages",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		fromEnv, err := envfile.ParseFile(args[0])
		if err != nil {
			return fmt.Errorf("reading from-file: %w", err)
		}

		toEnv, err := envfile.ParseFile(args[1])
		if err != nil {
			return fmt.Errorf("reading to-file: %w", err)
		}

		results := diff.Compare(fromEnv, toEnv)

		var output string
		switch formatFlag {
		case "json":
			output, err = diff.FormatJSON(results)
			if err != nil {
				return fmt.Errorf("formatting JSON: %w", err)
			}
		case "table":
			output = diff.FormatTable(results)
		default:
			output = diff.FormatText(results)
		}

		if outputFlag != "" {
			if err := os.WriteFile(outputFlag, []byte(output), 0644); err != nil {
				return fmt.Errorf("writing output file: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Output written to %s\n", outputFlag)
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), output)
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&formatFlag, "format", "f", "text", "Output format: text, table, json")
	rootCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Write output to file instead of stdout")
}
