package api

// Aws is the configuration structure for the AWS vault provider
type Aws struct {
	// Description (Optional) A description to provide to the secret
	Description string `yaml:"description,omitempty"`

	// Region (Required) The region where AWS Secrets Manager is hosted
	Region string `yaml:"region,omitempty"`

	// SecretName (Optional) The friendly name of the secret stored in AWS Secrets Manager
	// if omitted Terracreds will use the hostname value instead
	SecretName string `yaml:"secretName,omitempty"`
}

// Azure is the configuration structure for the Azure vault provider
type Azure struct {
	// SecretName (Optional) The name of the secret stored in Azure Key Vault
	// if omitted Terracreds will use the hostname value instead
	SecretName string `yaml:"secretName,omitempty"`

	// SubscriptionId (Required) The subscription ID where the target Key Vault has been created
	SubscriptionId string `yaml:"subscriptionId,omitempty"`

	// VaultUri (Required) The FQDN of the Azure Key Vault resource
	VaultUri string `yaml:"vaultUri,omitempty"`
}

// GCP is the configuration structure for the Goocle Cloud Secret Manager provider
type GCP struct {
	// ProjectId (Required) The name of the GCP project where the Secret Manager API has been enabled
	ProjectId string `yaml:"projectId,omitempty"`

	// SecretId (Optional) The name of the secret to create
	SecretId string `yaml:"secretId,omitempty"`
}

// HCVault is the configuration structure for the Hashicorp Vault provider
type HCVault struct {
	// EnvironmentTokenName (Required) The name of the environment variable that currently holds
	// the Vault token
	EnvironmentTokenName string `yaml:"environmentTokenName,omitempty"`

	// KeyVaultPath (Required) The name of the Key Vault store inside of Vault
	KeyVaultPath string `yaml:"keyVaultPath,omitempty"`

	// SecretName (Optional) The name of the secret stored inside of Vault
	// if omitted Terracreds will use the hostname value instead
	SecretName string `yaml:"secretName,omitempty"`

	// SecretPath (Required) The path to the secret itself inside of Vault
	SecretPath string `yaml:"secretPath,omitempty"`

	// VaultUri (Required) The URL of the Vault instance including its port
	VaultUri string `yaml:"vaultUri,omitempty"`
}

// Config struct for terracreds custom configuration
type Config struct {
	Logging    Logging  `yaml:"logging"`
	Aws        Aws      `yaml:"aws,omitempty"`
	Azure      Azure    `yaml:"azure,omitempty"`
	HashiVault HCVault  `yaml:"hcvault,omitempty"`
	GCP        GCP      `yaml:"gcp,omitempty"`
	Secrets    []string `yaml:"secrets,omitempty"`
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
