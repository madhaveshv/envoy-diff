package envfile

import (
	"testing"
)

func TestNormalize_NoOptions(t *testing.T) {
	input := map[string]string{"foo": "bar", "baz": "qux"}
	result := Normalize(input, NormalizeOptions{})
	if result["foo"] != "bar" || result["baz"] != "qux" {
		t.Errorf("expected unchanged map, got %v", result)
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	input := map[string]string{"db_host": "localhost", "Api_Key": "secret"}
	result := Normalize(input, NormalizeOptions{UppercaseKeys: true})

	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to exist")
	}
	if _, ok := result["API_KEY"]; !ok {
		t.Error("expected API_KEY to exist")
	}
	if _, ok := result["db_host"]; ok {
		t.Error("expected original lowercase key to be gone")
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	input := map[string]string{"HOST": "  localhost  ", "PORT": "\t8080\n"}
	result := Normalize(input, NormalizeOptions{TrimValues: true})

	if result["HOST"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", result["HOST"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected '8080', got %q", result["PORT"])
	}
}

func TestNormalize_RemoveEmpty(t *testing.T) {
	input := map[string]string{"HOST": "localhost", "EMPTY": "", "BLANK": "  "}
	result := Normalize(input, NormalizeOptions{TrimValues: true, RemoveEmpty: true})

	if _, ok := result["EMPTY"]; ok {
		t.Error("expected EMPTY to be removed")
	}
	if _, ok := result["BLANK"]; ok {
		t.Error("expected BLANK to be removed after trim")
	}
	if result["HOST"] != "localhost" {
		t.Errorf("expected HOST to remain, got %q", result["HOST"])
	}
}

func TestNormalize_KeyPrefix(t *testing.T) {
	input := map[string]string{"HOST": "localhost", "PORT": "5432"}
	result := Normalize(input, NormalizeOptions{KeyPrefix: "APP_"})

	if result["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %v", result)
	}
	if result["APP_PORT"] != "5432" {
		t.Errorf("expected APP_PORT=5432, got %v", result)
	}
}

func TestNormalize_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"key": "value"}
	_ = Normalize(input, NormalizeOptions{UppercaseKeys: true, KeyPrefix: "X_"})

	if _, ok := input["key"]; !ok {
		t.Error("expected original input to remain unchanged")
	}
}

func TestNormalize_UppercaseAndPrefix(t *testing.T) {
	input := map[string]string{"host": "db.internal"}
	result := Normalize(input, NormalizeOptions{UppercaseKeys: true, KeyPrefix: "PROD_"})

	if result["PROD_HOST"] != "db.internal" {
		t.Errorf("expected PROD_HOST=db.internal, got %v", result)
	}
}
