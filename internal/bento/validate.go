package bento

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/openmdd/koicha/internal/kafka"
)

// supportedVersions lists every SchemaVersion that this build can use.
var supportedVersions = map[SchemaVersion]struct{}{
	CurrentSchemaVersion: {},
}

// namePattern is the allowed format for Metadata.Name: lowercase letters,
// digits, and hyphens only (no leading or trailing hyphens).
var namePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9\-]*[a-z0-9]$|^[a-z0-9]$`)

// supportedProtocols lists the SecurityProtocol values that koicha can open a
// client connection for.
var supportedProtocols = map[kafka.SecurityProtocol]struct{}{
	kafka.SecurityProtocolPlaintext: {},
}

// Validate returns an error if b is not a valid, usable Bento.
//
// Checks performed:
//   - SchemaVersion is a supported version.
//   - Metadata.Name is non-empty and matches [a-z0-9][a-z0-9-]*.
//   - Spec.Kafka.BootstrapServers contains at least one address.
//   - Spec.Kafka.Auth.Protocol is a known SecurityProtocol.
func Validate(b Bento) error {
	var errs []error

	if _, ok := supportedVersions[b.SchemaVersion]; !ok {
		errs = append(errs, fmt.Errorf("unsupported schemaVersion %q: expected %q", b.SchemaVersion, CurrentSchemaVersion))
	}

	if err := ValidateName(b.Metadata.Name); err != nil {
		errs = append(errs, err)
	}

	if len(b.Spec.Kafka.BootstrapServers) == 0 {
		errs = append(errs, errors.New("spec.kafka.bootstrapServers must contain at least one address"))
	}

	if _, ok := supportedProtocols[b.Spec.Kafka.Auth.Protocol]; !ok {
		errs = append(errs, fmt.Errorf("unknown spec.kafka.auth.protocol %q", b.Spec.Kafka.Auth.Protocol))
	}

	return errors.Join(errs...)
}

func ValidateName(name string) error {
	if name == "" {
		return errors.New("metadata.name must not be empty")
	}
	if !namePattern.MatchString(name) {
		return fmt.Errorf("metadata.name %q must contain only lowercase letters, digits, and hyphens", name)
	}
	return nil
}
