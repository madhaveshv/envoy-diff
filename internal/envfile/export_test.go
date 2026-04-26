package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRender_DotenvFormat(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Render(env, ExportOptions{Format: FormatDotenv, Sorted: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO=\"") {
		t.Errorf("expected FOO in output, got: %s", out)
	}
	if !strings.Contains(out, "BAZ=\"") {
		t.Errorf("expected BAZ in output, got: %s", out)
	}
}

func TestRender_ShellFormat(t *testing.T) {
	env := map[string]string{"MY_VAR": "hello"}
	out, err := Render(env, ExportOptions{Format: FormatShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "export MY_VAR=") {
		t.Errorf("expected shell export prefix, got: %s", out)
	}
}

func TestRender_JSONFormat(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	out, err := Render(env, ExportOptions{Format: FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"KEY"`) {
		t.Errorf("expected JSON key, got: %s", out)
	}
	if !strings.Contains(out, `"value"`) {
		t.Errorf("expected JSON value, got: %s", out)
	}
}

func TestRender_UnsupportedFormat(t *testing.T) {
	_, err := Render(map[string]string{}, ExportOptions{Format: "xml"})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported export format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRender_SortedOutput(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"}
	out, err := Render(env, ExportOptions{Format: FormatDotenv, Sorted: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	alphaIdx := strings.Index(out, "ALPHA")
	middleIdx := strings.Index(out, "MIDDLE")
	zebraIdx := strings.Index(out, "ZEBRA")
	if !(alphaIdx < middleIdx && middleIdx < zebraIdx) {
		t.Errorf("expected sorted output, got: %s", out)
	}
}

func TestExport_WritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")
	env := map[string]string{"EXPORTED": "yes"}
	if err := Export(env, path, ExportOptions{Format: FormatDotenv}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read exported file: %v", err)
	}
	if !strings.Contains(string(data), "EXPORTED=") {
		t.Errorf("expected EXPORTED in file, got: %s", string(data))
	}
}

func TestExport_EmptyPath(t *testing.T) {
	err := Export(map[string]string{}, "", ExportOptions{})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}
