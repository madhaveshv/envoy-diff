package envfile

import (
	"testing"
)

func TestConvert_NoOptions(t *testing.T) {
	input := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Convert(input, ConvertOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("expected unchanged map, got %v", out)
	}
}

func TestConvert_UppercaseKeys(t *testing.T) {
	input := map[string]string{"foo": "1", "bar": "2"}
	out, err := Convert(input, ConvertOptions{KeyCase: "upper"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["FOO"]; !ok {
		t.Errorf("expected key FOO, got %v", out)
	}
	if _, ok := out["BAR"]; !ok {
		t.Errorf("expected key BAR, got %v", out)
	}
}

func TestConvert_LowercaseKeys(t *testing.T) {
	input := map[string]string{"FOO": "1", "BAR": "2"}
	out, err := Convert(input, ConvertOptions{KeyCase: "lower"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["foo"] != "1" || out["bar"] != "2" {
		t.Errorf("expected lowercase keys, got %v", out)
	}
}

func TestConvert_AddKeyPrefix(t *testing.T) {
	input := map[string]string{"HOST": "localhost"}
	out, err := Convert(input, ConvertOptions{KeyPrefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST, got %v", out)
	}
}

func TestConvert_StripPrefix(t *testing.T) {
	input := map[string]string{"PROD_HOST": "example.com", "PROD_PORT": "443"}
	out, err := Convert(input, ConvertOptions{StripPrefix: "PROD_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "example.com" || out["PORT"] != "443" {
		t.Errorf("expected stripped keys, got %v", out)
	}
}

func TestConvert_DoubleQuoteValues(t *testing.T) {
	input := map[string]string{"KEY": "value"}
	out, err := Convert(input, ConvertOptions{ValueQuoteStyle: "double"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != `"value"` {
		t.Errorf("expected double-quoted value, got %q", out["KEY"])
	}
}

func TestConvert_SingleQuoteValues(t *testing.T) {
	input := map[string]string{"KEY": "value"}
	out, err := Convert(input, ConvertOptions{ValueQuoteStyle: "single"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != `'value'` {
		t.Errorf("expected single-quoted value, got %q", out["KEY"])
	}
}

func TestConvert_InvalidKeyCase(t *testing.T) {
	_, err := Convert(map[string]string{}, ConvertOptions{KeyCase: "title"})
	if err == nil {
		t.Error("expected error for invalid key_case")
	}
}

func TestConvert_InvalidQuoteStyle(t *testing.T) {
	_, err := Convert(map[string]string{}, ConvertOptions{ValueQuoteStyle: "backtick"})
	if err == nil {
		t.Error("expected error for invalid value_quote_style")
	}
}

func TestConvert_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"foo": "bar"}
	_, err := Convert(input, ConvertOptions{KeyCase: "upper", KeyPrefix: "X_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := input["foo"]; !ok {
		t.Error("Convert mutated the input map")
	}
}
