package bento

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Store struct {
	Dir string
}

func (s Store) PathForName(name string) (string, error) {
	if err := ValidateName(name); err != nil {
		return "", err
	}
	if s.Dir == "" {
		return "", errors.New("bento store dir is required")
	}
	return filepath.Join(s.Dir, name+".yaml"), nil
}

func (s Store) Exists(name string) (bool, error) {
	path, err := s.PathForName(name)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	default:
		return false, fmt.Errorf("check bento file: %w", err)
	}
}

func (s Store) Save(b Bento) (string, error) {
	if err := Validate(b); err != nil {
		return "", err
	}

	path, err := s.PathForName(b.Metadata.Name)
	if err != nil {
		return "", err
	}

	data, err := Marshal(b)
	if err != nil {
		return "", fmt.Errorf("marshal bento: %w", err)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return "", fmt.Errorf("bento %q already exists", b.Metadata.Name)
		}
		return "", fmt.Errorf("create bento file: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = os.Remove(path)
		}
	}()

	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return "", fmt.Errorf("write bento file: %w", err)
	}
	if err := file.Close(); err != nil {
		return "", fmt.Errorf("close bento file: %w", err)
	}

	committed = true
	return path, nil
}

func (s Store) Load(name string) (Bento, error) {
	path, err := s.PathForName(name)
	if err != nil {
		return Bento{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Bento{}, fmt.Errorf("read bento file %q: %w", path, err)
	}

	b, err := Unmarshal(data)
	if err != nil {
		return Bento{}, fmt.Errorf("unmarshal bento file %q: %w", path, err)
	}
	if err := Validate(b); err != nil {
		return Bento{}, fmt.Errorf("validate bento file %q: %w", path, err)
	}
	if b.Metadata.Name != name {
		return Bento{}, fmt.Errorf("bento file %q contains metadata.name %q, want %q", path, b.Metadata.Name, name)
	}
	return b, nil
}

func (s Store) List() ([]Bento, error) {
	if s.Dir == "" {
		return nil, errors.New("bento store dir is required")
	}

	dirEntries, err := os.ReadDir(s.Dir)
	if err != nil {
		return nil, fmt.Errorf("read bento dir: %w", err)
	}

	bentos := make([]Bento, 0, len(dirEntries))

	for _, entry := range dirEntries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		path := filepath.Join(s.Dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read bento file %q: %w", path, err)
		}

		b, err := Unmarshal(data)
		if err != nil {
			return nil, fmt.Errorf("unmarshal bento file %q: %w", path, err)
		}
		if err := Validate(b); err != nil {
			return nil, fmt.Errorf("validate bento file %q: %w", path, err)
		}

		bentos = append(bentos, b)
	}
	return bentos, nil
}
