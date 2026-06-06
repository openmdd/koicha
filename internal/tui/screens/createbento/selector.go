package createbento

import (
	"strings"

	"github.com/openmdd/koicha/internal/tui/ui/theme"
)

func moveOption(options []selectOption, current int, delta int) int {
	if len(options) == 0 {
		return 0
	}

	next := current
	for {
		next += delta
		if next < 0 {
			next = len(options) - 1
		}
		if next >= len(options) {
			next = 0
		}
		if !options[next].Disabled || next == current {
			return next
		}
	}
}

func selectedTitle(options []selectOption, index int) string {
	if index < 0 || index >= len(options) {
		return ""
	}
	return options[index].Title
}

func renderSelector(styles theme.Styles, options []selectOption, selected int, focused bool) string {
	rows := make([]string, 0, len(options))
	for i, option := range options {
		prefix := "  "
		style := styles.Subtle
		if option.Disabled {
			style = styles.Disabled
		}
		if i == selected {
			prefix = "> "
			if focused && !option.Disabled {
				style = styles.Title
			}
		}

		row := prefix + option.Title
		if option.Detail != "" {
			row += " | " + option.Detail
		}
		rows = append(rows, style.Render(row))
	}
	return strings.Join(rows, "\n")
}

func renderInlineChoice(styles theme.Styles, options []selectOption, selected int, focused bool) string {
	parts := make([]string, 0, len(options))
	for i, option := range options {
		label := option.Title
		if option.Detail != "" {
			label += " " + option.Detail
		}
		if i == selected {
			label = "[" + label + "]"
		}

		switch {
		case option.Disabled:
			parts = append(parts, styles.Disabled.Render(label))
		case i == selected && focused:
			parts = append(parts, styles.Title.Render(label))
		case i == selected:
			parts = append(parts, styles.Error.Render(label))
		default:
			parts = append(parts, styles.Subtle.Render(label))
		}
	}
	return strings.Join(parts, "  ")
}
