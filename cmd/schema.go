package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-diff/internal/envfile"
)

var schemaFile string

var schemaCmd = &cobra.Command{
	Use:   "schema <env-file>",
	Short: "Validate an env file against a JSON schema definition",
	Args:  cobra.ExactArgs(1),
	RunE:  runSchema,
}

func init() {
	schemaCmd.Flags().StringVarP(&schemaFile, "schema", "s", "", "Path to JSON schema file (required)")
	_ = schemaCmd.MarkFlagRequired("schema")
	rootCmd.AddCommand(schemaCmd)
}

func runSchema(cmd *cobra.Command, args []string) error {
	env, err := envfile.Load(args[0])
	if err != nil {
		return fmt.Errorf("loading env file: %w", err)
	}

	schemaBytes, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("reading schema file: %w", err)
	}

	var fields []envfile.SchemaField
	if err := json.Unmarshal(schemaBytes, &fields); err != nil {
		return fmt.Errorf("parsing schema file: %w", err)
	}

	results := envfile.ValidateSchema(env, fields)
	if len(results) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✔ schema validation passed")
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✖ schema validation failed (%d issue(s)):\n", len(results))
	for _, r := range results {
		fmt.Fprintf(cmd.OutOrStdout(), "  [%s] %s\n", r.Key, r.Message)
	}
	os.Exit(1)
	return nil
}
