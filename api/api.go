package api

// Aws is the configuration structure for the AWS vault provider
type Aws struct {
	// Description (Optional) is a description to provide to the secret
	Description string `yaml:"description,omitempty"`

	// Region (Required) is the region where AWS Secrets Manager is hosted
	Region string `yaml:"region,omitempty"`

	// SecretName (Optional) is the friendly name of the secret stored in AWS Secrets Manager
	// if omitted Terracreds will use the hostname value instead
	SecretName string `yaml:"secretName,omitempty"`
}

// Azure is the configuration structure for the Azure vault provider
type Azure struct {
	// SecretName (Optional) is the name of the secret stored in Azure Key Vault
	// if omitted Terracreds will use the hostname value instead
	SecretName string `yaml:"secretName,omitempty"`

	// UseMSI (Required) is a flag to indicate if the Managed Identity of the Azure VM should be used for authentication
	UseMSI bool `yaml:"useMSI,omitempty"`

	// VaultUri (Required) is the FQDNS of the Azure Key Vault resource
	VaultUri string `yaml:"vaultUri,omitempty"`
}

type HashiVault struct {
	EnvTokenName string
	SecretName   string
	VaultUri     string
}

// Config struct for terracreds custom configuration
type Config struct {
	Logging    Logging    `yaml:"logging"`
	Aws        Aws        `yaml:"aws,omitempty"`
	Azure      Azure      `yaml:"azure,omitempty"`
	HashiVault HashiVault `yaml:"hashiVault,omitempty"`
}

// Logging struct defines the parameters for logging
type Logging struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

// CredentialResponse formatted for consumption by Terraform
type CredentialResponse struct {
	Token string `json:"token"`
}
