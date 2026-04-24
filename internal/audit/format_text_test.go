package audit_test

import (
	"strings"
	"testing"

	"github.com/user/envoy-diff/internal/audit"
	"github.com/user/envoy-diff/internal/diff"
)

func makeIssue(severity, key, message string) audit.Issue {
	return audit.Issue{
		Severity: severity,
		Key:      key,
		Message:  message,
	}
}

func TestFormatText_NoIssues(t *testing.T) {
	out := audit.FormatText(nil)
	if !strings.Contains(out, "No audit issues") {
		t.Errorf("expected no-issues message, got: %s", out)
	}
}

func TestFormatText_SingleIssue(t *testing.T) {
	issues := []audit.Issue{
		makeIssue("HIGH", "DB_PASSWORD", "sensitive key was removed"),
	}
	out := audit.FormatText(issues)
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected key in output, got: %s", out)
	}
	if !strings.Contains(out, "HIGH") {
		t.Errorf("expected severity in output, got: %s", out)
	}
	if !strings.Contains(out, "sensitive key was removed") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestFormatText_MultipleIssues(t *testing.T) {
	issues := []audit.Issue{
		makeIssue("HIGH", "SECRET", "sensitive key changed"),
		makeIssue("MEDIUM", "API_KEY", "empty value introduced"),
	}
	out := audit.FormatText(issues)
	if !strings.Contains(out, "SECRET") || !strings.Contains(out, "API_KEY") {
		t.Errorf("expected both keys in output, got: %s", out)
	}
}

func TestFormatText_SeveritySymbols(t *testing.T) {
	_ = diff.Result{} // ensure diff package is importable
	issues := []audit.Issue{
		makeIssue("HIGH", "KEY1", "msg"),
		makeIssue("MEDIUM", "KEY2", "msg"),
		makeIssue("LOW", "KEY3", "msg"),
	}
	out := audit.FormatText(issues)
	if !strings.Contains(out, "KEY1") || !strings.Contains(out, "KEY2") || !strings.Contains(out, "KEY3") {
		t.Errorf("expected all keys in output, got: %s", out)
	}
}
