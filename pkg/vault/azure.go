package vault

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/tonedefdev/terracreds/api"
	"github.com/tonedefdev/terracreds/pkg/helpers"
)

type AzureKeyVault struct {
	SecretName string `json:"secretName,omitempty"`
	UseMSI     bool   `json:"useMSI,omitempty"`
	VaultUri   string `json:"vaultUri,omitempty"`
}

// getVaultClientMSI returns a keyvault.BaseClient with an MSI authorizer for an Azure Key Vault resource
func getVaultClientMSI() keyvault.BaseClient {
	vaultClient := keyvault.New()
	msiConfig := auth.NewMSIConfig()
	msiConfig.Resource = "https://vault.azure.net"

	authorizer, err := msiConfig.Authorizer()
	if err != nil {
		helpers.CheckError(err)
	}

	vaultClient.Authorizer = authorizer
	return vaultClient
}

// formatSecretName replaces the periods from the hostname with dashes
// since Azure Key Vault can't store secrets that contain periods
func formatSecretName(secretName string) string {
	hostname := strings.Replace(secretName, ".", "-", -1)
	return hostname
}

func (akv AzureKeyVault) Create(secretValue string) error {
	ctx := context.Background()
	client := getVaultClientMSI()

	secretParams := keyvault.SecretSetParameters{
		ContentType: to.StringPtr("password"),
		Value:       &secretValue,
	}

	secret := formatSecretName(akv.SecretName)
	_, err := client.SetSecret(ctx, akv.VaultUri, secret, secretParams)
	return err
}

func (akv AzureKeyVault) Get() ([]byte, error) {
	ctx := context.Background()
	client := getVaultClientMSI()
	secret := formatSecretName(akv.SecretName)

	get, err := client.GetSecret(ctx, akv.VaultUri, secret, "")
	if err == nil {
		response := &api.CredentialResponse{
			Token: string(*get.Value),
		}
		token, err := json.Marshal(response)
		return token, err
	}

	return nil, err
}
