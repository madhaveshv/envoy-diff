package envfile

import "github.com/spf13/cobra"

// DiffFlags holds CLI flag values for the diff command.
type DiffFlags struct {
	Format    string
	Prefix    string
	Allowlist []string
	Exclude   []string
	Redact    bool
	Audit     bool
}

// RegisterDiffFlags registers diff-related flags on the given command.
func RegisterDiffFlags(cmd *cobra.Command) {
	cmd.Flags().String("format", "text", "Output format: text, table, json")
	cmd.Flags().String("prefix", "", "Only include keys with this prefix")
	cmd.Flags().StringSlice("allowlist", nil, "Only include these keys")
	cmd.Flags().StringSlice("exclude", nil, "Exclude these keys from diff")
	cmd.Flags().Bool("redact", false, "Redact sensitive values in output")
	cmd.Flags().Bool("audit", false, "Run audit rules on the diff")
}

// DiffFlagsFromArgs extracts DiffFlags from a cobra command's parsed flags.
func DiffFlagsFromArgs(cmd *cobra.Command) DiffFlags {
	format, _ := cmd.Flags().GetString("format")
	prefix, _ := cmd.Flags().GetString("prefix")
	allowlist, _ := cmd.Flags().GetStringSlice("allowlist")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	redact, _ := cmd.Flags().GetBool("redact")
	audit, _ := cmd.Flags().GetBool("audit")

	return DiffFlags{
		Format:    format,
		Prefix:    prefix,
		Allowlist: allowlist,
		Exclude:   exclude,
		Redact:    redact,
		Audit:     audit,
	}
}
