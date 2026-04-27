package envfile

import (
	"strings"
)

// NormalizeOptions controls how environment variable keys and values are normalized.
type NormalizeOptions struct {
	// UppercaseKeys converts all keys to uppercase.
	UppercaseKeys bool
	// TrimValues removes leading and trailing whitespace from values.
	TrimValues bool
	// RemoveEmpty drops entries with empty values.
	RemoveEmpty bool
	// KeyPrefix adds a prefix to all keys (applied after uppercasing).
	KeyPrefix string
}

// Normalize applies normalization rules to a map of environment variables
// and returns a new map with the transformations applied.
func Normalize(env map[string]string, opts NormalizeOptions) map[string]string {
	result := make(map[string]string, len(env))

	for k, v := range env {
		key := k
		value := v

		if opts.TrimValues {
			value = strings.TrimSpace(value)
		}

		if opts.RemoveEmpty && value == "" {
			continue
		}

		if opts.UppercaseKeys {
			key = strings.ToUpper(key)
		}

		if opts.KeyPrefix != "" {
			key = opts.KeyPrefix + key
		}

		result[key] = value
	}

	return result
}
