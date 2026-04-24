package diff

import (
	"encoding/json"
	"fmt"
)

// JSONResult represents the structured JSON output of a diff comparison.
type JSONResult struct {
	Summary JSONSummary `json:"summary"`
	Changes []JSONChange `json:"changes"`
}

// JSONSummary holds counts of each change type.
type JSONSummary struct {
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Modified int `json:"modified"`
	Total    int `json:"total"`
}

// JSONChange represents a single variable change in JSON format.
type JSONChange struct {
	Key      string `json:"key"`
	Type     string `json:"type"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
}

// FormatJSON formats a slice of Result entries as a JSON string.
func FormatJSON(results []Result) (string, error) {
	if len(results) == 0 {
		empty := JSONResult{
			Summary: JSONSummary{},
			Changes: []JSONChange{},
		}
		b, err := json.MarshalIndent(empty, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON: %w", err)
		}
		return string(b), nil
	}

	var changes []JSONChange
	summary := JSONSummary{}

	for _, r := range results {
		change := JSONChange{
			Key:  r.Key,
			Type: string(r.Type),
		}
		switch r.Type {
		case Added:
			change.NewValue = r.NewValue
			summary.Added++
		case Removed:
			change.OldValue = r.OldValue
			summary.Removed++
		case Modified:
			change.OldValue = r.OldValue
			change.NewValue = r.NewValue
			summary.Modified++
		}
		changes = append(changes, change)
	}

	summary.Total = summary.Added + summary.Removed + summary.Modified

	out := JSONResult{
		Summary: summary,
		Changes: changes,
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(b), nil
}
