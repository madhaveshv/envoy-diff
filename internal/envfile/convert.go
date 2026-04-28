package envfile

import (
	"fmt"
	"strings"
)

// ConvertOptions controls how environment variable keys are transformed.
type ConvertOptions struct {
	// KeyCase converts all keys to "upper" or "lower" case.
	KeyCase string
	// KeyPrefix adds a prefix to all keys.
	KeyPrefix string
	// StripPrefix removes a prefix from all keys (applied before adding KeyPrefix).
	StripPrefix string
	// ValueQuoteStyle wraps values in "double", 'single', or "none" (default) quotes.
	ValueQuoteStyle string
}

// Convert applies structural transformations to env map keys and values
// according to the provided ConvertOptions. It returns a new map and does
// not mutate the input.
func Convert(env map[string]string, opts ConvertOptions) (map[string]string, error) {
	if err := validateConvertOptions(opts); err != nil {
		return nil, err
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		newKey := k

		if opts.StripPrefix != "" {
			newKey = strings.TrimPrefix(newKey, opts.StripPrefix)
		}

		switch strings.ToLower(opts.KeyCase) {
		case "upper":
			newKey = strings.ToUpper(newKey)
		case "lower":
			newKey = strings.ToLower(newKey)
		}

		if opts.KeyPrefix != "" {
			newKey = opts.KeyPrefix + newKey
		}

		newVal := applyQuoteStyle(v, opts.ValueQuoteStyle)
		out[newKey] = newVal
	}
	return out, nil
}

func applyQuoteStyle(val, style string) string {
	switch strings.ToLower(style) {
	case "double":
		return fmt.Sprintf(`"%s"`, val)
	case "single":
		return fmt.Sprintf(`'%s'`, val)
	default:
		return val
	}
}

func validateConvertOptions(opts ConvertOptions) error {
	switch strings.ToLower(opts.KeyCase) {
	case "", "upper", "lower":
		// valid
	default:
		return fmt.Errorf("invalid key_case %q: must be \"upper\" or \"lower\"", opts.KeyCase)
	}
	switch strings.ToLower(opts.ValueQuoteStyle) {
	case "", "none", "double", "single":
		// valid
	default:
		return fmt.Errorf("invalid value_quote_style %q: must be \"double\", \"single\", or \"none\"", opts.ValueQuoteStyle)
	}
	return nil
}
