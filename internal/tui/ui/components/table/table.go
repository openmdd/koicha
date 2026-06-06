package table

import (
	"strings"

	btable "charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"

	"github.com/openmdd/koicha/internal/tui/ui/theme"
)

// Model wraps bubbles/table with reusable defaults and helpers.
type Model struct {
	table  btable.Model
	styles theme.Styles
}

func New(styles theme.Styles, columns []btable.Column, rows []btable.Row) Model {
	t := btable.New(
		btable.WithColumns(columns),
		btable.WithRows(rows),
		btable.WithFocused(true),
		btable.WithHeight(12),
	)

	baseStyles := btable.DefaultStyles()
	baseStyles.Header = styles.TableHead
	baseStyles.Selected = styles.TableFocus
	t.SetStyles(baseStyles)

	return Model{
		table:  t,
		styles: styles,
	}
}

func (m *Model) SetRows(rows []btable.Row) {
	m.table.SetRows(rows)
}

func (m *Model) SetSize(width, height int) {
	m.table.SetWidth(width)
	m.table.SetHeight(height)
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return cmd
}

func (m *Model) View() string {
	return m.table.View()
}

func (m *Model) SelectedRow() btable.Row {
	return m.table.SelectedRow()
}

func (m *Model) HelpLine(extra ...string) string {
	parts := []string{"up/down move", "enter select", "esc back", "q quit"}
	if len(extra) > 0 {
		parts = append(parts, extra...)
	}
	return m.styles.Help.Render(strings.Join(parts, " | "))
}
