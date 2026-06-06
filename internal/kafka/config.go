package kafka

const DefaultClientID = "koicha"

// Config contains the Kafka client settings selected by a bento.
//
// Admin/client methods that use franz-go should live in this package later.
type Config struct {
	BootstrapServers []string `yaml:"bootstrapServers"`
	ClientID         string   `yaml:"clientId,omitempty"`
	Auth             Auth     `yaml:"auth"`
}
