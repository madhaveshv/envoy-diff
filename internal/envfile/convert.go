package envfile

import (
	"fmt"
	"strings"
)

// ConvertOptions controls how an env map is transformed during conversion.
type ConvertOptions struct {
	UppercaseKeys bool
	LowercaseKeys bool
	KeyPrefix     string
	StripPrefix   string
	QuoteStyle    string // none, single, double
	TrimValues    bool
}

// Convert applies a set of transformation options to an env map and returns a new map.
func Convert(env map[string]string, opts ConvertOptions) (map[string]string, error) {
	if err := validateConvertOptions(opts); err != nil {
		return nil, err
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		newKey := k

		if opts.StripPrefix != "" && strings.HasPrefix(newKey, opts.StripPrefix) {
			newKey = strings.TrimPrefix(newKey, opts.StripPrefix)
		}

		if opts.KeyPrefix != "" {
			newKey = opts.KeyPrefix + newKey
		}

		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		} else if opts.LowercaseKeys {
			newKey = strings.ToLower(newKey)
		}

		newVal := v
		if opts.TrimValues {
			newVal = strings.TrimSpace(newVal)
		}

		newVal = applyQuoteStyle(newVal, opts.QuoteStyle)
		result[newKey] = newVal
	}
	return result, nil
}

func applyQuoteStyle(val, style string) string {
	switch style {
	case "single":
		return "'" + strings.ReplaceAll(val, "'", `\'`) + "'"
	case "double":
		return `"` + strings.ReplaceAll(val, `"`, `\"`) + `"`
	default:
		return val
	}
}

func validateConvertOptions(opts ConvertOptions) error {
	if opts.UppercaseKeys && opts.LowercaseKeys {
		return fmt.Errorf("convert: uppercase-keys and lowercase-keys are mutually exclusive")
	}
	valid := map[string]bool{"none": true, "single": true, "double": true, "": true}
	if !valid[opts.QuoteStyle] {
		return fmt.Errorf("convert: unsupported quote style %q", opts.QuoteStyle)
	}
	return nil
}
