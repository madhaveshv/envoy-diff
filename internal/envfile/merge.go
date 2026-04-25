package envfile

// MergeOptions controls how two env maps are merged.
type MergeOptions struct {
	// Overwrite controls whether values from src overwrite values in dst.
	Overwrite bool
	// SkipEmpty skips keys in src whose values are empty strings.
	SkipEmpty bool
}

// Merge combines two environment variable maps into a new map.
// The dst map is treated as the base; src is merged on top of it
// according to the provided MergeOptions.
//
// Neither dst nor src is mutated — a new map is always returned.
func Merge(dst, src map[string]string, opts MergeOptions) map[string]string {
	result := make(map[string]string, len(dst)+len(src))

	for k, v := range dst {
		result[k] = v
	}

	for k, v := range src {
		if opts.SkipEmpty && v == "" {
			continue
		}
		if _, exists := result[k]; exists && !opts.Overwrite {
			continue
		}
		result[k] = v
	}

	return result
}
