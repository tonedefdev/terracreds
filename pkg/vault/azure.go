package vault

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/tonedefdev/terracreds/api"
	"github.com/tonedefdev/terracreds/pkg/helpers"
)

type Azure struct {
	UseMSI   bool   `json:"useMSI,omitempty"`
	VaultUri string `json:"vaultUri,omitempty"`
}

func (az Azure) getVaultsClient() keyvault.BaseClient {
	vaultsClient := keyvault.New()
	msiConfig := auth.NewMSIConfig()
	msiConfig.Resource = "keyvault"
	authorizer, err := msiConfig.Authorizer()
	if err != nil {
		helpers.CheckError(err)
	}

	vaultsClient.Authorizer = authorizer
	return vaultsClient
}

func (az Azure) Get(secretName string) ([]byte, error) {
	ctx := context.Background()
	client := az.getVaultsClient()
	get, err := client.GetSecret(ctx, az.VaultUri, secretName, "")
	if err != nil {
		helpers.CheckError(err)
	}

	response := &api.CredentialResponse{
		Token: string(*get.Value),
	}

	token, err := json.Marshal(response)
	return token, err
}
