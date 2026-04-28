package envfile

import "github.com/spf13/cobra"

// RegisterLintFlags attaches lint-related flags to the given command.
func RegisterLintFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("lint-no-leading-underscore", false, "disallow keys that start with an underscore")
	cmd.Flags().Bool("lint-require-uppercase", false, "disallow keys with lowercase characters")
	cmd.Flags().Bool("lint-no-numeric-start", false, "disallow keys that start with a numeric character")
	cmd.Flags().Int("lint-max-key-length", 0, "maximum allowed key length (0 = unlimited)")
	cmd.Flags().Int("lint-max-value-length", 0, "maximum allowed value length (0 = unlimited)")
}

// LintOptionsFromFlags reads lint flags from a cobra command and returns
// a populated LintOptions struct.
func LintOptionsFromFlags(cmd *cobra.Command) LintOptions {
	noUnderscore, _ := cmd.Flags().GetBool("lint-no-leading-underscore")
	requireUpper, _ := cmd.Flags().GetBool("lint-require-uppercase")
	noNumeric, _ := cmd.Flags().GetBool("lint-no-numeric-start")
	maxKey, _ := cmd.Flags().GetInt("lint-max-key-length")
	maxVal, _ := cmd.Flags().GetInt("lint-max-value-length")

	return LintOptions{
		DisallowLeadingUnderscore: noUnderscore,
		DisallowLowercase:         requireUpper,
		DisallowNumericStart:      noNumeric,
		MaxKeyLength:              maxKey,
		MaxValueLength:            maxVal,
	}
}
