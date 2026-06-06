package bento

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/openmdd/koicha/internal/kafka"
)

func TestStoreSaveCreatesBentoFile(t *testing.T) {
	store := Store{Dir: t.TempDir()}
	b := validTestBento("local-dev")

	path, err := store.Save(b)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	wantPath := filepath.Join(store.Dir, "local-dev.yaml")
	if path != wantPath {
		t.Fatalf("Save() path = %q, want %q", path, wantPath)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(data), "name: local-dev") {
		t.Fatalf("saved YAML does not contain metadata name:\n%s", string(data))
	}
}

func TestStoreSaveRefusesOverwrite(t *testing.T) {
	store := Store{Dir: t.TempDir()}
	b := validTestBento("local-dev")

	path, err := store.Save(b)
	if err != nil {
		t.Fatalf("first Save() error = %v", err)
	}
	if err := os.WriteFile(path, []byte("original"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := store.Save(b); err == nil {
		t.Fatal("second Save() error = nil, want overwrite error")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "original" {
		t.Fatalf("existing file was overwritten: %q", string(data))
	}
}

func TestStoreRejectsInvalidName(t *testing.T) {
	store := Store{Dir: t.TempDir()}

	invalidNames := []string{
		"../secret",
		"bad/name",
		"bad.name",
		"-bad",
		"bad-",
	}
	for _, name := range invalidNames {
		t.Run(name, func(t *testing.T) {
			if _, err := store.PathForName(name); err == nil {
				t.Fatal("PathForName() error = nil, want invalid name error")
			}
		})
	}
}

func TestStoreSaveFileMode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows uses ACLs instead of Unix file mode bits")
	}

	store := Store{Dir: t.TempDir()}
	path, err := store.Save(validTestBento("local-dev"))
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("saved file mode = %v, want 0600", got)
	}
}

func TestStoreListReturnsSavedBentos(t *testing.T) {
	store := Store{Dir: t.TempDir()}

	if _, err := store.Save(validTestBento("local-dev")); err != nil {
		t.Fatalf("Save(local-dev) error = %v", err)
	}
	if _, err := store.Save(validTestBento("prod")); err != nil {
		t.Fatalf("Save(prod) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(store.Dir, "notes.txt"), []byte("ignored"), 0o600); err != nil {
		t.Fatalf("WriteFile(notes.txt) error = %v", err)
	}
	if err := os.Mkdir(filepath.Join(store.Dir, "nested"), 0o700); err != nil {
		t.Fatalf("Mkdir(nested) error = %v", err)
	}

	got, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("List() len = %d, want 2", len(got))
	}

	gotNames := []string{got[0].Metadata.Name, got[1].Metadata.Name}
	wantNames := []string{"local-dev", "prod"}
	for i := range wantNames {
		if gotNames[i] != wantNames[i] {
			t.Fatalf("List()[%d].Metadata.Name = %q, want %q", i, gotNames[i], wantNames[i])
		}
	}
}

func TestStoreListRejectsInvalidBentoFile(t *testing.T) {
	store := Store{Dir: t.TempDir()}
	path := filepath.Join(store.Dir, "invalid.yaml")
	if err := os.WriteFile(path, []byte("schemaVersion: v1\nmetadata:\n  name: invalid\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(invalid.yaml) error = %v", err)
	}

	if _, err := store.List(); err == nil {
		t.Fatal("List() error = nil, want validation error")
	}
}

func TestStoreLoadReturnsNamedBento(t *testing.T) {
	store := Store{Dir: t.TempDir()}
	if _, err := store.Save(validTestBento("local-dev")); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := store.Load("local-dev")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got.Metadata.Name != "local-dev" {
		t.Fatalf("Load().Metadata.Name = %q, want %q", got.Metadata.Name, "local-dev")
	}
}

func TestStoreLoadRejectsMetadataNameMismatch(t *testing.T) {
	store := Store{Dir: t.TempDir()}
	path := filepath.Join(store.Dir, "local-dev.yaml")
	data, err := Marshal(validTestBento("prod"))
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := store.Load("local-dev"); err == nil {
		t.Fatal("Load() error = nil, want metadata name mismatch error")
	}
}

func validTestBento(name string) Bento {
	return Bento{
		SchemaVersion: CurrentSchemaVersion,
		Metadata: Metadata{
			Name: name,
		},
		Spec: Spec{
			Kafka: kafka.Config{
				BootstrapServers: []string{"localhost:9092"},
				Auth:             kafka.Auth{Protocol: kafka.SecurityProtocolPlaintext},
			},
			Resources: ResourceView{AllowOutOfScope: true},
		},
	}
}
