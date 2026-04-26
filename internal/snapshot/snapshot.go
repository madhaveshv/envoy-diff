package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a saved state of environment variables at a point in time.
type Snapshot struct {
	Label     string            `json:"label"`
	Timestamp time.Time         `json:"timestamp"`
	Env       map[string]string `json:"env"`
}

// New creates a new Snapshot with the given label and env map.
func New(label string, env map[string]string) *Snapshot {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return &Snapshot{
		Label:     label,
		Timestamp: time.Now().UTC(),
		Env:       copy,
	}
}

// Save writes the snapshot as a JSON file to the given directory.
func Save(dir string, s *Snapshot) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("snapshot: create dir: %w", err)
	}
	filename := fmt.Sprintf("%s_%s.json", s.Label, s.Timestamp.Format("20060102T150405Z"))
	path := filepath.Join(dir, filename)
	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return "", fmt.Errorf("snapshot: encode: %w", err)
	}
	return path, nil
}

// Load reads a snapshot from the given JSON file path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open: %w", err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}
