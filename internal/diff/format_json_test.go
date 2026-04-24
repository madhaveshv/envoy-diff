package diff

import (
	"encoding/json"
	"testing"
)

func TestFormatJSON_Empty(t *testing.T) {
	out, err := FormatJSON([]Result{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result JSONResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if result.Summary.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Summary.Total)
	}
	if len(result.Changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(result.Changes))
	}
}

func TestFormatJSON_Added(t *testing.T) {
	results := []Result{
		{Key: "NEW_VAR", Type: Added, NewValue: "hello"},
	}
	out, err := FormatJSON(results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result JSONResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Summary.Added != 1 {
		t.Errorf("expected 1 added, got %d", result.Summary.Added)
	}
	if result.Changes[0].NewValue != "hello" {
		t.Errorf("expected new value 'hello', got %q", result.Changes[0].NewValue)
	}
	if result.Changes[0].OldValue != "" {
		t.Errorf("expected empty old value for added key")
	}
}

func TestFormatJSON_Removed(t *testing.T) {
	results := []Result{
		{Key: "OLD_VAR", Type: Removed, OldValue: "bye"},
	}
	out, err := FormatJSON(results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result JSONResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", result.Summary.Removed)
	}
	if result.Changes[0].OldValue != "bye" {
		t.Errorf("expected old value 'bye', got %q", result.Changes[0].OldValue)
	}
}

func TestFormatJSON_Mixed(t *testing.T) {
	results := []Result{
		{Key: "A", Type: Added, NewValue: "1"},
		{Key: "B", Type: Removed, OldValue: "2"},
		{Key: "C", Type: Modified, OldValue: "old", NewValue: "new"},
	}
	out, err := FormatJSON(results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result JSONResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Summary.Total != 3 {
		t.Errorf("expected total 3, got %d", result.Summary.Total)
	}
	if result.Summary.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", result.Summary.Modified)
	}
}
