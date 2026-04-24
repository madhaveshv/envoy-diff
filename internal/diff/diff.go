package diff

// ChangeType represents the type of change for an environment variable.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single environment variable change between two stages.
type Change struct {
	Key      string
	OldValue string
	NewValue string
	Type     ChangeType
}

// Result holds the full diff result between two env maps.
type Result struct {
	Changes []Change
}

// Added returns only the added changes.
func (r *Result) Added() []Change {
	return r.filter(Added)
}

// Removed returns only the removed changes.
func (r *Result) Removed() []Change {
	return r.filter(Removed)
}

// Modified returns only the modified changes.
func (r *Result) Modified() []Change {
	return r.filter(Modified)
}

func (r *Result) filter(ct ChangeType) []Change {
	var out []Change
	for _, c := range r.Changes {
		if c.Type == ct {
			out = append(out, c)
		}
	}
	return out
}

// Compare computes the diff between a base env map and a target env map.
func Compare(base, target map[string]string) *Result {
	result := &Result{}

	for key, baseVal := range base {
		if targetVal, ok := target[key]; ok {
			if baseVal != targetVal {
				result.Changes = append(result.Changes, Change{
					Key:      key,
					OldValue: baseVal,
					NewValue: targetVal,
					Type:     Modified,
				})
			}
		} else {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				OldValue: baseVal,
				NewValue: "",
				Type:     Removed,
			})
		}
	}

	for key, targetVal := range target {
		if _, ok := base[key]; !ok {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				OldValue: "",
				NewValue: targetVal,
				Type:     Added,
			})
		}
	}

	return result
}
