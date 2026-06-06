package shell

import (
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/openmdd/koicha/internal/bento"
	"github.com/openmdd/koicha/internal/tui/nav"
	"github.com/openmdd/koicha/internal/tui/ui/theme"
)

type screenEntry struct {
	id     nav.ScreenID
	screen nav.Screen
}

type currentBentoSavedMsg struct {
	err error
}

type App struct {
	deps Dependencies

	styles  theme.Styles
	current nav.Screen
	stack   []screenEntry

	currentBento    *bento.Bento
	startupDuration time.Duration

	width  int
	height int
}

func NewApp(deps Dependencies) App {
	app := App{
		deps:            deps,
		styles:          theme.NewStyles(),
		currentBento:    deps.CurrentBento,
		startupDuration: startupDuration(deps.StartedAt),
	}
	app.current = app.newScreen(nav.ScreenHome, nil)
	return app
}

func startupDuration(startedAt time.Time) time.Duration {
	duration := time.Since(startedAt)
	if duration <= 0 {
		return time.Nanosecond
	}
	return duration
}

func (a App) Init() tea.Cmd {
	return a.current.Init()
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
	case nav.NavigateMsg:
		return a.navigate(msg)
	case nav.BackMsg:
		return a.back()
	case nav.SelectBentoMsg:
		return a.selectBento(msg)
	case currentBentoSavedMsg:
		if msg.err != nil {
			// TODO: Surface this error to the user once a global notification mechanism exists.
		}
		return a, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+q", "ctrl+c", "q":
			return a, tea.Quit
		case "ctrl+r":
			if refreshable, ok := a.current.(nav.Refreshable); ok {
				return a, refreshable.Refresh()
			}
			return a, nil
		}
	}

	next, cmd := a.current.Update(msg)
	a.current = next
	return a, cmd
}

func (a App) View() tea.View {
	view := tea.NewView(a.current.View())
	view.AltScreen = true
	return view
}

func (a App) selectBento(msg nav.SelectBentoMsg) (tea.Model, tea.Cmd) {
	a.currentBento = &msg.Bento
	a.stack = nil
	a.current = a.newScreen(nav.ScreenHome, nil)
	a.current = a.resizeScreen(a.current)
	return a, tea.Batch(a.current.Init(), a.saveCurrentBentoCmd(msg.Bento))
}

func (a App) saveCurrentBentoCmd(b bento.Bento) tea.Cmd {
	path := a.deps.CurrentPath
	return func() tea.Msg {
		err := bento.SelectionStore{Path: path}.SaveSelectedName(b.Metadata.Name)
		return currentBentoSavedMsg{err: err}
	}
}

func (a App) navigate(msg nav.NavigateMsg) (tea.Model, tea.Cmd) {
	a.stack = append(a.stack, screenEntry{
		id:     a.current.ID(),
		screen: a.current,
	})
	a.current = a.newScreen(msg.To, msg.Params)
	a.current = a.resizeScreen(a.current)
	return a, a.current.Init()
}

func (a App) back() (tea.Model, tea.Cmd) {
	if len(a.stack) == 0 {
		return a, nil
	}

	last := a.stack[len(a.stack)-1]
	a.stack = a.stack[:len(a.stack)-1]
	a.current = last.screen
	a.current = a.resizeScreen(a.current)
	return a, nil
}

func (a App) resizeScreen(screen nav.Screen) nav.Screen {
	if a.width <= 0 || a.height <= 0 {
		return screen
	}

	resized, _ := screen.Update(tea.WindowSizeMsg{
		Width:  a.width,
		Height: a.height,
	})
	return resized
}
