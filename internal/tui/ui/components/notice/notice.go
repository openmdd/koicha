package notice

import (
	"strings"

	"github.com/openmdd/koicha/internal/tui/ui/theme"
)

// Model renders short contextual guidance without mixing it with metadata or
// keybinding help.
type Model struct {
	styles theme.Styles
}

func New(styles theme.Styles) Model {
	return Model{styles: styles}
}

func (m Model) Render(title, body string) string {
	parts := make([]string, 0, 2)
	if title != "" {
		parts = append(parts, m.styles.NoticeTitle.Render(title))
	}
	if body != "" {
		parts = append(parts, m.styles.NoticeBody.Render(body))
	}
	return strings.Join(parts, "\n")
}
