package home

import (
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/openmdd/koicha/internal/bento"
	"github.com/openmdd/koicha/internal/tui/nav"
	"github.com/openmdd/koicha/internal/tui/screens/stub"
	"github.com/openmdd/koicha/internal/tui/ui/components/layout"
	"github.com/openmdd/koicha/internal/tui/ui/components/list"
	"github.com/openmdd/koicha/internal/tui/ui/components/metarows"
	"github.com/openmdd/koicha/internal/tui/ui/components/notice"
	"github.com/openmdd/koicha/internal/tui/ui/components/statusbar"
	"github.com/openmdd/koicha/internal/tui/ui/theme"
)

type menuAction struct {
	screen nav.ScreenID
	params any
}

type Screen struct {
	layout    layout.Model
	statusbar statusbar.Model
	styles    theme.Styles
	version   string
	startup   time.Duration
	current   *bento.Bento

	menu   list.Model[menuAction]
	notice notice.Model
}

func New(styles theme.Styles, version string, startup time.Duration, current *bento.Bento) Screen {
	return Screen{
		layout:    layout.New(styles),
		statusbar: statusbar.New(styles),
		styles:    styles,
		version:   version,
		startup:   startup,
		current:   current,
		menu:      list.New(styles, menuItems()),
		notice:    notice.New(styles),
	}
}

func (s Screen) ID() nav.ScreenID { return nav.ScreenHome }

func (s Screen) Init() tea.Cmd { return nil }

func (s Screen) Refresh() tea.Cmd { return nil }

func (s Screen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.layout.SetSize(msg.Width, msg.Height)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			s.menu.MoveUp()
		case "down":
			s.menu.MoveDown()
		case "right", "enter":
			selected, ok := s.menu.Selected()
			if !ok {
				return s, nil
			}
			return s, func() tea.Msg {
				return nav.Navigate(selected.Value.screen, selected.Value.params)
			}
		}
	}
	return s, nil
}

func (s Screen) View() string {
	body := s.layout.Stack(
		s.renderBrandPanel(),
		s.layout.Panel(layout.PanelOptions{
			Title:     "Actions",
			Subtitle:  "Inspect Kafka state, run operational actions, or manage Bento connections.",
			Body:      s.menu.View(),
			BodyAlign: layout.AlignCenter,
			Active:    true,
		}),
	)

	help := s.statusbar.HelpListOpen()
	return s.layout.RenderBare("", body, help)
}

func (s Screen) renderBrandPanel() string {
	logo := koichaLogo
	if s.layout.IsNarrow(80) {
		logo = koichaCompactLogo
	}

	parts := []string{
		s.styles.Logo.Render(logo),
		"",
		s.renderMetadataBlock(),
	}
	if s.current == nil {
		parts = append(parts, "", s.notice.Render(
			"No Bento selected yet",
			"Open Config to create or inspect Bento profiles.",
		))
	}

	return s.layout.Panel(layout.PanelOptions{
		Body: strings.Join(parts, "\n"),
	})
}

func (s Screen) metadataRows() []metarows.Row {
	rows := []metarows.Row{
		{Label: "version", Value: s.version},
		{Label: "startup", Value: formatStartupDuration(s.startup)},
	}

	if s.current != nil {
		rows = append(rows, metarows.Row{Label: "current bento", Value: s.current.Metadata.Name})
		if s.current.Metadata.DisplayName != "" {
			rows = append(rows, metarows.Row{Label: "display name", Value: s.current.Metadata.DisplayName})
		}
		if s.current.Spec.Profile.Environment != "" {
			rows = append(rows, metarows.Row{Label: "profile", Value: s.current.Spec.Profile.Environment})
		}
		return rows
	}

	return append(rows, metarows.Row{Label: "current bento", Value: "not selected"})
}

func (s Screen) renderMetadataBlock() string {
	return metarows.Render(s.styles, s.metadataRows())
}

func formatStartupDuration(duration time.Duration) string {
	ms := float64(duration) / float64(time.Millisecond)
	value := strconv.FormatFloat(ms, 'f', 2, 64)
	return value + "ms"
}

func menuItems() []list.Item[menuAction] {
	return []list.Item[menuAction]{
		{
			Title: "inspect",
			Value: menuAction{
				screen: nav.ScreenStub,
				params: stub.Params{
					Title:   "Inspect",
					Message: "WIP",
				},
			},
		},
		{
			Title: "operate",
			Value: menuAction{
				screen: nav.ScreenStub,
				params: stub.Params{
					Title:   "Operate",
					Message: "WIP",
				},
			},
		},
		{
			Title: "config",
			Value: menuAction{screen: nav.ScreenConfig},
		},
	}
}

const koichaLogo = `
██╗  ██╗ ██████╗ ██╗ ██████╗██╗  ██╗ █████╗
██║ ██╔╝██╔═══██╗██║██╔════╝██║  ██║██╔══██╗
█████╔╝ ██║   ██║██║██║     ███████║███████║
██╔═██╗ ██║   ██║██║██║     ██╔══██║██╔══██║
██║  ██╗╚██████╔╝██║╚██████╗██║  ██║██║  ██║
╚═╝  ╚═╝ ╚═════╝ ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝
`

const koichaCompactLogo = `
koicha
`
