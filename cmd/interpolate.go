package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-diff/internal/envfile"
)

var interpolateCmd = &cobra.Command{
	Use:   "interpolate <file>",
	Short: "Resolve variable references in an env file and print the result",
	Args:  cobra.ExactArgs(1),
	RunE:  runInterpolate,
}

func init() {
	envfile.RegisterInterpolateFlags(interpolateCmd)
	interpolateCmd.Flags().StringP("format", "f", "dotenv", "Output format: dotenv, shell, json")
	RootCmd.AddCommand(interpolateCmd)
}

func runInterpolate(cmd *cobra.Command, args []string) error {
	env, err := envfile.Load(args[0])
	if err != nil {
		return fmt.Errorf("loading file: %w", err)
	}

	enabled, opts := envfile.InterpolateFlagsFromArgs(cmd)
	if enabled {
		env, err = envfile.Interpolate(env, opts)
		if err != nil {
			return fmt.Errorf("interpolation failed: %w", err)
		}
	}

	format, _ := cmd.Flags().GetString("format")
	out, err := envfile.Render(env, format)
	if err != nil {
		return err
	}
	fmt.Fprint(os.Stdout, out)
	return nil
}
