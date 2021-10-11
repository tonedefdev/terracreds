package vault

type TerraVault interface {
	Create()
	Delete()
	Get(secretName string) ([]byte, error)
}
