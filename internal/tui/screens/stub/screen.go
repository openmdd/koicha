package stub

import (
	tea "charm.land/bubbletea/v2"

	"github.com/openmdd/koicha/internal/tui/nav"
	"github.com/openmdd/koicha/internal/tui/ui/components/layout"
	"github.com/openmdd/koicha/internal/tui/ui/components/statusbar"
	"github.com/openmdd/koicha/internal/tui/ui/theme"
)

type Params struct {
	Title   string
	Message string
}

type Screen struct {
	layout    layout.Model
	statusbar statusbar.Model
	styles    theme.Styles
	title     string
	message   string
}

func New(styles theme.Styles, params Params) Screen {
	return Screen{
		layout:    layout.New(styles),
		statusbar: statusbar.New(styles),
		styles:    styles,
		title:     params.Title,
		message:   params.Message,
	}
}

func (s Screen) ID() nav.ScreenID { return nav.ScreenStub }

func (s Screen) Init() tea.Cmd { return nil }

func (s Screen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.layout.SetSize(msg.Width, msg.Height)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "left", "esc":
			return s, func() tea.Msg { return nav.Back() }
		}
	}
	return s, nil
}

func (s Screen) View() string {
	help := s.statusbar.HelpBackOnly()
	return s.layout.Render("", s.layout.Panel(layout.PanelOptions{
		Title:     s.title,
		Body:      s.styles.Subtle.Render(s.message),
		BodyAlign: layout.AlignCenter,
		Active:    true,
	}), help)
}
