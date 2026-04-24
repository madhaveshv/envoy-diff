package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Env represents a parsed set of environment variables as a key-value map.
type Env map[string]string

// ParseFile reads a .env file from the given path and returns an Env map.
// Lines starting with '#' are treated as comments and ignored.
// Empty lines are also ignored.
// Each valid line must be in KEY=VALUE format.
func ParseFile(path string) (Env, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envfile: failed to open %q: %w", path, err)
	}
	defer f.Close()

	env := make(Env)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			return nil, fmt.Errorf("envfile: invalid syntax at line %d: %q", lineNum, line)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		// Strip optional surrounding quotes from value
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		if key == "" {
			return nil, fmt.Errorf("envfile: empty key at line %d", lineNum)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envfile: error reading %q: %w", path, err)
	}

	return env, nil
}
