package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempInterpolateFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return p
}

func TestInterpolateCmd_NoFlag(t *testing.T) {
	p := writeTempInterpolateFile(t, "BASE=/app\nLOG=${BASE}/logs\n")
	out, err := executeCommand(RootCmd, "interpolate", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Without --interpolate flag, references should remain unexpanded
	if !strings.Contains(out, "${BASE}") {
		t.Errorf("expected raw reference in output, got: %s", out)
	}
}

func TestInterpolateCmd_WithFlag(t *testing.T) {
	p := writeTempInterpolateFile(t, "BASE=/app\nLOG_DIR=${BASE}/logs\n")
	out, err := executeCommand(RootCmd, "interpolate", "--interpolate", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "/app/logs") {
		t.Errorf("expected resolved value in output, got: %s", out)
	}
}

func TestInterpolateCmd_JSONFormat(t *testing.T) {
	p := writeTempInterpolateFile(t, "HOST=db\nDSN=postgres://${HOST}/mydb\n")
	out, err := executeCommand(RootCmd, "interpolate", "--interpolate", "--format", "json", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "{") {
		t.Errorf("expected JSON output, got: %s", out)
	}
	if !strings.Contains(out, "postgres://db/mydb") {
		t.Errorf("expected resolved DSN in JSON output, got: %s", out)
	}
}

func TestInterpolateCmd_MissingFile(t *testing.T) {
	_, err := executeCommand(RootCmd, "interpolate", "/nonexistent/path.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
