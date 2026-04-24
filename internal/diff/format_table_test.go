package diff_test

import (
	"strings"
	"testing"

	"github.com/user/envoy-diff/internal/diff"
)

func TestFormatTable_Empty(t *testing.T) {
	result := diff.FormatTable([]diff.Result{})
	if !strings.Contains(result, "No differences found") {
		t.Errorf("expected 'No differences found', got: %s", result)
	}
}

func TestFormatTable_Added(t *testing.T) {
	results := []diff.Result{
		{Key: "NEW_KEY", Status: diff.Added, ToValue: "hello"},
	}
	result := diff.FormatTable(results)
	if !strings.Contains(result, "NEW_KEY") {
		t.Errorf("expected NEW_KEY in output, got: %s", result)
	}
	if !strings.Contains(result, "added") {
		t.Errorf("expected 'added' status in output, got: %s", result)
	}
}

func TestFormatTable_Removed(t *testing.T) {
	results := []diff.Result{
		{Key: "OLD_KEY", Status: diff.Removed, FromValue: "world"},
	}
	result := diff.FormatTable(results)
	if !strings.Contains(result, "OLD_KEY") {
		t.Errorf("expected OLD_KEY in output, got: %s", result)
	}
	if !strings.Contains(result, "removed") {
		t.Errorf("expected 'removed' status in output, got: %s", result)
	}
}

func TestFormatTable_Modified(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.Modified, FromValue: "localhost", ToValue: "prod.db.internal"},
	}
	result := diff.FormatTable(results)
	if !strings.Contains(result, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", result)
	}
	if !strings.Contains(result, "modified") {
		t.Errorf("expected 'modified' status in output, got: %s", result)
	}
	if !strings.Contains(result, "localhost") {
		t.Errorf("expected from-value in output, got: %s", result)
	}
	if !strings.Contains(result, "prod.db.internal") {
		t.Errorf("expected to-value in output, got: %s", result)
	}
}

func TestFormatTable_TruncatesLongValues(t *testing.T) {
	long := strings.Repeat("x", 100)
	results := []diff.Result{
		{Key: "SECRET", Status: diff.Added, ToValue: long},
	}
	result := diff.FormatTable(results)
	if strings.Contains(result, long) {
		t.Errorf("expected long value to be truncated in table output")
	}
}

func TestFormatTable_HeaderPresent(t *testing.T) {
	results := []diff.Result{
		{Key: "FOO", Status: diff.Added, ToValue: "bar"},
	}
	result := diff.FormatTable(results)
	if !strings.Contains(result, "KEY") {
		t.Errorf("expected header row with KEY column, got: %s", result)
	}
	if !strings.Contains(result, "STATUS") {
		t.Errorf("expected header row with STATUS column, got: %s", result)
	}
}
