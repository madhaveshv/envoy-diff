package report

import (
	"fmt"
	"strings"

	auditfmt "github.com/user/envoy-diff/internal/audit"
	diffpkg "github.com/user/envoy-diff/internal/diff"
)

// FormatText renders the full report as human-readable text.
func FormatText(r *Report) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("=== envoy-diff report ===\n"))
	sb.WriteString(fmt.Sprintf("From: %s\n", r.FromFile))
	sb.WriteString(fmt.Sprintf("To:   %s\n\n", r.ToFile))

	sb.WriteString("--- Diff ---\n")
	sb.WriteString(diffpkg.FormatText(r.Diffs))
	sb.WriteString("\n")

	sb.WriteString("--- Audit ---\n")
	sb.WriteString(auditfmt.FormatText(r.Issues))

	return sb.String()
}

// FormatTable renders the diff section as a table and appends audit issues as text.
func FormatTable(r *Report) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("=== envoy-diff report ===\n"))
	sb.WriteString(fmt.Sprintf("From: %s\nTo:   %s\n\n", r.FromFile, r.ToFile))

	sb.WriteString("--- Diff ---\n")
	sb.WriteString(diffpkg.FormatTable(r.Diffs))
	sb.WriteString("\n")

	sb.WriteString("--- Audit ---\n")
	sb.WriteString(auditfmt.FormatText(r.Issues))

	return sb.String()
}
