package audit

import "github.com/user/envoy-diff/internal/diff"

// Report holds all findings produced by an audit run.
type Report struct {
	Findings []Finding
}

// HasIssues returns true when the report contains at least one finding.
func (r *Report) HasIssues() bool {
	return len(r.Findings) > 0
}

// CountBySeverity returns the number of findings for the given severity.
func (r *Report) CountBySeverity(s Severity) int {
	count := 0
	for _, f := range r.Findings {
		if f.Severity == s {
			count++
		}
	}
	return count
}

// Auditor runs a set of rules against a diff result.
type Auditor struct {
	rules []Rule
}

// New creates an Auditor with the provided rules.
// Pass nil to use DefaultRules.
func New(rules []Rule) *Auditor {
	if rules == nil {
		rules = DefaultRules()
	}
	return &Auditor{rules: rules}
}

// Run evaluates all rules against every entry in entries and returns a Report.
func (a *Auditor) Run(entries []diff.Entry) Report {
	var report Report
	for _, entry := range entries {
		for _, rule := range a.rules {
			if finding := rule.Check(entry); finding != nil {
				report.Findings = append(report.Findings, *finding)
			}
		}
	}
	return report
}
