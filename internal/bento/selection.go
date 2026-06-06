package bento

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type SelectionStore struct {
	Path string
}

func (s SelectionStore) SaveSelectedName(name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	if s.Path == "" {
		return errors.New("selected bento path is required")
	}
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o700); err != nil {
		return fmt.Errorf("prepare selected bento dir: %w", err)
	}

	data := []byte(name + "\n")
	if err := os.WriteFile(s.Path, data, 0o600); err != nil {
		return fmt.Errorf("write selected bento file: %w", err)
	}
	if runtime.GOOS != "windows" {
		if err := os.Chmod(s.Path, 0o600); err != nil {
			return fmt.Errorf("secure selected bento file: %w", err)
		}
	}
	return nil
}

func (s SelectionStore) LoadSelectedName() (string, error) {
	if s.Path == "" {
		return "", errors.New("selected bento path is required")
	}

	data, err := os.ReadFile(s.Path)
	switch {
	case err == nil:
	case errors.Is(err, os.ErrNotExist):
		return "", nil
	default:
		return "", fmt.Errorf("read selected bento file: %w", err)
	}

	name := strings.TrimSpace(string(data))
	if name == "" {
		return "", nil
	}
	if err := ValidateName(name); err != nil {
		return "", fmt.Errorf("validate selected bento name: %w", err)
	}
	return name, nil
}
