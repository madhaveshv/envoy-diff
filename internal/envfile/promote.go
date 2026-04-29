package envfile

import "fmt"

// PromoteOptions controls how env vars are promoted between stages.
type PromoteOptions struct {
	// OnlyKeys restricts promotion to a specific set of keys.
	OnlyKeys []string
	// ExcludeKeys skips these keys during promotion.
	ExcludeKeys []string
	// OverwriteExisting allows values in dst to be overwritten.
	OverwriteExisting bool
	// FailOnMissing returns an error if any OnlyKeys key is absent from src.
	FailOnMissing bool
}

// Promote copies env vars from src into dst according to opts.
// It returns the resulting merged map and any validation error.
func Promote(src, dst map[string]string, opts PromoteOptions) (map[string]string, error) {
	excludeSet := make(map[string]bool, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excludeSet[k] = true
	}

	allowSet := make(map[string]bool, len(opts.OnlyKeys))
	for _, k := range opts.OnlyKeys {
		allowSet[k] = true
	}

	if opts.FailOnMissing {
		for _, k := range opts.OnlyKeys {
			if _, ok := src[k]; !ok {
				return nil, fmt.Errorf("promote: required key %q not found in source", k)
			}
		}
	}

	result := make(map[string]string, len(dst))
	for k, v := range dst {
		result[k] = v
	}

	for k, v := range src {
		if excludeSet[k] {
			continue
		}
		if len(allowSet) > 0 && !allowSet[k] {
			continue
		}
		if _, exists := result[k]; exists && !opts.OverwriteExisting {
			continue
		}
		result[k] = v
	}

	return result, nil
}
