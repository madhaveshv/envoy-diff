package envfile

import (
	"testing"
)

var sampleEnv = map[string]string{
	"APP_HOST":     "localhost",
	"APP_PORT":     "8080",
	"DB_HOST":      "db.internal",
	"DB_PASSWORD":  "secret",
	"LOG_LEVEL":    "info",
	"SECRET_TOKEN": "tok123",
}

func TestFilter_NoOptions(t *testing.T) {
	result := Filter(sampleEnv, FilterOptions{})
	if len(result) != len(sampleEnv) {
		t.Errorf("expected %d keys, got %d", len(sampleEnv), len(result))
	}
}

func TestFilter_ByPrefix(t *testing.T) {
	result := Filter(sampleEnv, FilterOptions{Prefix: "APP_"})
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
	if _, ok := result["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in result")
	}
}

func TestFilter_ByAllowlist(t *testing.T) {
	result := Filter(sampleEnv, FilterOptions{Keys: []string{"LOG_LEVEL", "DB_HOST"}})
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["LOG_LEVEL"]; !ok {
		t.Error("expected LOG_LEVEL in result")
	}
}

func TestFilter_ByExcludeKeys(t *testing.T) {
	result := Filter(sampleEnv, FilterOptions{ExcludeKeys: []string{"DB_PASSWORD", "SECRET_TOKEN"}})
	if _, ok := result["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be excluded")
	}
	if _, ok := result["SECRET_TOKEN"]; ok {
		t.Error("SECRET_TOKEN should be excluded")
	}
	if len(result) != len(sampleEnv)-2 {
		t.Errorf("expected %d keys, got %d", len(sampleEnv)-2, len(result))
	}
}

func TestFilter_PrefixAndExclude(t *testing.T) {
	result := Filter(sampleEnv, FilterOptions{
		Prefix:      "DB_",
		ExcludeKeys: []string{"DB_PASSWORD"},
	})
	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestFilter_EmptyEnv(t *testing.T) {
	result := Filter(map[string]string{}, FilterOptions{Prefix: "APP_"})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d keys", len(result))
	}
}
