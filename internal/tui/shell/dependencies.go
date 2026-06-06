package shell

import (
	"time"

	"github.com/openmdd/koicha/internal/bento"
)

type Dependencies struct {
	Version      string
	StartedAt    time.Time
	BentoDir     string
	CurrentPath  string
	CurrentBento *bento.Bento
}
