package nav

import tea "charm.land/bubbletea/v2"

type ScreenID string

const (
	ScreenHome        ScreenID = "home"
	ScreenConfig      ScreenID = "config"
	ScreenBentos      ScreenID = "bentos"
	ScreenBento       ScreenID = "bento"
	ScreenCreateBento ScreenID = "create-bento"
	ScreenStub        ScreenID = "stub"
)

// Screen describes one navigation destination.
type Screen interface {
	ID() ScreenID
	Init() tea.Cmd
	Update(msg tea.Msg) (Screen, tea.Cmd)
	View() string
}

// Refreshable is implemented by screens that support forced reload.
type Refreshable interface {
	Refresh() tea.Cmd
}
