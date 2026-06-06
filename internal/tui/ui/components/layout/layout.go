package layout

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/openmdd/koicha/internal/tui/ui/theme"
)

const (
	appFrameMargin   = 10
	frameInnerMargin = 6
	panelOuterMargin = 4
	panelMinWidth    = 56
	panelMaxWidth    = 104
	panelInnerMargin = 8
	minInnerWidth    = 20
)

type Align int

const (
	AlignCenter Align = iota
	AlignLeft
	AlignRight
)

// Model is a shared frame for screen-like content.
type Model struct {
	styles theme.Styles
	Width  int
	Height int
}

// PanelOptions describes one reusable bordered content block.
type PanelOptions struct {
	Title     string
	Subtitle  string
	Body      string
	BodyAlign Align
	Active    bool
}

// New creates a reusable layout model.
func New(styles theme.Styles) Model {
	return Model{
		styles: styles,
		Width:  theme.DefaultWidth,
		Height: theme.DefaultHeight,
	}
}

// SetSize updates current viewport size.
func (m *Model) SetSize(width, height int) {
	if width > 0 {
		m.Width = width
	}
	if height > 0 {
		m.Height = height
	}
}

// PanelWidth returns a responsive width for screen panels.
func (m *Model) PanelWidth() int {
	return clamp(m.FrameInnerWidth()-panelOuterMargin, panelMinWidth, panelMaxWidth)
}

// InnerWidth returns the content width inside a screen panel.
func (m *Model) InnerWidth() int {
	return max(m.PanelWidth()-panelInnerMargin, minInnerWidth)
}

// FrameWidth returns the width of the shared application frame.
func (m *Model) FrameWidth() int {
	margin := max(m.Width/appFrameMargin, 0)
	return clamp(m.Width-margin, 1, m.Width)
}

// FrameHeight returns the height of the shared application frame.
func (m *Model) FrameHeight() int {
	margin := max(m.Height/appFrameMargin, 0)
	return clamp(m.Height-margin, 1, m.Height)
}

// FrameInnerWidth returns the content width inside the shared application frame.
func (m *Model) FrameInnerWidth() int {
	return max(m.FrameWidth()-frameInnerMargin, minInnerWidth)
}

// FrameInnerHeight returns the content height inside the shared application frame.
func (m *Model) FrameInnerHeight() int {
	return max(m.FrameHeight()-frameInnerMargin, 1)
}

// IsNarrow reports whether the current viewport is narrower than width.
func (m *Model) IsNarrow(width int) bool {
	return m.Width < width
}

// Center renders content centered inside the current inner panel width.
func (m *Model) Center(content string) string {
	return m.Align(content, AlignCenter)
}

// Align renders content inside the current inner panel width.
func (m *Model) Align(content string, align Align) string {
	return lipgloss.NewStyle().
		Width(m.InnerWidth()).
		Align(lipglossPosition(align)).
		Render(content)
}

// Panel renders a reusable bordered panel with optional title and subtitle.
func (m *Model) Panel(options PanelOptions) string {
	sections := make([]string, 0, 2)
	if options.Subtitle != "" {
		sections = append(sections, m.Center(m.styles.Subtle.Render(options.Subtitle)))
	}
	if options.Body != "" {
		sections = append(sections, m.Align(options.Body, options.BodyAlign))
	}
	return m.renderPanelBorder(options.Title, strings.Join(sections, "\n\n"), options.Active)
}

// Stack joins screen sections vertically with one empty line between them.
func (m *Model) Stack(sections ...string) string {
	visible := make([]string, 0, len(sections))
	for _, section := range sections {
		if section != "" {
			if len(visible) > 0 {
				visible = append(visible, "")
			}
			visible = append(visible, section)
		}
	}
	return lipgloss.JoinVertical(lipgloss.Center, visible...)
}

