package envfile

import (
	"regexp"
	"testing"
)

func TestRedact_DefaultPatterns(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD":   "hunter2",
		"API_KEY":       "abc123",
		"SECRET_TOKEN":  "xyz",
		"APP_NAME":      "myapp",
		"PORT":          "8080",
		"AUTH_HEADER":   "Bearer tok",
		"PRIVATE_KEY":   "-----BEGIN",
		"DATABASE_HOST": "localhost",
	}

	result := Redact(env, nil)

	sensitive := []string{"DB_PASSWORD", "API_KEY", "SECRET_TOKEN", "AUTH_HEADER", "PRIVATE_KEY"}
	for _, k := range sensitive {
		if result[k] != redactedValue {
			t.Errorf("expected %s to be redacted, got %q", k, result[k])
		}
	}

	plain := []string{"APP_NAME", "PORT", "DATABASE_HOST"}
	for _, k := range plain {
		if result[k] == redactedValue {
			t.Errorf("expected %s NOT to be redacted", k)
		}
	}
}

func TestRedact_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "secret123"}
	_ = Redact(env, nil)
	if env["DB_PASSWORD"] != "secret123" {
		t.Error("Redact mutated the input map")
	}
}

func TestRedact_CustomPatterns(t *testing.T) {
	env := map[string]string{
		"INTERNAL_CODE": "1234",
		"PUBLIC_URL":    "https://example.com",
	}

	opts := &RedactOptions{
		SensitivePatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)internal`),
		},
	}

	result := Redact(env, opts)

	if result["INTERNAL_CODE"] != redactedValue {
		t.Errorf("expected INTERNAL_CODE to be redacted")
	}
	if result["PUBLIC_URL"] == redactedValue {
		t.Errorf("expected PUBLIC_URL NOT to be redacted")
	}
}

func TestRedact_EmptyMap(t *testing.T) {
	result := Redact(map[string]string{}, nil)
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result))
	}
}

func TestDefaultSensitivePatterns_NotEmpty(t *testing.T) {
	patterns := DefaultSensitivePatterns()
	if len(patterns) == 0 {
		t.Error("expected at least one default sensitive pattern")
	}
}
