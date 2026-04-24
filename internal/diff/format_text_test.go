package diff_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
)

func TestFormatText_Empty(t *testing.T) {
	result := diff.FormatText([]diff.Result{})
	if strings.TrimSpace(result) != "No differences found." {
		t.Errorf("expected 'No differences found.', got %q", result)
	}
}

func TestFormatText_Added(t *testing.T) {
	results := []diff.Result{
		{Key: "NEW_VAR", Status: diff.Added, NewValue: "hello"},
	}
	out := diff.FormatText(results)
	if !strings.Contains(out, "+ NEW_VAR") {
		t.Errorf("expected added line, got: %s", out)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("expected value 'hello', got: %s", out)
	}
}

func TestFormatText_Removed(t *testing.T) {
	results := []diff.Result{
		{Key: "OLD_VAR", Status: diff.Removed, OldValue: "bye"},
	}
	out := diff.FormatText(results)
	if !strings.Contains(out, "- OLD_VAR") {
		t.Errorf("expected removed line, got: %s", out)
	}
	if !strings.Contains(out, "bye") {
		t.Errorf("expected value 'bye', got: %s", out)
	}
}

func TestFormatText_Modified(t *testing.T) {
	results := []diff.Result{
		{Key: "MOD_VAR", Status: diff.Modified, OldValue: "old", NewValue: "new"},
	}
	out := diff.FormatText(results)
	if !strings.Contains(out, "~ MOD_VAR") {
		t.Errorf("expected modified line, got: %s", out)
	}
	if !strings.Contains(out, "old") || !strings.Contains(out, "new") {
		t.Errorf("expected both old and new values, got: %s", out)
	}
}

func TestFormatText_TruncatesLongValues(t *testing.T) {
	long := strings.Repeat("x", 200)
	results := []diff.Result{
		{Key: "BIG_VAR", Status: diff.Added, NewValue: long},
	}
	out := diff.FormatText(results)
	if strings.Contains(out, long) {
		t.Errorf("expected long value to be truncated, but full value was present")
	}
	if !strings.Contains(out, "...") {
		t.Errorf("expected truncation indicator '...', got: %s", out)
	}
}
