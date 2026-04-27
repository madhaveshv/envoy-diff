package envfile

import "github.com/spf13/cobra"

// ValidateFlags holds CLI flag values for the validate command.
type ValidateFlags struct {
	RequireUppercase bool
	DisallowEmpty    bool
	ForbiddenKeys    []string
	KeyPattern       string
}

// RegisterValidateFlags binds validate-related flags onto a cobra command.
func RegisterValidateFlags(cmd *cobra.Command, f *ValidateFlags) {
	cmd.Flags().BoolVar(&f.RequireUppercase, "require-uppercase", false, "warn if any key is not uppercase")
	cmd.Flags().BoolVar(&f.DisallowEmpty, "disallow-empty", false, "warn if any value is empty")
	cmd.Flags().StringSliceVar(&f.ForbiddenKeys, "forbidden-keys", nil, "comma-separated list of keys that must not be present")
	cmd.Flags().StringVar(&f.KeyPattern, "key-pattern", "", "regex pattern that all keys must match")
}

// ValidateOptionsFromFlags converts CLI flags into a ValidateOptions struct.
func ValidateOptionsFromFlags(f ValidateFlags) ValidateOptions {
	return ValidateOptions{
		RequireUppercase: f.RequireUppercase,
		DisallowEmpty:    f.DisallowEmpty,
		ForbiddenKeys:    f.ForbiddenKeys,
		KeyPattern:       f.KeyPattern,
	}
}
