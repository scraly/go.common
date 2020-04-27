package secret

import (
	"fmt"

	"github.com/scraly/go.common/pkg/log"

	vault "github.com/hashicorp/vault/api"
	"go.uber.org/zap"
)

// Secret retriever
type Secret struct {
	Token  string
	Client *vault.Client
}

// New initialize a vault connection
func New(token, addr string) (*Secret, error) {
	client, err := vault.NewClient(vault.DefaultConfig())
	if err != nil {
		return nil, err
	}

	// Assign access tokenl
	client.SetToken(token)

	// Set VAULT cluster address
	if err := client.SetAddress(addr); err != nil {
		return nil, err
	}

	return &Secret{
		Client: client,
		Token:  token,
	}, nil
}

// GetFromVault returns a secret from vault
func (secret *Secret) GetFromVault(s string) (string, error) {
	conf, err := secret.Client.Logical().Read(s)
	if err != nil {
		return "", err
	} else if conf == nil {
		log.Bg().Warn("no value found", zap.String("path", s))
		return "", nil
	}

	value, exists := conf.Data["data"]
	if !exists {
		log.Bg().Warn("no 'data' field found (you must add a field with a key named data)", zap.String("path", s))
		return "", nil
	}

	return fmt.Sprintf("%v", value), nil
}
