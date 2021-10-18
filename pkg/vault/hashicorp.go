package vault

import (
	"fmt"
	"os"

	hcvault "github.com/hashicorp/vault/api"
	"github.com/tonedefdev/terracreds/pkg/helpers"
)

type HashiVault struct {
	EnvTokenName string
	KeyVaultPath string
	SecretName   string
	SecretPath   string
	VaultUri     string
}

func (hc *HashiVault) newHashiVaultClient() *hcvault.Client {
	config := hcvault.DefaultConfig()
	config.Address = hc.VaultUri

	client, err := hcvault.NewClient(config)
	if err != nil {
		helpers.CheckError(err)
	}

	client.SetToken(os.Getenv(hc.EnvTokenName))

	return client
}

func (hc *HashiVault) Create(secretValue string) error {
	client := hc.newHashiVaultClient()
	secret := make(map[string]interface{})

	key := hc.SecretName
	data := make(map[string]interface{})
	data[key] = secretValue
	secret["data"] = data

	kvPath := fmt.Sprintf("%s/data/%s", hc.KeyVaultPath, hc.SecretPath)
	_, err := client.Logical().Write(kvPath, secret)
	if err != nil {
		return err
	}

	return err
}

func (hc *HashiVault) Delete() error {
	client := hc.newHashiVaultClient()

	kvPath := fmt.Sprintf("%s/data/%s", hc.KeyVaultPath, hc.SecretPath)
	_, err := client.Logical().Delete(kvPath)
	if err != nil {
		return err
	}

	return err
}

func (hc *HashiVault) Get() ([]byte, error) {
	client := hc.newHashiVaultClient()

	kvPath := fmt.Sprintf("%s/data/%s", hc.KeyVaultPath, hc.SecretPath)
	secret, err := client.Logical().Read(kvPath)

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data type assertion failed: %T %#v", secret.Data["data"], secret.Data["data"])
	}

	key := hc.SecretName
	value, ok := data[key].(string)
	if !ok {
		return nil, fmt.Errorf("value type assertion failed: %T %#v", data[key], data[key])
	}

	return []byte(value), err
}
