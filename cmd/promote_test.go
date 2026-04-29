package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempPromoteFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "env.env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return p
}

func TestPromoteCmd_BasicPromotion(t *testing.T) {
	src := writeTempPromoteFile(t, "A=1\nB=2\n")
	dst := writeTempPromoteFile(t, "C=3\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"promote", src, dst})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, key := range []string{"A", "B", "C"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected %q in output, got: %s", key, out)
		}
	}
}

func TestPromoteCmd_OnlyFlag(t *testing.T) {
	src := writeTempPromoteFile(t, "A=1\nB=2\n")
	dst := writeTempPromoteFile(t, "")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"promote", "--only=A", src, dst})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "A=1") {
		t.Errorf("expected A=1 in output")
	}
	if strings.Contains(out, "B") {
		t.Errorf("expected B to be excluded")
	}
}

func TestPromoteCmd_MissingFile(t *testing.T) {
	rootCmd.SetArgs([]string{"promote", "nonexistent.env", "also-missing.env"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestPromoteCmd_FailOnMissing(t *testing.T) {
	src := writeTempPromoteFile(t, "A=1\n")
	dst := writeTempPromoteFile(t, "")

	rootCmd.SetArgs([]string{"promote", "--only=MISSING", "--fail-on-missing", src, dst})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for missing key")
	}
}
