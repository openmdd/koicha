package kafka

import "context"

// SecretResolver resolves secret material without storing it in bento files.

type Kind string

const (
	KindUsernamePassword Kind = "username_password"
	KindPassphrase       Kind = "passphrase"
)

type Secret interface {
	Kind() Kind
}

type SecretResolver interface {
	Resolve(ctx context.Context, request SecretRequest) (Secret, error)
}

type SecretRequest struct {
	ID     string
	Kind   Kind
	Target string // e.g. "prod/goods-catalog"
	Reason string // e.g. "Kafka SASL authentication"
}

type UsernamePassword struct {
	Username string
	Password []byte
}

func (UsernamePassword) Kind() Kind {
	return KindUsernamePassword
}
