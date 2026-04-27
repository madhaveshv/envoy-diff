package envfile

import "github.com/spf13/cobra"

// InterpolateOptions is defined in interpolate.go.

// RegisterInterpolateFlags registers interpolation-related flags on a cobra command.
func RegisterInterpolateFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("interpolate", false, "Resolve variable references in values (e.g. ${VAR})")
	cmd.Flags().Bool("interpolate-os", false, "Fall back to OS environment when resolving variable references")
	cmd.Flags().Bool("interpolate-strict", false, "Fail if a referenced variable cannot be resolved")
}

// InterpolateFlagsFromArgs reads interpolation flags from a cobra command and
// returns whether interpolation is enabled and the corresponding options.
func InterpolateFlagsFromArgs(cmd *cobra.Command) (bool, InterpolateOptions) {
	enabled, _ := cmd.Flags().GetBool("interpolate")
	useOS, _ := cmd.Flags().GetBool("interpolate-os")
	strict, _ := cmd.Flags().GetBool("interpolate-strict")
	return enabled, InterpolateOptions{
		UseOSEnv:      useOS,
		FailOnMissing: strict,
	}
}
