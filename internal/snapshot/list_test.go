package snapshot_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envoy-diff/internal/snapshot"
)

func TestList_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	entries, err := snapshot.List(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestList_MultipleSnapshots(t *testing.T) {
	dir := t.TempDir()
	env := map[string]string{"A": "1"}
	for _, label := range []string{"dev", "staging", "prod"} {
		s := snapshot.New(label, env)
		if _, err := snapshot.Save(dir, s); err != nil {
			t.Fatalf("Save(%s) failed: %v", label, err)
		}
	}
	entries, err := snapshot.List(dir)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestFormatList_Empty(t *testing.T) {
	out := snapshot.FormatList(nil)
	if !strings.Contains(out, "No snapshots") {
		t.Errorf("expected 'No snapshots' message, got: %q", out)
	}
}

func TestFormatList_WithEntries(t *testing.T) {
	entries := []snapshot.Entry{
		{Label: "staging", Timestamp: "2024-01-15T12:00:00Z", Path: "/tmp/staging.json"},
		{Label: "prod", Timestamp: "2024-01-16T08:30:00Z", Path: "/tmp/prod.json"},
	}
	out := snapshot.FormatList(entries)
	if !strings.Contains(out, "staging") {
		t.Errorf("expected 'staging' in output, got: %q", out)
	}
	if !strings.Contains(out, "prod") {
		t.Errorf("expected 'prod' in output, got: %q", out)
	}
	if !strings.Contains(out, "LABEL") {
		t.Errorf("expected header row in output")
	}
}
