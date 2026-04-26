package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-diff/internal/envfile"
	"github.com/yourorg/envoy-diff/internal/snapshot"
)

var snapshotDir string

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage environment variable snapshots",
}

var snapshotSaveCmd = &cobra.Command{
	Use:   "save [label] [file]",
	Short: "Save a snapshot of an env file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		label, filePath := args[0], args[1]
		env, err := envfile.Load(filePath)
		if err != nil {
			return fmt.Errorf("load env file: %w", err)
		}
		s := snapshot.New(label, env)
		path, err := snapshot.Save(snapshotDir, s)
		if err != nil {
			return fmt.Errorf("save snapshot: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Snapshot saved: %s\n", path)
		return nil
	},
}

var snapshotListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved snapshots",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := snapshot.List(snapshotDir)
		if err != nil {
			return fmt.Errorf("list snapshots: %w", err)
		}
		fmt.Fprint(cmd.OutOrStdout(), snapshot.FormatList(entries))
		return nil
	},
}

func init() {
	defaultDir := os.Getenv("ENVOY_SNAPSHOT_DIR")
	if defaultDir == "" {
		defaultDir = ".envoy-snapshots"
	}
	snapshotCmd.PersistentFlags().StringVar(&snapshotDir, "dir", defaultDir, "directory to store snapshots")
	snapshotCmd.AddCommand(snapshotSaveCmd)
	snapshotCmd.AddCommand(snapshotListCmd)
	rootCmd.AddCommand(snapshotCmd)
}
