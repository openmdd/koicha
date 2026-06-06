package appconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Paths struct {
	ConfigDir        string
	BentoDir         string
	CacheDir         string
	CurrentBentoPath string
}

func Init() (Paths, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return Paths{}, fmt.Errorf("resolve user config dir: %w", err)
	}

	paths := Paths{
		ConfigDir: filepath.Join(baseDir, "koicha"),
	}
	paths.BentoDir = filepath.Join(paths.ConfigDir, "bento")
	paths.CacheDir = filepath.Join(paths.ConfigDir, ".caches")
	paths.CurrentBentoPath = filepath.Join(paths.CacheDir, "current-bento")

	if err := secureDir(paths.ConfigDir); err != nil {
		return Paths{}, fmt.Errorf("prepare koicha config dir: %w", err)
	}
	if err := secureDir(paths.BentoDir); err != nil {
		return Paths{}, fmt.Errorf("prepare koicha bento dir: %w", err)
	}
	if err := secureDir(paths.CacheDir); err != nil {
		return Paths{}, fmt.Errorf("prepare koicha cache dir: %w", err)
	}
	return paths, nil
}

func secureDir(path string) error {
	if err := os.MkdirAll(path, 0o700); err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		return nil
	}
	return os.Chmod(path, 0o700)
}
