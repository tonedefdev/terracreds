package api

type Azure struct {
	SecretName string `yaml:"secretName,omitempty"`
	UseMSI     bool   `yaml:"useMSI,omitempty"`
	VaultUri   string `yaml:"vaultUri,omitempty"`
}

// Config struct for terracreds custom configuration
type Config struct {
	Logging Logging `yaml:"logging"`
	Azure   Azure   `yaml:"azure,omitempty"`
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
