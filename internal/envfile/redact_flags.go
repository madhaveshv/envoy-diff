package envfile

import "regexp"

// RedactFlagsFromArgs builds a RedactOptions from CLI flag values.
// If noRedact is true, nil is returned (disabling all redaction).
// extraPatterns is a slice of raw regex strings supplied by the user.
func RedactFlagsFromArgs(noRedact bool, extraPatterns []string) (*RedactOptions, error) {
	if noRedact {
		return &RedactOptions{
			SensitivePatterns: []*regexp.Regexp{},
		}, nil
	}

	patterns := DefaultSensitivePatterns()

	for _, raw := range extraPatterns {
		re, err := regexp.Compile(raw)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, re)
	}

	return &RedactOptions{SensitivePatterns: patterns}, nil
}
