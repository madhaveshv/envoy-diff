package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// FormatText writes a human-readable diff summary to w.
func FormatText(w io.Writer, result *Result) {
	changes := make([]Change, len(result.Changes))
	copy(changes, result.Changes)
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	if len(changes) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	fmt.Fprintf(w, "%-12s %-30s %-20s %s\n", "STATUS", "KEY", "OLD VALUE", "NEW VALUE")
	fmt.Fprintln(w, strings.Repeat("-", 80))

	for _, c := range changes {
		var status, oldVal, newVal string
		switch c.Type {
		case Added:
			status = "[+] added"
			oldVal = "-"
			newVal = c.NewValue
		case Removed:
			status = "[-] removed"
			oldVal = c.OldValue
			newVal = "-"
		case Modified:
			status = "[~] modified"
			oldVal = c.OldValue
			newVal = c.NewValue
		}
		fmt.Fprintf(w, "%-12s %-30s %-20s %s\n", status, c.Key, truncate(oldVal, 20), truncate(newVal, 40))
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
