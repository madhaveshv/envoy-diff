package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-diff/internal/audit"
	"envoy-diff/internal/diff"
	"envoy-diff/internal/envfile"
	"envoy-diff/internal/report"
)

var diffCmd = &cobra.Command{
	Use:   "diff <before> <after>",
	Short: "Diff two environment files",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	envfile.RegisterDiffFlags(diffCmd)
	RootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	flags := envfile.DiffFlagsFromArgs(cmd)

	before, after, err := envfile.LoadPair(args[0], args[1])
	if err != nil {
		return fmt.Errorf("loading env files: %w", err)
	}

	filterOpts := envfile.FilterFlagsFromArgs(cmd)
	before = envfile.Filter(before, filterOpts)
	after = envfile.Filter(after, filterOpts)

	if flags.Redact {
		before = envfile.Redact(before, envfile.DefaultSensitivePatterns)
		after = envfile.Redact(after, envfile.DefaultSensitivePatterns)
	}

	changes := diff.Compare(before, after)

	var issues []audit.Issue
	if flags.Audit {
		auditor := audit.New(audit.DefaultRules())
		issues = auditor.Run(changes)
	}

	r := report.New(changes, issues)

	var output string
	switch flags.Format {
	case "json":
		output, err = report.FormatJSON(r)
		if err != nil {
			return fmt.Errorf("formatting output: %w", err)
		}
	case "table":
		output = report.FormatTable(r)
	default:
		output = report.FormatText(r)
	}

	fmt.Fprintln(os.Stdout, output)
	return nil
}
