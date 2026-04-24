package audit_test

import (
	"testing"

	"github.com/user/envoy-diff/internal/audit"
	"github.com/user/envoy-diff/internal/diff"
)

func TestAuditor_SensitiveKeyRemoved(t *testing.T) {
	a := audit.New(nil)
	entries := []diff.Entry{
		{Key: "DB_PASSWORD", Status: diff.Removed, OldValue: "secret"},
	}
	report := a.Run(entries)
	if !report.HasIssues() {
		t.Fatal("expected findings, got none")
	}
	if report.CountBySeverity(audit.SeverityHigh) != 1 {
		t.Errorf("expected 1 HIGH finding, got %d", report.CountBySeverity(audit.SeverityHigh))
	}
}

func TestAuditor_SensitiveKeyChanged(t *testing.T) {
	a := audit.New(nil)
	entries := []diff.Entry{
		{Key: "API_TOKEN", Status: diff.Modified, OldValue: "old", NewValue: "new"},
	}
	report := a.Run(entries)
	if report.CountBySeverity(audit.SeverityHigh) != 1 {
		t.Errorf("expected 1 HIGH finding, got %d", report.CountBySeverity(audit.SeverityHigh))
	}
}

func TestAuditor_EmptyValueIntroduced(t *testing.T) {
	a := audit.New(nil)
	entries := []diff.Entry{
		{Key: "LOG_LEVEL", Status: diff.Added, NewValue: ""},
	}
	report := a.Run(entries)
	if report.CountBySeverity(audit.SeverityMedium) != 1 {
		t.Errorf("expected 1 MEDIUM finding, got %d", report.CountBySeverity(audit.SeverityMedium))
	}
}

func TestAuditor_NoIssues(t *testing.T) {
	a := audit.New(nil)
	entries := []diff.Entry{
		{Key: "APP_ENV", Status: diff.Added, NewValue: "production"},
		{Key: "PORT", Status: diff.Modified, OldValue: "8080", NewValue: "9090"},
	}
	report := a.Run(entries)
	if report.HasIssues() {
		t.Errorf("expected no findings, got %+v", report.Findings)
	}
}

func TestAuditor_CustomRules(t *testing.T) {
	customRule := audit.Rule{
		Name: "no-debug",
		Check: func(e diff.Entry) *audit.Finding {
			if e.NewValue == "debug" {
				return &audit.Finding{
					Key:      e.Key,
					Severity: audit.SeverityLow,
					Message:  "debug mode enabled",
				}
			}
			return nil
		},
	}
	a := audit.New([]audit.Rule{customRule})
	entries := []diff.Entry{
		{Key: "LOG_LEVEL", Status: diff.Added, NewValue: "debug"},
	}
	report := a.Run(entries)
	if report.CountBySeverity(audit.SeverityLow) != 1 {
		t.Errorf("expected 1 LOW finding, got %d", report.CountBySeverity(audit.SeverityLow))
	}
}
