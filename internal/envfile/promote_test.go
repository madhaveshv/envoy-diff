package envfile

import (
	"testing"
)

func TestPromote_NoOptions(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"C": "3"}
	out, err := Promote(src, dst, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "2" || out["C"] != "3" {
		t.Errorf("unexpected result: %v", out)
	}
}

func TestPromote_DoesNotOverwriteByDefault(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, err := Promote(src, dst, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "old" {
		t.Errorf("expected old value to be preserved, got %q", out["A"])
	}
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, err := Promote(src, dst, PromoteOptions{OverwriteExisting: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "new" {
		t.Errorf("expected new value, got %q", out["A"])
	}
}

func TestPromote_OnlyKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}
	out, err := Promote(src, dst, PromoteOptions{OnlyKeys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["C"] != "3" {
		t.Errorf("expected A and C, got %v", out)
	}
	if _, ok := out["B"]; ok {
		t.Errorf("expected B to be excluded")
	}
}

func TestPromote_ExcludeKeys(t *testing.T) {
	src := map[string]string{"A": "1", "SECRET": "s"}
	dst := map[string]string{}
	out, err := Promote(src, dst, PromoteOptions{ExcludeKeys: []string{"SECRET"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["SECRET"]; ok {
		t.Errorf("expected SECRET to be excluded")
	}
	if out["A"] != "1" {
		t.Errorf("expected A to be present")
	}
}

func TestPromote_FailOnMissing(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{}
	_, err := Promote(src, dst, PromoteOptions{
		OnlyKeys:      []string{"A", "MISSING"},
		FailOnMissing: true,
	})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestPromote_DoesNotMutateDst(t *testing.T) {
	src := map[string]string{"X": "new"}
	dst := map[string]string{"X": "orig"}
	Promote(src, dst, PromoteOptions{OverwriteExisting: true})
	if dst["X"] != "orig" {
		t.Errorf("dst was mutated")
	}
}
