# UI components

This directory contains reusable building blocks for koicha TUI components.
It is a UI-kit layer inside the TUI package tree.

## Packages

- `theme` - shared lipgloss style palette.
- `components/layout` - common frame (`header + body + footer/help`).
- `components/statusbar` - one-line status and hotkeys bar.
- `components/list` - lightweight selectable list.
- `components/table` - reusable wrapper around `bubbles/table`.

## Runtime layer

- Runtime navigation and concrete screens live in `internal/tui`.
- `internal/tui` composes these components and keeps app flow logic out of UI-kit.