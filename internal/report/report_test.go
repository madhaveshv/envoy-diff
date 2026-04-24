package report_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envoy-diff/internal/audit"
	"github.com/user/envoy-diff/internal/report"
)

var from = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"API_KEY":     "abc123",
}

var to = map[string]string{
	"DB_HOST":     "prod-db.example.com",
	"DB_PASSWORD": "newsecret",
	"NEW_VAR":     "value",
}

func TestNew_HasChanges(t *testing.T) {
	r := report.New("staging.env", "prod.env", from, to, audit.DefaultRules())
	if !r.HasChanges() {
		t.Error("expected report to have changes")
	}
}

func TestNew_HasIssues(t *testing.T) {
	r := report.New("staging.env", "prod.env", from, to, audit.DefaultRules())
	if !r.HasIssues() {
		t.Error("expected report to have audit issues for sensitive key changes")
	}
}

func TestNew_NoChanges(t *testing.T) {
	r := report.New("a.env", "b.env", from, from, audit.DefaultRules())
	if r.HasChanges() {
		t.Error("expected no changes when envs are identical")
	}
	if r.HasIssues() {
		t.Error("expected no issues when envs are identical")
	}
}

func TestFormatText_ContainsHeaders(t *testing.T) {
	r := report.New("staging.env", "prod.env", from, to, audit.DefaultRules())
	out := report.FormatText(r)
	for _, want := range []string{"envoy-diff", "staging.env", "prod.env", "Diff", "Audit"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestFormatJSON_ValidJSON(t *testing.T) {
	r := report.New("staging.env", "prod.env", from, to, audit.DefaultRules())
	out, err := report.FormatJSON(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := parsed["summary"]; !ok {
		t.Error("expected 'summary' key in JSON output")
	}
}
