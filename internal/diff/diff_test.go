package diff

import (
	"testing"
)

func TestCompare_Added(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	result := Compare(base, target)

	if len(result.Added()) != 1 {
		t.Fatalf("expected 1 added, got %d", len(result.Added()))
	}
	if result.Added()[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY, got %s", result.Added()[0].Key)
	}
}

func TestCompare_Removed(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD_KEY": "gone"}
	target := map[string]string{"FOO": "bar"}

	result := Compare(base, target)

	if len(result.Removed()) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(result.Removed()))
	}
	if result.Removed()[0].Key != "OLD_KEY" {
		t.Errorf("expected OLD_KEY, got %s", result.Removed()[0].Key)
	}
}

func TestCompare_Modified(t *testing.T) {
	base := map[string]string{"FOO": "old_value"}
	target := map[string]string{"FOO": "new_value"}

	result := Compare(base, target)

	if len(result.Modified()) != 1 {
		t.Fatalf("expected 1 modified, got %d", len(result.Modified()))
	}
	c := result.Modified()[0]
	if c.OldValue != "old_value" || c.NewValue != "new_value" {
		t.Errorf("unexpected values: old=%s new=%s", c.OldValue, c.NewValue)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	base := map[string]string{"FOO": "bar", "BAZ": "qux"}
	target := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := Compare(base, target)

	if len(result.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(result.Changes))
	}
}

func TestCompare_MixedChanges(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	target := map[string]string{"A": "1", "B": "changed", "D": "4"}

	result := Compare(base, target)

	if len(result.Added()) != 1 {
		t.Errorf("expected 1 added, got %d", len(result.Added()))
	}
	if len(result.Removed()) != 1 {
		t.Errorf("expected 1 removed, got %d", len(result.Removed()))
	}
	if len(result.Modified()) != 1 {
		t.Errorf("expected 1 modified, got %d", len(result.Modified()))
	}
}
