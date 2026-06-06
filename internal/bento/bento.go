package bento

import "github.com/openmdd/koicha/internal/kafka"

// SchemaVersion is the bento file format version.
type SchemaVersion string

// CurrentSchemaVersion is the schema version written by Marshal and expected by Validate.
const CurrentSchemaVersion SchemaVersion = "v1"

// Bento is a portable profile that combines Kafka connection settings with the
// resources koicha should show by default.
//
// A bento is stored as a single YAML file. The filename on disk is derived from
// Metadata.Name with a ".yaml" suffix (e.g. "local-dev.yaml").
type Bento struct {
	SchemaVersion SchemaVersion `yaml:"schemaVersion"`
	Metadata      Metadata      `yaml:"metadata"`
	Spec          Spec          `yaml:"spec"`
}

// Metadata describes a bento without affecting Kafka connectivity or resource
// selection.
//
// Name is the machine-readable identifier and matches the bento filename on
// disk (without the ".yaml" extension). It must be unique within a config
// directory and should use only lowercase letters, digits, and hyphens.
type Metadata struct {
	Name        string            `yaml:"name"`
	DisplayName string            `yaml:"displayName,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
}

// Spec contains the selectable runtime profile.
type Spec struct {
	Profile   Profile      `yaml:"profile"`
	Kafka     kafka.Config `yaml:"kafka"`
	Resources ResourceView `yaml:"resources"`
}

// Profile identifies the environment or team context for a bento.
type Profile struct {
	Environment string `yaml:"environment,omitempty"`
}
