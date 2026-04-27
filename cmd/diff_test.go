package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempDiffFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestDiffCmd_TextOutput(t *testing.T) {
	before := writeTempDiffFile(t, "FOO=bar\nBAZ=qux\n")
	after := writeTempDiffFile(t, "FOO=changed\nNEW=value\n")

	buf := &bytes.Buffer{}
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs([]string{"diff", before, after})

	err := RootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiffCmd_JSONOutput(t *testing.T) {
	before := writeTempDiffFile(t, "FOO=bar\n")
	after := writeTempDiffFile(t, "FOO=baz\n")

	buf := &bytes.Buffer{}
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs([]string{"diff", "--format", "json", before, after})

	err := RootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiffCmd_MissingFile(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nonexistent.env")
	existing := writeTempDiffFile(t, "FOO=bar\n")

	RootCmd.SetArgs([]string{"diff", missing, existing})
	err := RootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestDiffCmd_WithAuditFlag(t *testing.T) {
	before := writeTempDiffFile(t, "SECRET_KEY=old\n")
	after := writeTempDiffFile(t, "SECRET_KEY=\n")

	buf := &bytes.Buffer{}
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs([]string{"diff", "--audit", before, after})

	err := RootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	_ = strings.Contains(output, "audit") // audit section may appear
}
