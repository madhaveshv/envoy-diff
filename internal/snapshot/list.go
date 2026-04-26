package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Entry is a lightweight summary of a snapshot file on disk.
type Entry struct {
	Label     string `json:"label"`
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
}

// List returns all snapshot entries found in dir, sorted by filename (chronological).
func List(dir string) ([]Entry, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("snapshot: glob: %w", err)
	}
	sort.Strings(matches)
	var entries []Entry
	for _, p := range matches {
		entry, err := summaryFromFile(p)
		if err != nil {
			continue // skip malformed files
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func summaryFromFile(path string) (Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return Entry{}, err
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return Entry{}, err
	}
	return Entry{
		Label:     s.Label,
		Timestamp: s.Timestamp.Format("2006-01-02T15:04:05Z"),
		Path:      path,
	}, nil
}

// FormatList returns a human-readable listing of snapshot entries.
func FormatList(entries []Entry) string {
	if len(entries) == 0 {
		return "No snapshots found.\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-20s  %-24s  %s\n", "LABEL", "TIMESTAMP", "PATH"))
	sb.WriteString(strings.Repeat("-", 72) + "\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%-20s  %-24s  %s\n", e.Label, e.Timestamp, e.Path))
	}
	return sb.String()
}
