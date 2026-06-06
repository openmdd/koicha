package config

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/openmdd/koicha/internal/bento"
	"github.com/openmdd/koicha/internal/tui/nav"
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
	menu      list.Model[menuAction]
	notice    notice.Model
	current   *bento.Bento
	bentoDir  string
}

func New(styles theme.Styles, current *bento.Bento, bentoDir string) Screen {
	return Screen{
		layout:    layout.New(styles),
		statusbar: statusbar.New(styles),
		styles:    styles,
		menu:      list.New(styles, menuItems()),
		notice:    notice.New(styles),
		current:   current,
		bentoDir:  bentoDir,
	}
}

func (s Screen) ID() nav.ScreenID { return nav.ScreenConfig }

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
		case "left", "esc":
			return s, func() tea.Msg { return nav.Back() }
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
		s.layout.Panel(layout.PanelOptions{
			Title:     "Bento Config",
			Subtitle:  "Manage Kafka connection profiles for Koicha.",
			Body:      s.renderSummary(),
			BodyAlign: layout.AlignLeft,
		}),
		s.layout.Panel(layout.PanelOptions{
			Title:     "Bento Actions",
			Body:      s.menu.View(),
			BodyAlign: layout.AlignCenter,
			Active:    true,
		}),
	)

	help := s.statusbar.HelpListOpenBack()
	return s.layout.Render("", body, help)
}

func (s Screen) renderSummary() string {
	if s.current != nil {
		rows := []metarows.Row{
			{Label: "current bento", Value: s.current.Metadata.Name},
			{Label: "storage", Value: s.bentoDir},
		}
		if s.current.Spec.Profile.Environment != "" {
			rows = append(rows, metarows.Row{Label: "profile", Value: s.current.Spec.Profile.Environment})
		}
		return metarows.Render(s.styles, rows)
	}

	return strings.Join([]string{
		metarows.Render(s.styles, []metarows.Row{
			{Label: "current bento", Value: "not selected"},
			{Label: "storage", Value: s.bentoDir},
		}),
		"",
		s.notice.Render(
			"No Bento selected yet",
			"Create a Bento or inspect existing profiles to choose one.",
		),
	}, "\n")
}

func menuItems() []list.Item[menuAction] {
	return []list.Item[menuAction]{
		{
			Title: "inspect bento profiles",
			Value: menuAction{screen: nav.ScreenBentos},
		},
		{
			Title: "create bento",
			Value: menuAction{screen: nav.ScreenCreateBento},
		},
	}
}
