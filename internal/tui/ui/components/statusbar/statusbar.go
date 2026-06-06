package statusbar

import (
	"strings"

	"github.com/openmdd/koicha/internal/tui/ui/theme"
)

// Model renders one-line status/help hints.
type Model struct {
	styles theme.Styles
}

func New(styles theme.Styles) Model {
	return Model{styles: styles}
}

func (m Model) Help(actions ...string) string {
	parts := make([]string, 0, len(actions)+2)
	for _, action := range actions {
		if action != "" {
			parts = append(parts, formatAction(action))
		}
	}
	parts = append(parts, "Ctrl+R refresh", "Ctrl+Q quit")
	return m.Render(strings.Join(parts, " | "), "")
}

func (m Model) HelpListOpen() string {
	return m.Help("up/down move", "enter/right open")
}

func (m Model) HelpListOpenBack() string {
	return m.Help("up/down move", "enter/right open", "left/esc back")
}

func (m Model) HelpListSelectBack() string {
	return m.Help("up/down move", "enter/right select", "left/esc back")
}

func (m Model) HelpBackOnly() string {
	return m.Help("left/esc back")
}

func (m Model) Render(left, right string) string {
	if left == "" && right == "" {
		return ""
	}

	switch {
	case left == "":
		return m.styles.StatusBar.Render(right)
	case right == "":
		return m.styles.StatusBar.Render(left)
	default:
		return m.styles.StatusBar.Render(strings.Join([]string{left, right}, " | "))
	}
}

func formatAction(action string) string {
	replacements := []struct {
		old string
		new string
	}{
		{old: "up/down", new: "↑/↓"},
		{old: "shift+tab", new: "Shift+Tab"},
		{old: "tab", new: "Tab"},
		{old: "enter/right", new: "↵/→"},
		{old: "enter", new: "↵"},
		{old: "right", new: "→"},
		{old: "left/esc", new: "←/Esc"},
		{old: "left", new: "←"},
		{old: "esc", new: "Esc"},
		{old: "ctrl+n/ctrl+p", new: "Ctrl+N/Ctrl+P"},
		{old: "ctrl+n", new: "Ctrl+N"},
		{old: "ctrl+p", new: "Ctrl+P"},
	}
	for _, replacement := range replacements {
		if strings.HasPrefix(action, replacement.old) {
			return replacement.new + strings.TrimPrefix(action, replacement.old)
		}
	}
	return action
}
