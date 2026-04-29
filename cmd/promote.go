package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-diff/internal/envfile"
)

var promoteCmd = &cobra.Command{
	Use:   "promote <src> <dst>",
	Short: "Promote env vars from a source file into a destination file",
	Args:  cobra.ExactArgs(2),
	RunE:  runPromote,
}

func init() {
	envfile.RegisterPromoteFlags(promoteCmd)
	promoteCmd.Flags().StringP("format", "f", "dotenv", "output format: dotenv, shell, json")
	promoteCmd.Flags().StringP("output", "o", "", "write result to file instead of stdout")
	rootCmd.AddCommand(promoteCmd)
}

func runPromote(cmd *cobra.Command, args []string) error {
	srcPath, dstPath := args[0], args[1]

	src, err := envfile.Load(srcPath)
	if err != nil {
		return fmt.Errorf("loading source: %w", err)
	}
	dst, err := envfile.Load(dstPath)
	if err != nil {
		return fmt.Errorf("loading destination: %w", err)
	}

	opts, err := envfile.PromoteOptionsFromFlags(cmd)
	if err != nil {
		return err
	}

	result, err := envfile.Promote(src, dst, opts)
	if err != nil {
		return err
	}

	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")

	rendered, err := envfile.Render(result, format)
	if err != nil {
		return err
	}

	if output != "" {
		if err := os.WriteFile(output, []byte(rendered), 0644); err != nil {
			return fmt.Errorf("writing output file: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "wrote promoted env to %s\n", output)
		return nil
	}

	fmt.Fprint(cmd.OutOrStdout(), rendered)
	return nil
}
