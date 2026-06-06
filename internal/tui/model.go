package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/openmdd/koicha/internal/bento"
	"github.com/openmdd/koicha/internal/tui/shell"
)

type Dependencies struct {
	Version      string
	StartedAt    time.Time
	BentoDir     string
	CurrentPath  string
	CurrentBento *bento.Bento
}

func NewModel(deps Dependencies) tea.Model {
	return shell.NewApp(shell.Dependencies{
		Version:      deps.Version,
		StartedAt:    deps.StartedAt,
		BentoDir:     deps.BentoDir,
		CurrentPath:  deps.CurrentPath,
		CurrentBento: deps.CurrentBento,
	})
}
