package metarows

import (
	"strings"

	"github.com/openmdd/koicha/internal/tui/ui/theme"
)

// Row is a key/value metadata entry rendered as one line.
type Row struct {
	Label string
	Value string
}

// Render prints rows with aligned labels and shared styling.
func Render(styles theme.Styles, rows []Row) string {
	if len(rows) == 0 {
		return ""
	}

	labelWidth := 0
	valueWidth := 0
	for _, row := range rows {
		labelWidth = max(labelWidth, runeWidth(row.Label))
		valueWidth = max(valueWidth, runeWidth(row.Value))
	}

	lines := make([]string, 0, len(rows))
	for _, row := range rows {
		labelPadding := strings.Repeat(" ", max(labelWidth-runeWidth(row.Label), 0))
		valuePadding := strings.Repeat(" ", max(valueWidth-runeWidth(row.Value), 0))
		lines = append(lines, styles.Subtle.Render(row.Label+labelPadding+": ")+styles.Error.Render(row.Value+valuePadding))
	}
	return strings.Join(lines, "\n")
}

func runeWidth(value string) int {
	return len([]rune(value))
}
