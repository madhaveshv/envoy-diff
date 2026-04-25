package envfile

import "strings"

// FilterFlagsFromArgs parses CLI-style flag strings into a FilterOptions struct.
// Supported inputs:
//   - prefix=<value>  sets the key prefix filter
//   - key=<k1,k2>    sets the allowlist of keys
//   - exclude=<k1,k2> sets the denylist of keys
func FilterFlagsFromArgs(prefix, keys, exclude string) FilterOptions {
	opts := FilterOptions{}

	if prefix != "" {
		opts.Prefix = prefix
	}

	if keys != "" {
		for _, k := range strings.Split(keys, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				opts.Keys = append(opts.Keys, k)
			}
		}
	}

	if exclude != "" {
		for _, k := range strings.Split(exclude, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				opts.ExcludeKeys = append(opts.ExcludeKeys, k)
			}
		}
	}

	return opts
}
