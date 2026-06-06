package list

import (
	"fmt"
	"strings"

	"github.com/openmdd/koicha/internal/tui/ui/theme"
)

// Item represents one list row with typed payload.
type Item[T any] struct {
	Title    string
	Subtitle string
	Detail   string
	Badge    string
	Disabled bool
	Value    T
}

// Model is a generic reusable selectable list.
type Model[T any] struct {
	styles theme.Styles
	items  []Item[T]
	cursor int
}

func New[T any](styles theme.Styles, items []Item[T]) Model[T] {
	return Model[T]{
		styles: styles,
		items:  items,
	}
}

func (m *Model[T]) SetItems(items []Item[T]) {
	m.items = items
	if len(m.items) == 0 {
		m.cursor = 0
		return
	}
	if m.cursor > len(m.items)-1 {
		m.cursor = len(m.items) - 1
	}
	m.ensureSelectable(1)
}

func (m *Model[T]) MoveUp() {
	m.move(-1)
}

func (m *Model[T]) MoveDown() {
	m.move(1)
}

func (m *Model[T]) Selected() (Item[T], bool) {
	if len(m.items) == 0 {
		var zero Item[T]
		return zero, false
	}
	if m.items[m.cursor].Disabled {
		var zero Item[T]
		return zero, false
	}
	return m.items[m.cursor], true
}

func (m *Model[T]) Len() int {
	return len(m.items)
}

func (m *Model[T]) View() string {
	if len(m.items) == 0 {
		return m.styles.Subtle.Render("No entries yet")
	}

	var b strings.Builder
	titleWidth := m.titleWidth()
	for i, item := range m.items {
		prefix := "  "
		lineStyle := m.styles.Subtle
		if item.Disabled {
			lineStyle = m.styles.Disabled
		}
		if i == m.cursor && !item.Disabled {
			prefix = "> "
			lineStyle = m.styles.Title
		} else if i == m.cursor {
			prefix = "> "
		}
		padding := strings.Repeat(" ", titleWidth-displayWidth(item.Title))
		b.WriteString(lineStyle.Render(fmt.Sprintf("%s%s%s", prefix, item.Title, padding)))
		if item.Detail != "" {
			b.WriteString(lineStyle.Render(" | "))
			b.WriteString(m.styles.Error.Render(item.Detail))
		}
		if item.Badge != "" {
			b.WriteString(lineStyle.Render(" | "))
			b.WriteString(m.styles.Disabled.Render(item.Badge))
		}
		if item.Subtitle != "" {
			b.WriteString("\n")
			b.WriteString(m.styles.Subtle.Render("   " + item.Subtitle))
		}
		if i < len(m.items)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (m *Model[T]) move(delta int) {
	if len(m.items) == 0 {
		return
	}

	next := m.cursor
	for {
		next += delta
		if next < 0 || next >= len(m.items) {
			return
		}
		m.cursor = next
		if !m.items[m.cursor].Disabled {
			return
		}
	}
}

func (m *Model[T]) ensureSelectable(delta int) {
	if len(m.items) == 0 || !m.items[m.cursor].Disabled {
		return
	}
	m.move(delta)
	if m.items[m.cursor].Disabled {
		m.move(-delta)
	}
}

func (m *Model[T]) titleWidth() int {
	width := 0
	for _, item := range m.items {
		width = max(width, displayWidth(item.Title))
	}
	return width
}

func displayWidth(value string) int {
	return len([]rune(value))
}
