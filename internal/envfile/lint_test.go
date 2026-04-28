package envfile

import (
	"testing"
)

func TestLint_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	issues := Lint(env, LintOptions{})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestLint_LeadingUnderscore(t *testing.T) {
	env := map[string]string{"_INTERNAL": "value"}
	issues := Lint(env, LintOptions{DisallowLeadingUnderscore: true})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "_INTERNAL" {
		t.Errorf("unexpected key: %s", issues[0].Key)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := map[string]string{"appName": "value"}
	issues := Lint(env, LintOptions{DisallowLowercase: true})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestLint_NumericStart(t *testing.T) {
	env := map[string]string{"1INVALID": "value"}
	issues := Lint(env, LintOptions{DisallowNumericStart: true})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestLint_MaxKeyLength(t *testing.T) {
	env := map[string]string{"VERY_LONG_KEY_NAME_EXCEEDING_LIMIT": "v"}
	issues := Lint(env, LintOptions{MaxKeyLength: 10})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestLint_MaxValueLength(t *testing.T) {
	env := map[string]string{"KEY": "this-value-is-too-long"}
	issues := Lint(env, LintOptions{MaxValueLength: 5})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestLint_MultipleRules(t *testing.T) {
	env := map[string]string{
		"_bad":    "value",
		"GOOD_KEY": "ok",
	}
	opts := LintOptions{
		DisallowLeadingUnderscore: true,
		DisallowLowercase:         true,
	}
	issues := Lint(env, opts)
	// _bad triggers both leading underscore and lowercase
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
}
