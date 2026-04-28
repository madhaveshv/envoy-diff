package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// LintOptions controls which lint rules are applied.
type LintOptions struct {
	DisallowLeadingUnderscore bool
	DisallowLowercase         bool
	DisallowNumericStart      bool
	MaxKeyLength              int
	MaxValueLength            int
}

// LintIssue represents a single lint finding.
type LintIssue struct {
	Key     string
	Message string
}

func (i LintIssue) String() string {
	return fmt.Sprintf("[%s] %s", i.Key, i.Message)
}

var numericStart = regexp.MustCompile(`^[0-9]`)

// Lint checks the provided env map against the given options and returns
// a slice of LintIssues. An empty slice means no issues were found.
func Lint(env map[string]string, opts LintOptions) []LintIssue {
	var issues []LintIssue

	for k, v := range env {
		if opts.DisallowLeadingUnderscore && strings.HasPrefix(k, "_") {
			issues = append(issues, LintIssue{Key: k, Message: "key starts with underscore"})
		}
		if opts.DisallowLowercase && k != strings.ToUpper(k) {
			issues = append(issues, LintIssue{Key: k, Message: "key contains lowercase characters"})
		}
		if opts.DisallowNumericStart && numericStart.MatchString(k) {
			issues = append(issues, LintIssue{Key: k, Message: "key starts with a numeric character"})
		}
		if opts.MaxKeyLength > 0 && len(k) > opts.MaxKeyLength {
			issues = append(issues, LintIssue{Key: k, Message: fmt.Sprintf("key length %d exceeds maximum %d", len(k), opts.MaxKeyLength)})
		}
		if opts.MaxValueLength > 0 && len(v) > opts.MaxValueLength {
			issues = append(issues, LintIssue{Key: k, Message: fmt.Sprintf("value length %d exceeds maximum %d", len(v), opts.MaxValueLength)})
		}
	}

	return issues
}
