package audit

import (
	"encoding/json"
	"testing"
)

func TestFormatJSON_NoIssues(t *testing.T) {
	out, err := FormatJSON([]Issue{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var report map[string]interface{}
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if report["total"].(float64) != 0 {
		t.Errorf("expected total 0, got %v", report["total"])
	}
}

func TestFormatJSON_SingleIssue(t *testing.T) {
	issues := []Issue{
		{Severity: SeverityHigh, Rule: "sensitive-removed", Key: "DB_PASSWORD", Message: "sensitive key removed"},
	}

	out, err := FormatJSON(issues)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var report map[string]interface{}
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if report["total"].(float64) != 1 {
		t.Errorf("expected total 1, got %v", report["total"])
	}

	issuesList := report["issues"].([]interface{})
	if len(issuesList) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issuesList))
	}

	issue := issuesList[0].(map[string]interface{})
	if issue["key"] != "DB_PASSWORD" {
		t.Errorf("expected key DB_PASSWORD, got %v", issue["key"])
	}
	if issue["severity"] != string(SeverityHigh) {
		t.Errorf("expected severity %s, got %v", SeverityHigh, issue["severity"])
	}
}

func TestFormatJSON_MultipleIssues(t *testing.T) {
	issues := []Issue{
		{Severity: SeverityHigh, Rule: "sensitive-removed", Key: "SECRET_KEY", Message: "sensitive key removed"},
		{Severity: SeverityMedium, Rule: "empty-value", Key: "API_URL", Message: "empty value introduced"},
	}

	out, err := FormatJSON(issues)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var report map[string]interface{}
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if report["total"].(float64) != 2 {
		t.Errorf("expected total 2, got %v", report["total"])
	}

	issuesList := report["issues"].([]interface{})
	if len(issuesList) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issuesList))
	}
}
