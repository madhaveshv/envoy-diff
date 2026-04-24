package envfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SupportedExtensions lists valid env file extensions.
var SupportedExtensions = []string{".env", ".txt", ""}

// LoadPair loads two env files and returns their parsed key-value maps.
// It validates that both paths exist and have supported extensions.
func LoadPair(fromPath, toPath string) (map[string]string, map[string]string, error) {
	fromEnv, err := Load(fromPath)
	if err != nil {
		return nil, nil, fmt.Errorf("loading 'from' file %q: %w", fromPath, err)
	}

	toEnv, err := Load(toPath)
	if err != nil {
		return nil, nil, fmt.Errorf("loading 'to' file %q: %w", toPath, err)
	}

	return fromEnv, toEnv, nil
}

// Load reads and parses a single env file after validating it.
func Load(path string) (map[string]string, error) {
	if err := validatePath(path); err != nil {
		return nil, err
	}
	return ParseFile(path)
}

// validatePath checks that the file exists, is readable, and has a supported extension.
func validatePath(path string) error {
	ext := strings.ToLower(filepath.Ext(path))

	supported := false
	for _, e := range SupportedExtensions {
		if ext == e {
			supported = true
			break
		}
	}
	if !supported {
		return fmt.Errorf("unsupported file extension %q (supported: .env, .txt, or no extension)", ext)
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %q", path)
		}
		return fmt.Errorf("cannot access file %q: %w", path, err)
	}

	if info.IsDir() {
		return fmt.Errorf("path %q is a directory, not a file", path)
	}

	return nil
}
