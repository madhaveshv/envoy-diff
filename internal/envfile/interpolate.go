package envfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var varPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// InterpolateOptions controls how variable interpolation is performed.
type InterpolateOptions struct {
	// UseOSEnv allows falling back to the OS environment when a variable is
	// not found in the provided map.
	UseOSEnv bool
	// FailOnMissing returns an error if a referenced variable cannot be resolved.
	FailOnMissing bool
}

// Interpolate replaces variable references in env values with their resolved
// values. References can be in the form $VAR or ${VAR}. The provided map is
// consulted first; if UseOSEnv is set, os.Getenv is used as a fallback.
func Interpolate(env map[string]string, opts InterpolateOptions) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		resolved, err := interpolateValue(v, env, opts)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

func interpolateValue(value string, env map[string]string, opts InterpolateOptions) (string, error) {
	var lastErr error
	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		name := strings.TrimPrefix(match, "$")
		name = strings.TrimPrefix(name, "{")
		name = strings.TrimSuffix(name, "}")

		if val, ok := env[name]; ok {
			return val
		}
		if opts.UseOSEnv {
			if val, ok := os.LookupEnv(name); ok {
				return val
			}
		}
		if opts.FailOnMissing {
			lastErr = fmt.Errorf("undefined variable: %s", name)
		}
		return match
	})
	if lastErr != nil {
		return "", lastErr
	}
	return result, nil
}
