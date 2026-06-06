package bentos

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/openmdd/koicha/internal/bento"
	"github.com/openmdd/koicha/internal/tui/nav"
	"github.com/openmdd/koicha/internal/tui/ui/components/layout"
	"github.com/openmdd/koicha/internal/tui/ui/components/list"
	"github.com/openmdd/koicha/internal/tui/ui/components/statusbar"
	"github.com/openmdd/koicha/internal/tui/ui/theme"
)

type bentoListItem = list.Item[bento.Bento]

type Screen struct {
	layout    layout.Model
	statusbar statusbar.Model
	styles    theme.Styles
	store     bento.Store

	bentos  list.Model[bento.Bento]
	loading bool
	loadErr error
}

type bentosLoadedMsg struct {
	bentos []bento.Bento
	err    error
}

func New(styles theme.Styles, store bento.Store) Screen {
	return Screen{
		layout:    layout.New(styles),
		statusbar: statusbar.New(styles),
		styles:    styles,
		store:     store,
		bentos:    list.New[bento.Bento](styles, nil),
		loading:   true,
	}
}

func (s Screen) ID() nav.ScreenID { return nav.ScreenBentos }

func (s Screen) Init() tea.Cmd {
	return s.loadBentosCmd()
}

func (s Screen) Refresh() tea.Cmd { return s.loadBentosCmd() }

func (s Screen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.layout.SetSize(msg.Width, msg.Height)
	case bentosLoadedMsg:
		s.loading = false
		if msg.err != nil {
			s.loadErr = msg.err
			return s, nil
		}
		s.loadErr = nil
		s.bentos.SetItems(bentoItems(msg.bentos))
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			s.bentos.MoveUp()
		case "down":
			s.bentos.MoveDown()
		case "left", "esc":
			return s, func() tea.Msg { return nav.Back() }
		case "right", "enter":
			selected, ok := s.bentos.Selected()
			if !ok {
				return s, nil
			}
			return s, func() tea.Msg {
				return nav.SelectBento(selected.Value)
			}
		}
	}
	return s, nil
}

func (s Screen) View() string {
	body := s.layout.Panel(layout.PanelOptions{
		Title:     "Bento profiles",
		Subtitle:  "Select a Bento profile to inspect.",
		Body:      s.renderBentoList(),
		BodyAlign: layout.AlignLeft,
		Active:    true,
	})

	help := s.statusbar.HelpListSelectBack()
	return s.layout.Render("", body, help)
}

func (s Screen) renderBentoList() string {
	if s.loading {
		return s.layout.Center(s.styles.Subtle.Render("Loading Bento profiles..."))
	}
	if s.loadErr != nil {
		return strings.Join([]string{
			s.styles.Danger.Render("Could not load Bento profiles."),
			s.styles.Subtle.Render(s.loadErr.Error()),
		}, "\n")
	}

	return s.bentos.View()
}

func (s Screen) loadBentosCmd() tea.Cmd {
	bentoStore := s.store

	return func() tea.Msg {
		bentos, err := bentoStore.List()
		return bentosLoadedMsg{
			bentos: bentos,
			err:    err,
		}
	}
}

func bentoItems(bentos []bento.Bento) []bentoListItem {
	items := make([]bentoListItem, 0, len(bentos))
	for _, item := range bentos {
		items = append(items, bentoListItem{
			Title:  item.Metadata.Name,
			Detail: profileLabel(item.Spec.Profile),
			Value:  item,
		})
	}
	return items
}

func profileLabel(profile bento.Profile) string {
	return profile.Environment
}
