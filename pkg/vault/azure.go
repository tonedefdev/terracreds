package vault

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/tonedefdev/terracreds/pkg/helpers"
)

type AzureKeyVault struct {
	SecretList []string
	SecretName string
	UseMSI     bool
	VaultUri   string
}

// getVaultClientMSI returns a keyvault.BaseClient with an MSI authorizer for an Azure Key Vault resource
func getVaultClientMSI() keyvault.BaseClient {
	const vaultUri = "https://vault.azure.net"
	vaultClient := keyvault.New()
	msiConfig := auth.NewMSIConfig()
	msiConfig.Resource = vaultUri

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

// Create stores a secret in an Azure Key Vault
func (akv *AzureKeyVault) Create(secretValue string, method string) error {
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

// Delete removes a secret stored in an Azure Key Vault
func (akv *AzureKeyVault) Delete() error {
	ctx := context.Background()
	client := getVaultClientMSI()

	secret := formatSecretName(akv.SecretName)
	_, err := client.DeleteSecret(ctx, akv.VaultUri, secret)
	return err
}

// Get retrieves a secrete stored in an Azure Key Vault
func (akv *AzureKeyVault) Get() ([]byte, error) {
	ctx := context.Background()
	client := getVaultClientMSI()

	secret := formatSecretName(akv.SecretName)
	get, err := client.GetSecret(ctx, akv.VaultUri, secret, "")
	return []byte(*get.Value), err
}

func (akv *AzureKeyVault) List(secretNames []string) ([]string, error) {
	var secretValues []string
	ctx := context.Background()
	client := getVaultClientMSI()

	for _, secret := range secretNames {
		get, err := client.GetSecret(ctx, akv.VaultUri, secret, "")

		if err != nil {
			return nil, err
		}

		secretValues = append(secretValues, *get.Value)
	}

	return secretValues, nil
}
