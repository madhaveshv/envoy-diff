package envfile

import "strings"

// FilterOptions controls which keys are included in the output.
type FilterOptions struct {
	// Prefix filters keys to only those starting with the given prefix.
	Prefix string
	// Keys is an explicit allowlist of keys to include. If empty, all keys pass.
	Keys []string
	// ExcludeKeys is an explicit denylist of keys to exclude.
	ExcludeKeys []string
}

// Filter returns a new map containing only the entries that match the options.
func Filter(env map[string]string, opts FilterOptions) map[string]string {
	excluded := make(map[string]struct{}, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excluded[k] = struct{}{}
	}

	allowed := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		allowed[k] = struct{}{}
	}

	result := make(map[string]string)
	for k, v := range env {
		if _, skip := excluded[k]; skip {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		if len(allowed) > 0 {
			if _, ok := allowed[k]; !ok {
				continue
			}
		}
		result[k] = v
	}
	return result
}
