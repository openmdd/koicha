package kafka

// SecurityProtocol is Kafka's transport/auth protocol family.
type SecurityProtocol string

const (
	// SecurityProtocolPlaintext is an unauthenticated, unencrypted connection.
	SecurityProtocolPlaintext SecurityProtocol = "PLAINTEXT"
)

// Auth describes how koicha should authenticate to Kafka.
//
// Only PLAINTEXT is modeled for now. Auth variants that require secrets (e.g.
// SASL_SSL) should reference a dedicated secret provider rather than storing
// credentials directly in bento files.
type Auth struct {
	Protocol SecurityProtocol `yaml:"protocol"`
}
