package vault

import (
	"fmt"
	"os"

	hcvault "github.com/hashicorp/vault/api"
	"github.com/tonedefdev/terracreds/pkg/helpers"
)

type HashiVault struct {
	EnvTokenName string
	SecretName   string
	VaultUri     string
}

func (hc *HashiVault) newHashiVaultClient() *hcvault.Client {
	config := hcvault.DefaultConfig()
	config.Address = hc.VaultUri

	client, err := hcvault.NewClient(config)
	if err != nil {
		helpers.CheckError(err)
	}

	return client
}

func (hc *HashiVault) Create(secretValue string) error {
	return nil
}

func (hc *HashiVault) Delete() error {
	return nil
}

func (hc *HashiVault) Get() ([]byte, error) {
	client := hc.newHashiVaultClient()
	client.SetToken(os.Getenv(hc.EnvTokenName))

	secret, err := client.Logical().Read("kv-v2/data/creds")

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data type assertion failed: %T %#v", secret.Data["data"], secret.Data["data"])
	}

	key := "password"
	value, ok := data[key].(string)
	if !ok {
		return nil, fmt.Errorf("value type assertion failed: %T %#v", data[key], data[key])
	}

	return []byte(value), err
}
