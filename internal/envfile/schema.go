package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaField defines the expected shape of a single env variable.
type SchemaField struct {
	Key      string
	Required bool
	Pattern  string // optional regex pattern for the value
}

// SchemaResult holds a single schema violation.
type SchemaResult struct {
	Key     string
	Message string
}

// ValidateSchema checks an env map against a slice of SchemaFields.
// It returns a list of violations (missing required keys, pattern mismatches).
func ValidateSchema(env map[string]string, fields []SchemaField) []SchemaResult {
	var results []SchemaResult

	for _, field := range fields {
		val, ok := env[field.Key]

		if field.Required && (!ok || strings.TrimSpace(val) == "") {
			results = append(results, SchemaResult{
				Key:     field.Key,
				Message: "required key is missing or empty",
			})
			continue
		}

		if ok && field.Pattern != "" {
			re, err := regexp.Compile(field.Pattern)
			if err != nil {
				results = append(results, SchemaResult{
					Key:     field.Key,
					Message: fmt.Sprintf("invalid pattern %q: %v", field.Pattern, err),
				})
				continue
			}
			if !re.MatchString(val) {
				results = append(results, SchemaResult{
					Key:     field.Key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, field.Pattern),
				})
			}
		}
	}

	return results
}