// Render composes header, body and footer sections.
func (m *Model) Render(header, body, footer string) string {
	content := m.renderContent(header, body, footer)
	framedContent := lipgloss.Place(
		m.FrameInnerWidth(),
		m.FrameInnerHeight(),
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
	framed := m.styles.AppFrame.
		Width(m.FrameWidth()).
		Height(m.FrameHeight()).
		Render(framedContent)
	centered := lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		framed,
	)
	return m.styles.App.Render(centered)
}

// RenderBare composes a screen without the shared application frame.
func (m *Model) RenderBare(header, body, footer string) string {
	content := m.renderContent(header, body, footer)
	centered := lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
	return m.styles.App.Render(centered)
}

func (m *Model) renderContent(header, body, footer string) string {
	parts := make([]string, 0, 3)
	if header != "" {
		parts = append(parts, m.styles.Header.Render(header))
	}
	if body != "" {
		parts = append(parts, body)
	}
	if footer != "" {
		parts = append(parts, footer)
	}

	return strings.Join(parts, "\n")
}

// Size returns current viewport size for responsive rendering.
func (m *Model) Size() (int, int) {
	return m.Width, m.Height
}

func clamp(value, minValue, maxValue int) int {
	if maxValue < minValue {
		return maxValue
	}
	return min(max(value, minValue), maxValue)
}

func (m *Model) panelTitle(title string, active bool) string {
	style := m.styles.PanelTitle.Inherit(m.styles.PanelBorder)
	if active {
		style = m.styles.PanelTitle.Inherit(m.styles.ActivePanel)
	}
	return style.Render(" " + title + " ")
}

func (m *Model) renderPanelBorder(title, content string, active bool) string {
	width := m.PanelWidth()
	borderStyle := m.styles.PanelBorder
	if active {
		borderStyle = m.styles.ActivePanel
	}

	lines := []string{m.panelTopBorder(title, width, active)}
	lines = append(lines, m.panelEmptyLine(width, borderStyle))
	for _, line := range splitLines(content) {
		lines = append(lines, m.panelContentLine(line, width, borderStyle))
	}
	lines = append(lines, m.panelEmptyLine(width, borderStyle))
	lines = append(lines, m.panelBottomBorder(width, borderStyle))
	return strings.Join(lines, "\n")
}

func (m *Model) panelTopBorder(title string, width int, active bool) string {
	borderStyle := m.styles.PanelBorder
	if active {
		borderStyle = m.styles.ActivePanel
	}
	if title == "" {
		return borderStyle.Render("╭" + strings.Repeat("─", max(width-2, 0)) + "╮")
	}

	label := m.panelTitle(title, active)
	labelWidth := lipgloss.Width(label)
	fillWidth := max(width-2-labelWidth, 0)
	return borderStyle.Render("╭") +
		label +
		borderStyle.Render(strings.Repeat("─", fillWidth)+"╮")
}

func (m *Model) panelBottomBorder(width int, borderStyle lipgloss.Style) string {
	return borderStyle.Render("╰" + strings.Repeat("─", max(width-2, 0)) + "╯")
}

func (m *Model) panelEmptyLine(width int, borderStyle lipgloss.Style) string {
	return m.panelContentLine("", width, borderStyle)
}

func (m *Model) panelContentLine(line string, width int, borderStyle lipgloss.Style) string {
	const horizontalPadding = 3

	innerWidth := max(width-2-horizontalPadding*2, 0)
	visibleWidth := lipgloss.Width(line)
	padding := strings.Repeat(" ", max(innerWidth-visibleWidth, 0))
	return borderStyle.Render("│") +
		strings.Repeat(" ", horizontalPadding) +
		line +
		padding +
		strings.Repeat(" ", horizontalPadding) +
		borderStyle.Render("│")
}

func splitLines(content string) []string {
	if content == "" {
		return nil
	}
	return strings.Split(content, "\n")
}

func lipglossPosition(align Align) lipgloss.Position {
	switch align {
	case AlignLeft:
		return lipgloss.Left
	case AlignRight:
		return lipgloss.Right
	default:
		return lipgloss.Center
	}
}
