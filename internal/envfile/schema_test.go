package envfile

import (
	"testing"
)

func TestValidateSchema_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"PORT":    "8080",
	}
	fields := []SchemaField{
		{Key: "APP_ENV", Required: true},
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
	}
	results := ValidateSchema(env, fields)
	if len(results) != 0 {
		t.Fatalf("expected no violations, got %d: %+v", len(results), results)
	}
}

func TestValidateSchema_MissingRequired(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	fields := []SchemaField{
		{Key: "APP_ENV", Required: true},
		{Key: "PORT", Required: true},
	}
	results := ValidateSchema(env, fields)
	if len(results) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(results))
	}
	if results[0].Key != "APP_ENV" {
		t.Errorf("expected violation for APP_ENV, got %s", results[0].Key)
	}
}

func TestValidateSchema_EmptyRequiredValue(t *testing.T) {
	env := map[string]string{"APP_ENV": "   "}
	fields := []SchemaField{{Key: "APP_ENV", Required: true}}
	results := ValidateSchema(env, fields)
	if len(results) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(results))
	}
}

func TestValidateSchema_PatternMismatch(t *testing.T) {
	env := map[string]string{"PORT": "not-a-number"}
	fields := []SchemaField{{Key: "PORT", Pattern: `^\d+$`}}
	results := ValidateSchema(env, fields)
	if len(results) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(results))
	}
	if results[0].Key != "PORT" {
		t.Errorf("expected violation for PORT, got %s", results[0].Key)
	}
}

func TestValidateSchema_PatternMatch(t *testing.T) {
	env := map[string]string{"PORT": "3000"}
	fields := []SchemaField{{Key: "PORT", Pattern: `^\d+$`}}
	results := ValidateSchema(env, fields)
	if len(results) != 0 {
		t.Fatalf("expected no violations, got %d", len(results))
	}
}

func TestValidateSchema_InvalidPattern(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	fields := []SchemaField{{Key: "KEY", Pattern: `[invalid`}}
	results := ValidateSchema(env, fields)
	if len(results) != 1 {
		t.Fatalf("expected 1 violation for invalid pattern, got %d", len(results))
	}
}

func TestValidateSchema_OptionalMissingKeySkipped(t *testing.T) {
	env := map[string]string{}
	fields := []SchemaField{{Key: "OPTIONAL_KEY", Required: false, Pattern: `^\d+$`}}
	results := ValidateSchema(env, fields)
	if len(results) != 0 {
		t.Fatalf("expected no violations for missing optional key, got %d", len(results))
	}
}
