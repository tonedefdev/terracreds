package vault

import (
	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/keyvault/mgmt/keyvault"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/tonedefdev/terracreds/pkg/helpers"
)

type Azure struct {
	Client         keyvault.BaseClient
	SubscriptionID string
	VaultUri       string
}

func (akv *Azure) getKeyVaultClient() {
	akv.Client = keyvault.New(akv.SubscriptionID)

	msiConfig := auth.NewMSIConfig()
	authorizer, err := msiConfig.Authorizer()
	if err != nil {
		helpers.CheckError(err)
	}

	akv.Client.Authorizer = authorizer
	akv.Client.BaseURI = akv.VaultUri
}

func (akv *Azure) Get() {

}
