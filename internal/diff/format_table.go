package diff

import (
	"fmt"
	"strings"
)

const (
	colKeyWidth   = 30
	colValWidth   = 40
	tableHeader   = "%-4s %-30s %-40s %-40s"
	tableRow      = "%-4s %-30s %-40s %-40s"
	tableDivider  = "----+--------------------------------+------------------------------------------+------------------------------------------"
)

// FormatTable renders diff results as an aligned table for terminal output.
func FormatTable(results []Result) string {
	if len(results) == 0 {
		return "No differences found.\n"
	}

	var sb strings.Builder

	header := fmt.Sprintf(tableHeader, "OP", "KEY", "OLD VALUE", "NEW VALUE")
	sb.WriteString(header + "\n")
	sb.WriteString(tableDivider + "\n")

	for _, r := range results {
		op := statusSymbol(r.Status)
		key := padOrTruncate(r.Key, colKeyWidth)
		oldVal := padOrTruncate(truncate(r.OldValue, colValWidth), colValWidth)
		newVal := padOrTruncate(truncate(r.NewValue, colValWidth), colValWidth)
		line := fmt.Sprintf(tableRow, op, key, oldVal, newVal)
		sb.WriteString(line + "\n")
	}

	return sb.String()
}

func statusSymbol(s Status) string {
	switch s {
	case Added:
		return "+"
	case Removed:
		return "-"
	case Modified:
		return "~"
	default:
		return " "
	}
}

func padOrTruncate(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}
