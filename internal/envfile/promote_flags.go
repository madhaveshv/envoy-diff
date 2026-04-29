package envfile

import "github.com/spf13/cobra"

// RegisterPromoteFlags attaches promote-related flags to cmd.
func RegisterPromoteFlags(cmd *cobra.Command) {
	cmd.Flags().StringSlice("only", nil, "promote only these keys")
	cmd.Flags().StringSlice("exclude", nil, "skip these keys during promotion")
	cmd.Flags().Bool("overwrite", false, "overwrite existing keys in destination")
	cmd.Flags().Bool("fail-on-missing", false, "fail if any --only key is absent from source")
}

// PromoteOptionsFromFlags builds a PromoteOptions from cobra flags.
func PromoteOptionsFromFlags(cmd *cobra.Command) (PromoteOptions, error) {
	only, err := cmd.Flags().GetStringSlice("only")
	if err != nil {
		return PromoteOptions{}, err
	}
	exclude, err := cmd.Flags().GetStringSlice("exclude")
	if err != nil {
		return PromoteOptions{}, err
	}
	overwrite, err := cmd.Flags().GetBool("overwrite")
	if err != nil {
		return PromoteOptions{}, err
	}
	failOnMissing, err := cmd.Flags().GetBool("fail-on-missing")
	if err != nil {
		return PromoteOptions{}, err
	}
	return PromoteOptions{
		OnlyKeys:          only,
		ExcludeKeys:       exclude,
		OverwriteExisting: overwrite,
		FailOnMissing:     failOnMissing,
	}, nil
}
