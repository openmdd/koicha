package nav

import (
	tea "charm.land/bubbletea/v2"

	"github.com/openmdd/koicha/internal/bento"
)

type NavigateMsg struct {
	To     ScreenID
	Params any
}

type BackMsg struct{}

type SelectBentoMsg struct {
	Bento bento.Bento
}

func Navigate(to ScreenID, params any) tea.Msg {
	return NavigateMsg{
		To:     to,
		Params: params,
	}
}

func Back() tea.Msg {
	return BackMsg{}
}

func SelectBento(b bento.Bento) tea.Msg {
	return SelectBentoMsg{Bento: b}
}

func ParamsAs[T any](value any) (T, bool) {
	typed, ok := value.(T)
	return typed, ok
}
