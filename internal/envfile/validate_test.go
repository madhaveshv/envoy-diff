package envfile

import (
	"testing"
)

func TestValidate_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	issues := Validate(env, ValidateOptions{})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestValidate_RequireUppercase(t *testing.T) {
	env := map[string]string{
		"app_name": "myapp",
		"PORT":     "8080",
	}
	issues := Validate(env, ValidateOptions{RequireUppercase: true})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "app_name" {
		t.Errorf("expected issue on 'app_name', got %q", issues[0].Key)
	}
	if issues[0].Severity != "warning" {
		t.Errorf("expected severity 'warning', got %q", issues[0].Severity)
	}
}

func TestValidate_DisallowEmpty(t *testing.T) {
	env := map[string]string{
		"PORT":    "",
		"APP":     "myapp",
	}
	issues := Validate(env, ValidateOptions{DisallowEmpty: true})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "PORT" {
		t.Errorf("expected issue on 'PORT', got %q", issues[0].Key)
	}
}

func TestValidate_ForbiddenKeys(t *testing.T) {
	env := map[string]string{
		"DEBUG":   "true",
		"APP":     "myapp",
	}
	issues := Validate(env, ValidateOptions{ForbiddenKeys: []string{"DEBUG"}})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "error" {
		t.Errorf("expected severity 'error', got %q", issues[0].Severity)
	}
}

func TestValidate_KeyPattern(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"port":     "8080",
	}
	issues := Validate(env, ValidateOptions{KeyPattern: `^[A-Z_]+$`})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "port" {
		t.Errorf("expected issue on 'port', got %q", issues[0].Key)
	}
}

func TestValidate_InvalidKeyPattern(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	issues := Validate(env, ValidateOptions{KeyPattern: "[invalid("})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue for bad regex, got %d", len(issues))
	}
	if issues[0].Key != "<options>" {
		t.Errorf("expected issue key '<options>', got %q", issues[0].Key)
	}
}

func TestValidationIssue_String(t *testing.T) {
	issue := ValidationIssue{Key: "FOO", Message: "something wrong", Severity: "error"}
	s := issue.String()
	if s != "[ERROR] FOO: something wrong" {
		t.Errorf("unexpected string: %q", s)
	}
}
