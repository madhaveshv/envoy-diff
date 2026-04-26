package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-diff/internal/snapshot"
)

func TestNew_CopiesEnv(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := snapshot.New("staging", env)
	env["FOO"] = "mutated"
	if s.Env["FOO"] != "bar" {
		t.Errorf("expected original value 'bar', got %q", s.Env["FOO"])
	}
	if s.Label != "staging" {
		t.Errorf("expected label 'staging', got %q", s.Label)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	env := map[string]string{"KEY": "value", "NUM": "42"}
	s := snapshot.New("prod", env)
	path, err := snapshot.Save(dir, s)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("saved file not found: %v", err)
	}
	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Label != "prod" {
		t.Errorf("expected label 'prod', got %q", loaded.Label)
	}
	if loaded.Env["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", loaded.Env["KEY"])
	}
}

func TestLoad_InvalidFile(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(bad, []byte("not json{"), 0o644)
	_, err := snapshot.Load(bad)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
