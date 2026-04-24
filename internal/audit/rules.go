package audit

import (
	"regexp"
	"strings"

	"github.com/user/envoy-diff/internal/diff"
)

// Severity represents the risk level of an audit finding.
type Severity string

const (
	SeverityHigh   Severity = "HIGH"
	SeverityMedium Severity = "MEDIUM"
	SeverityLow    Severity = "LOW"
)

// Finding represents a single audit result for a variable change.
type Finding struct {
	Key      string
	Severity Severity
	Message  string
}

// Rule defines an audit rule applied to diff entries.
type Rule struct {
	Name    string
	Check   func(entry diff.Entry) *Finding
}

var sensitivePattern = regexp.MustCompile(
	`(?i)(secret|password|passwd|token|api_key|apikey|private_key|auth|credential)`,
)

// DefaultRules returns the built-in set of audit rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			Name: "sensitive-key-removed",
			Check: func(e diff.Entry) *Finding {
				if e.Status == diff.Removed && sensitivePattern.MatchString(e.Key) {
					return &Finding{
						Key:      e.Key,
						Severity: SeverityHigh,
						Message:  "sensitive variable was removed",
					}
				}
				return nil
			},
		},
		{
			Name: "sensitive-key-changed",
			Check: func(e diff.Entry) *Finding {
				if e.Status == diff.Modified && sensitivePattern.MatchString(e.Key) {
					return &Finding{
						Key:      e.Key,
						Severity: SeverityHigh,
						Message:  "sensitive variable value was changed",
					}
				}
				return nil
			},
		},
		{
			Name: "empty-value-introduced",
			Check: func(e diff.Entry) *Finding {
				if (e.Status == diff.Added || e.Status == diff.Modified) &&
					strings.TrimSpace(e.NewValue) == "" {
					return &Finding{
						Key:      e.Key,
						Severity: SeverityMedium,
						Message:  "variable introduced with an empty value",
					}
				}
				return nil
			},
		},
	}
}
