package audit

import (
	"fmt"
	"io"
	"strings"
)

const (
	symbolHigh   = "[!!]"
	symbolMedium = "[! ]"
	symbolLow    = "[  ]"
)

func severitySymbol(s Severity) string {
	switch s {
	case SeverityHigh:
		return symbolHigh
	case SeverityMedium:
		return symbolMedium
	default:
		return symbolLow
	}
}

// FormatText writes a human-readable audit report to w.
func FormatText(w io.Writer, report Report) {
	if !report.HasIssues() {
		fmt.Fprintln(w, "audit: no issues found")
		return
	}

	fmt.Fprintf(w, "audit: %d issue(s) found\n", len(report.Findings))
	fmt.Fprintln(w, strings.Repeat("-", 50))
	for _, f := range report.Findings {
		fmt.Fprintf(w, "%s %-30s %s\n", severitySymbol(f.Severity), f.Key, f.Message)
	}
	fmt.Fprintln(w, strings.Repeat("-", 50))
	fmt.Fprintf(w, "  HIGH: %d  MEDIUM: %d  LOW: %d\n",
		report.CountBySeverity(SeverityHigh),
		report.CountBySeverity(SeverityMedium),
		report.CountBySeverity(SeverityLow),
	)
}
