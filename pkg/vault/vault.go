package vault

type TerraVault interface {
	Create(secretValue string) error
	Get() ([]byte, error)
}
