package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_ValidEnvFile(t *testing.T) {
	path := writeTempEnvFile(t, "KEY=value\nFOO=bar\n")

	env, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", env["KEY"])
	}
	if env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", env["FOO"])
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_UnsupportedExtension(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.yaml")
	_ = os.WriteFile(path, []byte("KEY=value\n"), 0644)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for unsupported extension, got nil")
	}
}

func TestLoad_DirectoryPath(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := Load(tmpDir)
	if err == nil {
		t.Fatal("expected error for directory path, got nil")
	}
}

func TestLoadPair_BothValid(t *testing.T) {
	fromPath := writeTempEnvFile(t, "KEY=old\n")
	toPath := writeTempEnvFile(t, "KEY=new\n")

	fromEnv, toEnv, err := LoadPair(fromPath, toPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fromEnv["KEY"] != "old" {
		t.Errorf("expected from KEY=old, got %q", fromEnv["KEY"])
	}
	if toEnv["KEY"] != "new" {
		t.Errorf("expected to KEY=new, got %q", toEnv["KEY"])
	}
}

func TestLoadPair_FromFileMissing(t *testing.T) {
	toPath := writeTempEnvFile(t, "KEY=value\n")
	_, _, err := LoadPair("/nonexistent/.env", toPath)
	if err == nil {
		t.Fatal("expected error when from file is missing")
	}
}

func TestLoadPair_ToFileMissing(t *testing.T) {
	fromPath := writeTempEnvFile(t, "KEY=value\n")
	_, _, err := LoadPair(fromPath, "/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error when to file is missing")
	}
}
