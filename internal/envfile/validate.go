package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationIssue represents a single validation problem found in an env map.
type ValidationIssue struct {
	Key     string
	Message string
	Severity string // "error" | "warning"
}

func (v ValidationIssue) String() string {
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(v.Severity), v.Key, v.Message)
}

// ValidateOptions controls which validation rules are applied.
type ValidateOptions struct {
	RequireUppercase bool
	DisallowEmpty    bool
	ForbiddenKeys    []string
	KeyPattern       string // optional regex that all keys must match
}

var defaultKeyPattern = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// Validate checks an env map against the provided options and returns any issues found.
func Validate(env map[string]string, opts ValidateOptions) []ValidationIssue {
	var issues []ValidationIssue

	forbidden := make(map[string]struct{}, len(opts.ForbiddenKeys))
	for _, k := range opts.ForbiddenKeys {
		forbidden[k] = struct{}{}
	}

	var keyRe *regexp.Regexp
	if opts.KeyPattern != "" {
		var err error
		keyRe, err = regexp.Compile(opts.KeyPattern)
		if err != nil {
			issues = append(issues, ValidationIssue{
				Key:      "<options>",
				Message:  fmt.Sprintf("invalid key_pattern regex: %v", err),
				Severity: "error",
			})
		}
	}

	for k, v := range env {
		if opts.RequireUppercase && !defaultKeyPattern.MatchString(k) {
			issues = append(issues, ValidationIssue{
				Key:      k,
				Message:  "key does not match expected uppercase format [A-Z_][A-Z0-9_]*",
				Severity: "warning",
			})
		}
		if opts.DisallowEmpty && v == "" {
			issues = append(issues, ValidationIssue{
				Key:      k,
				Message:  "value is empty",
				Severity: "warning",
			})
		}
		if _, found := forbidden[k]; found {
			issues = append(issues, ValidationIssue{
				Key:      k,
				Message:  "key is explicitly forbidden",
				Severity: "error",
			})
		}
		if keyRe != nil && !keyRe.MatchString(k) {
			issues = append(issues, ValidationIssue{
				Key:      k,
				Message:  fmt.Sprintf("key does not match required pattern %q", opts.KeyPattern),
				Severity: "warning",
			})
		}
	}
	return issues
}
