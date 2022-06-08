package vault

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

type AzureKeyVault struct {
	SecretName     string
	SubscriptionId string
	VaultUri       string
}

// getDefaultAzureClient returns a pointer to an azsecrets.Client using the default
// authorization scheme for azidentity
func getDefaultAzureClient(akv *AzureKeyVault) (*azsecrets.Client, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	vaultClient, err := azsecrets.NewClient(akv.VaultUri, cred, nil)
	if err != nil {
		return nil, err
	}

	return vaultClient, err
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
	client, err := getDefaultAzureClient(akv)
	if err != nil {
		return err
	}

	content := "password"
	options := azsecrets.SetSecretOptions{
		ContentType: &content,
	}

	secret := formatSecretName(akv.SecretName)
	_, err = client.SetSecret(ctx, secret, secretValue, &options)
	return err
}

// Delete removes a secret stored in an Azure Key Vault
func (akv *AzureKeyVault) Delete() error {
	ctx := context.Background()
	client, err := getDefaultAzureClient(akv)
	if err != nil {
		return err
	}

	options := azsecrets.BeginDeleteSecretOptions{}
	secret := formatSecretName(akv.SecretName)

	_, err = client.BeginDeleteSecret(ctx, secret, &options)
	return err
}

// Get retrieves a secrete stored in an Azure Key Vault
func (akv *AzureKeyVault) Get() ([]byte, error) {
	ctx := context.Background()
	client, err := getDefaultAzureClient(akv)
	if err != nil {
		return nil, err
	}

	options := azsecrets.GetSecretOptions{}
	secret := formatSecretName(akv.SecretName)

	get, err := client.GetSecret(ctx, secret, &options)
	if err != nil {
		return nil, err
	}
	return []byte(*get.Value), err
}

func (akv *AzureKeyVault) List(secretNames []string) ([]string, error) {
	var secretValues []string
	ctx := context.Background()
	client, err := getDefaultAzureClient(akv)
	if err != nil {
		return nil, err
	}

	for _, secret := range secretNames {
		options := azsecrets.GetSecretOptions{}
		secret := formatSecretName(secret)

		get, err := client.GetSecret(ctx, secret, &options)
		if err != nil {
			return nil, err
		}

		secretValues = append(secretValues, *get.Value)
	}

	return secretValues, nil
}
