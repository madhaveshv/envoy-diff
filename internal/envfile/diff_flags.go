package envfile

import "github.com/spf13/cobra"

// DiffFlags holds options that control how env file diffs are filtered and displayed.
type DiffFlags struct {
	Prefix      string
	Allowlist   []string
	ExcludeKeys []string
	Redact      bool
	Audit       bool
	Format      string
}

// RegisterDiffFlags attaches diff-related flags to a cobra command.
func RegisterDiffFlags(cmd *cobra.Command) {
	cmd.Flags().String("prefix", "", "Only include keys with this prefix")
	cmd.Flags().StringSlice("allowlist", nil, "Comma-separated list of keys to include")
	cmd.Flags().StringSlice("exclude", nil, "Comma-separated list of keys to exclude")
	cmd.Flags().Bool("redact", false, "Redact sensitive values in output")
	cmd.Flags().Bool("audit", false, "Run audit rules on the diff")
	cmd.Flags().String("format", "text", "Output format: text, table, json")
}

// DiffFlagsFromArgs extracts DiffFlags from a cobra command's parsed flags.
func DiffFlagsFromArgs(cmd *cobra.Command) DiffFlags {
	prefix, _ := cmd.Flags().GetString("prefix")
	allowlist, _ := cmd.Flags().GetStringSlice("allowlist")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	redact, _ := cmd.Flags().GetBool("redact")
	audit, _ := cmd.Flags().GetBool("audit")
	format, _ := cmd.Flags().GetString("format")

	return DiffFlags{
		Prefix:      prefix,
		Allowlist:   allowlist,
		ExcludeKeys: exclude,
		Redact:      redact,
		Audit:       audit,
		Format:      format,
	}
}
