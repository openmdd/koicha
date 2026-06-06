package bento

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSelectionStoreLoadMissingFile(t *testing.T) {
	store := SelectionStore{Path: filepath.Join(t.TempDir(), "current-bento")}

	name, err := store.LoadSelectedName()
	if err != nil {
		t.Fatalf("LoadSelectedName() error = %v", err)
	}
	if name != "" {
		t.Fatalf("LoadSelectedName() = %q, want empty name", name)
	}
}

func TestSelectionStoreSaveAndLoadSelectedName(t *testing.T) {
	store := SelectionStore{Path: filepath.Join(t.TempDir(), ".caches", "current-bento")}

	if err := store.SaveSelectedName("local-dev"); err != nil {
		t.Fatalf("SaveSelectedName() error = %v", err)
	}

	name, err := store.LoadSelectedName()
	if err != nil {
		t.Fatalf("LoadSelectedName() error = %v", err)
	}
	if name != "local-dev" {
		t.Fatalf("LoadSelectedName() = %q, want %q", name, "local-dev")
	}
}

func TestSelectionStoreSaveRejectsInvalidName(t *testing.T) {
	store := SelectionStore{Path: filepath.Join(t.TempDir(), "current-bento")}

	if err := store.SaveSelectedName("../secret"); err == nil {
		t.Fatal("SaveSelectedName() error = nil, want invalid name error")
	}
}

func TestSelectionStoreLoadRejectsInvalidName(t *testing.T) {
	store := SelectionStore{Path: filepath.Join(t.TempDir(), "current-bento")}
	if err := os.WriteFile(store.Path, []byte("../secret\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := store.LoadSelectedName(); err == nil {
		t.Fatal("LoadSelectedName() error = nil, want invalid name error")
	}
}

func TestSelectionStoreSaveFileMode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows uses ACLs instead of Unix file mode bits")
	}

	store := SelectionStore{Path: filepath.Join(t.TempDir(), "current-bento")}
	if err := store.SaveSelectedName("local-dev"); err != nil {
		t.Fatalf("SaveSelectedName() error = %v", err)
	}

	info, err := os.Stat(store.Path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("selected bento file mode = %v, want 0600", got)
	}
}
