package envfile

import (
	"os"
	"testing"
)

func TestInterpolate_NoReferences(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Interpolate(env, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("expected unchanged values, got %v", out)
	}
}

func TestInterpolate_DollarBrace(t *testing.T) {
	env := map[string]string{"BASE": "/app", "LOG_DIR": "${BASE}/logs"}
	out, err := Interpolate(env, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["LOG_DIR"] != "/app/logs" {
		t.Errorf("expected /app/logs, got %q", out["LOG_DIR"])
	}
}

func TestInterpolate_DollarPlain(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "DSN": "postgres://$HOST/db"}
	out, err := Interpolate(env, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "postgres://localhost/db" {
		t.Errorf("expected resolved DSN, got %q", out["DSN"])
	}
}

func TestInterpolate_UseOSEnv(t *testing.T) {
	os.Setenv("_TEST_ENVOY_VAR", "fromOS")
	defer os.Unsetenv("_TEST_ENVOY_VAR")

	env := map[string]string{"VAL": "${_TEST_ENVOY_VAR}"}
	out, err := Interpolate(env, InterpolateOptions{UseOSEnv: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["VAL"] != "fromOS" {
		t.Errorf("expected fromOS, got %q", out["VAL"])
	}
}

func TestInterpolate_FailOnMissing(t *testing.T) {
	env := map[string]string{"VAL": "${UNDEFINED_XYZ}"}
	_, err := Interpolate(env, InterpolateOptions{FailOnMissing: true})
	if err == nil {
		t.Fatal("expected error for undefined variable, got nil")
	}
}

func TestInterpolate_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"A": "1", "B": "${A}-2"}
	original := map[string]string{"A": "1", "B": "${A}-2"}
	Interpolate(env, InterpolateOptions{})
	for k, v := range original {
		if env[k] != v {
			t.Errorf("input mutated at key %q", k)
		}
	}
}
