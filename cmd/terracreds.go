package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/MakeNowJust/heredoc"
	"github.com/fatih/color"
	"github.com/tonedefdev/terracreds/api"
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/tonedefdev/terracreds/pkg/platform"
	"github.com/tonedefdev/terracreds/pkg/vault"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

// TerraCreds interface implements these methods for a credential's lifecycle
type TerraCreds interface {
	// Create or store a secret in a vault
	Create(cfg *api.Config, hostname string, token any, user *user.User, vault vault.TerraVault) error
	// Delete or forget a secret in a vault
	Delete(cfg *api.Config, command string, hostname string, user *user.User, vault vault.TerraVault) error
	// Get or retrieve a secret in a vault
	Get(cfg *api.Config, hostname string, user *user.User, vault vault.TerraVault) ([]byte, error)
	// List the secrets from within a vault
	List(c *cli.Context, cfg *api.Config, secretNames []string, user *user.User, vault vault.TerraVault) ([]string, error)
}

// CopyTerraCreds will create a copy of the binary to the destination path.
func CopyTerraCreds(dest string) error {
	from, err := os.Open(string(os.Args[0]))
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}

	fmt.Fprintf(color.Output, "%s: Copied binary '%s' to '%s'\n", color.CyanString("INFO"), string(os.Args[0]), dest)
	return err
}

// GenerateTerracreds creates the binary to use this package as a credential helper and optionally the terraform.rc file
func GenerateTerraCreds(c *cli.Context, version string, confirm string) error {
	var cliConfig string
	var tfPlugins string
	var binary string

	if runtime.GOOS == "windows" {
		userProfile := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		cliConfig = filepath.Join(userProfile, "terraform.rc")
		tfPlugins = filepath.Join(userProfile, "terraform.d", "plugins")
		binary = filepath.Join(tfPlugins, "terraform-credentials-terracreds.exe")
	}

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		userProfile := os.Getenv("HOME")
		cliConfig = filepath.Join(userProfile, ".terraform.d", ".terraformrc")
		tfPlugins = filepath.Join(userProfile, ".terraform.d", "plugins")
		binary = filepath.Join(tfPlugins, "terraform-credentials-terracreds")
	}

	err := helpers.NewDirectory(tfPlugins)
	if err != nil {
		return err
	}

	err = CopyTerraCreds(binary)
	if err != nil {
		return err
	}

	if c.Bool("create-cli-config") == true {
		const verbiage = "This command will delete any settings in your .terraformrc file\n\n    Enter 'yes' to coninue or press 'enter' or 'return' to cancel: "
		fmt.Fprintf(color.Output, "%s: %s", color.YellowString("WARNING"), verbiage)
		fmt.Scanln(&confirm)
		fmt.Print("\n")

		if confirm == "yes" {
			doc := heredoc.Doc(`
			credentials_helper "terracreds" {
				args = []
			}`)

			err := helpers.WriteToFile(cliConfig, doc)
			return err
		}
	}

	return nil
}

// GetSecretName returns the name of the secret from the config or returns the hostname value from the CLI
func GetSecretName(cfg *api.Config, hostname string) string {
	if cfg.Aws.SecretName != "" {
		return cfg.Aws.SecretName
	}
	if cfg.Azure.SecretName != "" {
		return cfg.Azure.SecretName
	}
	if cfg.HashiVault.SecretName != "" {
		return cfg.HashiVault.SecretName
	}
	return hostname
}

// InitTerraCreds initializes the configuration for Terracreds
func (cmd *Config) InitTerraCreds() {
	if cmd.ConfigFile.EnvironmentValue != "" {
		cmd.ConfigFile.Path = filepath.Join(cmd.ConfigFile.EnvironmentValue, cmd.ConfigFile.Name)
	} else {
		binPath := helpers.GetBinaryPath(os.Args[0], runtime.GOOS)
		cmd.ConfigFile.Path = filepath.Join(binPath, cmd.ConfigFile.Name)
	}

	err := helpers.CreateConfigFile(cmd.ConfigFile.Path)
	if err != nil {
		helpers.CheckError(err)
	}

	err = cmd.LoadConfig(cmd.ConfigFile.Path)
	if err != nil {
		helpers.CheckError(err)
	}
}

// LoadConfig loads the config file if it exists
func (cmd *Config) LoadConfig(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, &cmd.Cfg)
	return err
}

// NewTerraCreds is the constructor to create a TerraCreds interface
func NewTerraCreds(os string) TerraCreds {
	switch os {
	case "darwin":
		return &platform.Mac{}
	case "linux":
		return &platform.Linux{}
	case "windows":
		return &platform.Windows{}
	default:
		return nil
	}
}

// NewTerrVault is the constructor to create a TerraVault interface for the vault provider defined in the Cfg
func (cmdCfg *Config) NewTerraVault(hostname string) vault.TerraVault {
	if cmdCfg.Cfg.Aws.Region != "" {
		vault := &vault.AwsSecretsManager{
			Description: cmdCfg.Cfg.Aws.Description,
			Region:      cmdCfg.Cfg.Aws.Region,
			SecretName:  hostname,
		}

		if cmdCfg.Cfg.Aws.SecretName != "" {
			vault.SecretName = cmdCfg.Cfg.Aws.SecretName
		}

		return vault
	}

	if cmdCfg.Cfg.Azure.VaultUri != "" {
		vault := &vault.AzureKeyVault{
			SecretName:     hostname,
			SubscriptionId: cmdCfg.Cfg.Azure.SubscriptionId,
			VaultUri:       cmdCfg.Cfg.Azure.VaultUri,
		}

		if cmdCfg.Cfg.Azure.SecretName != "" {
			vault.SecretName = cmdCfg.Cfg.Azure.SecretName
		}

		return vault
	}

	if cmdCfg.Cfg.GCP.ProjectId != "" {
		vault := &vault.GCPSecretManager{
			ProjectId: cmdCfg.Cfg.GCP.ProjectId,
			SecretId:  hostname,
		}

		if cmdCfg.Cfg.GCP.SecretId != "" {
			vault.SecretId = cmdCfg.Cfg.GCP.SecretId
		}

		return vault
	}

	if cmdCfg.Cfg.HashiVault.VaultUri != "" {
		vault := &vault.HashiVault{
			EnvTokenName: cmdCfg.Cfg.HashiVault.EnvironmentTokenName,
			KeyVaultPath: cmdCfg.Cfg.HashiVault.KeyVaultPath,
			SecretName:   hostname,
			SecretPath:   cmdCfg.Cfg.HashiVault.SecretPath,
			VaultUri:     cmdCfg.Cfg.HashiVault.VaultUri,
		}

		if cmdCfg.Cfg.HashiVault.SecretName != "" {
			vault.SecretName = cmdCfg.Cfg.HashiVault.SecretName
		}

		return vault
	}

	return nil
}
