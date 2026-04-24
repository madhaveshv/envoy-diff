package report

import (
	"encoding/json"
)

type jsonReport struct {
	FromFile string        `json:"from_file"`
	ToFile   string        `json:"to_file"`
	Diffs    []jsonDiff    `json:"diffs"`
	Issues   []jsonIssue   `json:"issues"`
	Summary  jsonSummary   `json:"summary"`
}

type jsonDiff struct {
	Key      string `json:"key"`
	Status   string `json:"status"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
}

type jsonIssue struct {
	Severity string `json:"severity"`
	Key      string `json:"key"`
	Message  string `json:"message"`
}

type jsonSummary struct {
	TotalDiffs  int `json:"total_diffs"`
	TotalIssues int `json:"total_issues"`
}

// FormatJSON renders the full report as a JSON string.
func FormatJSON(r *Report) (string, error) {
	diffs := make([]jsonDiff, 0, len(r.Diffs))
	for _, d := range r.Diffs {
		diffs = append(diffs, jsonDiff{
			Key:      d.Key,
			Status:   string(d.Status),
			OldValue: d.OldValue,
			NewValue: d.NewValue,
		})
	}

	issues := make([]jsonIssue, 0, len(r.Issues))
	for _, i := range r.Issues {
		issues = append(issues, jsonIssue{
			Severity: i.Severity,
			Key:      i.Key,
			Message:  i.Message,
		})
	}

	payload := jsonReport{
		FromFile: r.FromFile,
		ToFile:   r.ToFile,
		Diffs:    diffs,
		Issues:   issues,
		Summary: jsonSummary{
			TotalDiffs:  len(diffs),
			TotalIssues: len(issues),
		},
	}

	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
