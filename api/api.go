package api

// Config struct for terracreds custom configuration
type Config struct {
	Logging struct {
		Enabled bool   `yaml:"enabled"`
		Path    string `yaml:"path"`
	} `yaml:"logging"`
}

// CredentialResponse formatted for consumption by Terraform
type CredentialResponse struct {
	Token string `json:"token"`
}

// Terracreds implements these crud operations
type Terracreds interface {
	Create()
	Delete()
	Get()
}
