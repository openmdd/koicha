package main

import (
	"log"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/openmdd/koicha/internal/appconfig"
	"github.com/openmdd/koicha/internal/bento"
	"github.com/openmdd/koicha/internal/tui"
)

const version = "0.0.0-pre-alpha"

func main() {
	logger := log.New(os.Stderr, "", 0)

	paths, err := appconfig.Init()
	if err != nil {
		logger.Fatalf("koicha config failed: %v", err)
	}

	currentBento, err := loadCurrentBento(paths)
	if err != nil {
		logger.Printf("koicha current bento ignored: %v", err)
	}

	model := tui.NewModel(tui.Dependencies{
		Version:      version,
		StartedAt:    time.Now(),
		BentoDir:     paths.BentoDir,
		CurrentPath:  paths.CurrentBentoPath,
		CurrentBento: currentBento,
	})

	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		logger.Fatalf("koicha failed: %v", err)
	}
}

func loadCurrentBento(paths appconfig.Paths) (*bento.Bento, error) {
	selectionStore := bento.SelectionStore{Path: paths.CurrentBentoPath}
	name, err := selectionStore.LoadSelectedName()
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, nil
	}

	store := bento.Store{Dir: paths.BentoDir}
	current, err := store.Load(name)
	if err != nil {
		return nil, err
	}
	return &current, nil
}
