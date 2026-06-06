package shell

import (
	"github.com/openmdd/koicha/internal/bento"
	"github.com/openmdd/koicha/internal/tui/nav"
	"github.com/openmdd/koicha/internal/tui/screens/bentos"
	"github.com/openmdd/koicha/internal/tui/screens/config"
	"github.com/openmdd/koicha/internal/tui/screens/createbento"
	"github.com/openmdd/koicha/internal/tui/screens/home"
	"github.com/openmdd/koicha/internal/tui/screens/stub"
)

func (a App) newScreen(id nav.ScreenID, params any) nav.Screen {
	store := bento.Store{Dir: a.deps.BentoDir}

	switch id {
	case nav.ScreenHome:
		return home.New(a.styles, a.deps.Version, a.startupDuration, a.currentBento)
	case nav.ScreenConfig:
		return config.New(a.styles, a.currentBento, a.deps.BentoDir)
	case nav.ScreenBentos:
		return bentos.New(a.styles, store)
	case nav.ScreenCreateBento:
		return createbento.New(a.styles, store)
	case nav.ScreenStub:
		if p, ok := nav.ParamsAs[stub.Params](params); ok {
			return stub.New(a.styles, p)
		}
		return stub.New(a.styles, stub.Params{
			Title:   "Stub",
			Message: "Screen is not ready yet.",
		})
	default:
		return stub.New(a.styles, stub.Params{
			Title:   "Unknown route",
			Message: "Requested screen is not registered.",
		})
	}
}
