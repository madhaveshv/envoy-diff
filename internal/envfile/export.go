package envfile

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// ExportFormat defines the output format for exported env vars.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv"
	FormatShell  ExportFormat = "shell"
	FormatJSON   ExportFormat = "json"
)

// ExportOptions controls how env vars are written.
type ExportOptions struct {
	Format ExportFormat
	Sorted bool
}

// Export writes the given env map to a file at the given path.
func Export(env map[string]string, path string, opts ExportOptions) error {
	if path == "" {
		return fmt.Errorf("export path must not be empty")
	}

	content, err := Render(env, opts)
	if err != nil {
		return fmt.Errorf("render failed: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	return nil
}

// Render returns the string representation of the env map in the given format.
func Render(env map[string]string, opts ExportOptions) (string, error) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if opts.Sorted {
		sort.Strings(keys)
	}

	switch opts.Format {
	case FormatShell:
		return renderShell(keys, env), nil
	case FormatJSON:
		return renderJSON(keys, env), nil
	case FormatDotenv, "":
		return renderDotenv(keys, env), nil
	default:
		return "", fmt.Errorf("unsupported export format: %q", opts.Format)
	}
}

func renderDotenv(keys []string, env map[string]string) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%q\n", k, env[k])
	}
	return sb.String()
}

func renderShell(keys []string, env map[string]string) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "export %s=%q\n", k, env[k])
	}
	return sb.String()
}

func renderJSON(keys []string, env map[string]string) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", k, env[k], comma)
	}
	sb.WriteString("}\n")
	return sb.String()
}
