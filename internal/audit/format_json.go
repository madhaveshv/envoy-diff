package audit

import (
	"encoding/json"
)

type jsonIssue struct {
	Severity string `json:"severity"`
	Rule     string `json:"rule"`
	Key      string `json:"key"`
	Message  string `json:"message"`
}

type jsonAuditReport struct {
	Total  int          `json:"total"`
	Issues []jsonIssue  `json:"issues"`
}

// FormatJSON renders audit issues as a JSON string.
// Returns an empty report structure if there are no issues.
func FormatJSON(issues []Issue) (string, error) {
	report := jsonAuditReport{
		Total:  len(issues),
		Issues: make([]jsonIssue, 0, len(issues)),
	}

	for _, iss := range issues {
		report.Issues = append(report.Issues, jsonIssue{
			Severity: string(iss.Severity),
			Rule:     iss.Rule,
			Key:      iss.Key,
			Message:  iss.Message,
		})
	}

	bytes, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
