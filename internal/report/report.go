// Package report combines diff and audit results into a unified report.
package report

import (
	"github.com/user/envoy-diff/internal/audit"
	"github.com/user/envoy-diff/internal/diff"
)

// Report holds the combined results of a diff and audit run.
type Report struct {
	FromFile string
	ToFile   string
	Diffs    []diff.Result
	Issues   []audit.Issue
}

// New creates a Report by running a diff between two env maps and auditing the results.
func New(fromFile, toFile string, from, to map[string]string, rules []audit.Rule) *Report {
	diffs := diff.Compare(from, to)
	auditor := audit.New(rules)
	issues := auditor.Audit(diffs)
	return &Report{
		FromFile: fromFile,
		ToFile:   toFile,
		Diffs:    diffs,
		Issues:   issues,
	}
}

// HasIssues returns true if the report contains any audit issues.
func (r *Report) HasIssues() bool {
	return len(r.Issues) > 0
}

// HasChanges returns true if the report contains any diff changes.
func (r *Report) HasChanges() bool {
	return len(r.Diffs) > 0
}
