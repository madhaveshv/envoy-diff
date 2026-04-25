package envfile

import (
	"testing"
)

func TestMerge_NoOverwrite(t *testing.T) {
	dst := map[string]string{"A": "1", "B": "2"}
	src := map[string]string{"B": "99", "C": "3"}

	result := Merge(dst, src, MergeOptions{Overwrite: false})

	if result["A"] != "1" {
		t.Errorf("expected A=1, got %s", result["A"])
	}
	if result["B"] != "2" {
		t.Errorf("expected B=2 (no overwrite), got %s", result["B"])
	}
	if result["C"] != "3" {
		t.Errorf("expected C=3, got %s", result["C"])
	}
}

func TestMerge_WithOverwrite(t *testing.T) {
	dst := map[string]string{"A": "1", "B": "2"}
	src := map[string]string{"B": "99", "C": "3"}

	result := Merge(dst, src, MergeOptions{Overwrite: true})

	if result["B"] != "99" {
		t.Errorf("expected B=99 (overwrite), got %s", result["B"])
	}
	if result["C"] != "3" {
		t.Errorf("expected C=3, got %s", result["C"])
	}
}

func TestMerge_SkipEmpty(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"B": "", "C": "3"}

	result := Merge(dst, src, MergeOptions{Overwrite: true, SkipEmpty: true})

	if _, ok := result["B"]; ok {
		t.Error("expected B to be skipped due to empty value")
	}
	if result["C"] != "3" {
		t.Errorf("expected C=3, got %s", result["C"])
	}
}

func TestMerge_DoesNotMutateDst(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"A": "2", "B": "3"}

	_ = Merge(dst, src, MergeOptions{Overwrite: true})

	if dst["A"] != "1" {
		t.Errorf("dst was mutated: expected A=1, got %s", dst["A"])
	}
	if _, ok := dst["B"]; ok {
		t.Error("dst was mutated: B should not exist in original dst")
	}
}

func TestMerge_EmptyMaps(t *testing.T) {
	result := Merge(map[string]string{}, map[string]string{}, MergeOptions{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d keys", len(result))
	}
}
