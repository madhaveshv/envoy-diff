package envfile

import (
	"regexp"
	"strings"
)

const redactedValue = "[REDACTED]"

// RedactOptions controls which keys are redacted.
type RedactOptions struct {
	// SensitivePatterns is a list of regex patterns matched against key names.
	// Keys matching any pattern will have their values replaced with [REDACTED].
	SensitivePatterns []*regexp.Regexp
}

// DefaultSensitivePatterns returns a set of common patterns for sensitive keys.
func DefaultSensitivePatterns() []*regexp.Regexp {
	raw := []string{
		`(?i)password`,
		`(?i)secret`,
		`(?i)token`,
		`(?i)api[_-]?key`,
		`(?i)private[_-]?key`,
		`(?i)auth`,
		`(?i)credential`,
		`(?i)passphrase`,
	}
	patterns := make([]*regexp.Regexp, 0, len(raw))
	for _, r := range raw {
		patterns = append(patterns, regexp.MustCompile(r))
	}
	return patterns
}

// Redact returns a copy of env with sensitive values replaced by [REDACTED].
// If opts is nil, DefaultSensitivePatterns are used.
func Redact(env map[string]string, opts *RedactOptions) map[string]string {
	patterns := DefaultSensitivePatterns()
	if opts != nil {
		patterns = opts.SensitivePatterns
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if isSensitive(k, patterns) {
			out[k] = redactedValue
		} else {
			out[k] = v
		}
	}
	return out
}

func isSensitive(key string, patterns []*regexp.Regexp) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if p.MatchString(upper) || p.MatchString(key) {
			return true
		}
	}
	return false
}
