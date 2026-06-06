package theme

import "charm.land/lipgloss/v2"

const (
	DefaultWidth  = 100
	DefaultHeight = 30
)

// Styles holds all reusable lipgloss styles for screens/components.
type Styles struct {
	App         lipgloss.Style
	AppFrame    lipgloss.Style
	Panel       lipgloss.Style
	AccentPanel lipgloss.Style
	PanelBorder lipgloss.Style
	ActivePanel lipgloss.Style
	PanelTitle  lipgloss.Style
	Header      lipgloss.Style
	Logo        lipgloss.Style
	Title       lipgloss.Style
	Subtle      lipgloss.Style
	Error       lipgloss.Style
	Danger      lipgloss.Style
	Disabled    lipgloss.Style
	NoticeTitle lipgloss.Style
	NoticeBody  lipgloss.Style
	Help        lipgloss.Style
	StatusBar   lipgloss.Style
	TableCell   lipgloss.Style
	TableHead   lipgloss.Style
	TableFocus  lipgloss.Style
}

// NewStyles builds a consistent default style palette.
func NewStyles() Styles {
	const (
		primary   = "#76a32e"
		secondary = "#2e76a3"
		accent    = "#a32e76"
		notice    = "#d6b25e"
		danger    = "#d65e5e"
		disabled  = "#6f7679"
	)

	return Styles{
		App: lipgloss.NewStyle().
			Padding(0, 1),
		AppFrame: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(primary)).
			Padding(1, 2),
		Panel: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(primary)).
			Padding(1, 2),
		AccentPanel: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(accent)).
			Padding(1, 2),
		PanelBorder: lipgloss.NewStyle().
			Foreground(lipgloss.Color(primary)),
		ActivePanel: lipgloss.NewStyle().
			Foreground(lipgloss.Color(accent)),
		PanelTitle: lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true),
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color(primary)).
			Bold(true),
		Logo: lipgloss.NewStyle().
			Foreground(lipgloss.Color(primary)).
			Bold(true),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color(secondary)).
			Bold(true),
		Subtle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#b6c3a1")),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color(accent)).
			Bold(true),
		Danger: lipgloss.NewStyle().
			Foreground(lipgloss.Color(danger)).
			Bold(true),
		Disabled: lipgloss.NewStyle().
			Foreground(lipgloss.Color(disabled)),
		NoticeTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(notice)).
			Bold(true),
		NoticeBody: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c8b98a")),
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8f9599")),
		StatusBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8f9599")).
			Padding(0, 1),
		TableCell: lipgloss.NewStyle().
			Padding(0, 1),
		TableHead: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f4f8ed")).
			Background(lipgloss.Color(secondary)).
			Bold(true),
		TableFocus: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f4f8ed")).
			Background(lipgloss.Color(accent)).
			Bold(true),
	}
}
