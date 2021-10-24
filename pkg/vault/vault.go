package vault

// TerraVault implements an interface that handles secret lifecycle mananagement
// for a credential vault provider
type TerraVault interface {
	Create(secretValue string, method string) error
	Delete() error
	Get() ([]byte, error)
}
