package main

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/tonedefdev/terracreds/api"
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/tonedefdev/terracreds/pkg/platform"
	"github.com/tonedefdev/terracreds/pkg/vault"
)

// Terracreds interface implements these methods for a credential's lifecycle
type Terracreds interface {
	// Create or store an API token in a vault
	Create(cfg api.Config, hostname string, token interface{}, user *user.User, vault vault.TerraVault) error
	// Delete or forget an API token in a vault
	Delete(cfg api.Config, command string, hostname string, user *user.User, vault vault.TerraVault) error
	// Get or retrieve an API token in a vault
	Get(cfg api.Config, hostname string, user *user.User, vault vault.TerraVault) ([]byte, error)
}

// returnProvider returns the correct struct for the specific operating system
func returnProvider(os string) Terracreds {
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

// returnsVaultProvider handles returning the correct vault provider based on the
// api.Config struct
func returnVaultProvider(cfg *api.Config, hostname string) vault.TerraVault {
	if cfg.Aws.Region != "" {
		vault := &vault.AwsSecretsManager{
			Description: cfg.Aws.Description,
			Region:      cfg.Aws.Region,
			SecretName:  hostname,
		}

		// Fallback to the cfg's secret name if it isn't an empty string
		if cfg.Aws.SecretName != "" {
			vault.SecretName = cfg.Aws.SecretName
		}

		return vault
	}

	if cfg.Azure.VaultUri != "" {
		vault := &vault.AzureKeyVault{
			SecretName: hostname,
			UseMSI:     cfg.Azure.UseMSI,
			VaultUri:   cfg.Azure.VaultUri,
		}

		// Fallback to the cfg's secret name if it isn't an empty string
		if cfg.Azure.SecretName != "" {
			vault.SecretName = cfg.Azure.SecretName
		}

		return vault
	}

	if cfg.HashiVault.VaultUri != "" {
		vault := &vault.HashiVault{
			EnvTokenName: cfg.HashiVault.EnvironmentTokenName,
			KeyVaultPath: cfg.HashiVault.KeyVaultPath,
			SecretName:   hostname,
			SecretPath:   cfg.HashiVault.SecretPath,
			VaultUri:     cfg.HashiVault.VaultUri,
		}

		if cfg.HashiVault.SecretName != "" {
			vault.SecretName = cfg.HashiVault.SecretName
		}

		return vault
	}

	return nil
}

func main() {
	var cfg api.Config
	version := "2.0.0"

	err := helpers.LoadConfig(&cfg)
	if err != nil {
		helpers.CheckError(err)
	}

	provider := returnProvider(runtime.GOOS)
	if provider == nil {
		fmt.Fprintf(color.Output, "%s: Terracreds cannot run on this platform: '%s'\n", color.RedString("ERROR"), runtime.GOOS)
		return
	}

	app := &cli.App{
		Name:      "terracreds",
		Usage:     "a credential helper for Terraform Cloud/Enterprise that leverages your vault provider of choice for securely storing your API tokens or other secrets.\n\n   Visit https://github.com/tonedefdev/terracreds for more information",
		UsageText: "Directly store credentials from Terraform using 'terraform login' or manually store them using 'terracreds create -n app.terraform.io -t myAPItoken'",
		Version:   version,
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Manually create or update a credential object in the vault provider of your choice that contains the Terraform Cloud/Enterprise authorization token or another secret",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "hostname",
						Aliases: []string{"n"},
						Value:   "place_holder",
						Usage:   "The name of the Terraform Cloud/Enterprise server's hostname or the name of the secret. This is also the display name of the credential object",
					},
					&cli.StringFlag{
						Name:    "apiToken",
						Aliases: []string{"t"},
						Value:   "",
						Usage:   "The Terraform Cloud/Enterprise API authorization token or other secret value to be securely stored in your vault provider of choice",
					},
				},
				Action: func(c *cli.Context) error {
					if len(os.Args) == 2 {
						fmt.Fprintf(color.Output, "%s: No hostname or token was specified. Use 'terracreds create -h' to print help info\n", color.RedString("ERROR"))
						return nil
					}

					vaultProvider := returnVaultProvider(&cfg, c.String("hostname"))
					hostname := helpers.GetSecretName(&cfg, c.String("hostname"))

					user, err := user.Current()
					helpers.CheckError(err)
					err = Terracreds.Create(provider, cfg, hostname, c.String("apiToken"), user, vaultProvider)
					if err != nil {
						helpers.CheckError(err)
					}

					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "Delete a stored credential in the vault provider of your choice",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "hostname",
						Aliases: []string{"n"},
						Value:   "place_holder",
						Usage:   "The name of the Terraform Cloud/Enterprise server's hostname or the name of the secret. This is also the display name of the credential object.",
					},
				},
				Action: func(c *cli.Context) error {
					if len(os.Args) == 2 {
						fmt.Fprintf(color.Output, "%s: No hostname was specified. Use 'terracreds delete -h' for help info\n", color.RedString("ERROR"))
						return nil
					}

					if !strings.Contains(os.Args[2], "-n") && !strings.Contains(os.Args[2], "--hostname") {
						msg := fmt.Sprintf("A hostname was not expected here: '%s'", os.Args[2])
						helpers.Logging(cfg, msg, "WARNING")
						fmt.Fprintf(color.Output, "%s: %s Did you mean `terracreds delete --hostname/-n %s'?\n", color.YellowString("WARNING"), msg, os.Args[2])
						return nil
					}

					vaultProvider := returnVaultProvider(&cfg, c.String("hostname"))
					hostname := helpers.GetSecretName(&cfg, c.String("hostname"))

					user, err := user.Current()
					helpers.CheckError(err)
					err = Terracreds.Delete(provider, cfg, os.Args[1], hostname, user, vaultProvider)
					if err != nil {
						helpers.CheckError(err)
					}

					return nil
				},
			},
			{
				Name:  "forget",
				Usage: "(Terraform Only) Forget a stored credential in your vault provider of choice when 'terraform logout' has been called",
				Action: func(c *cli.Context) error {
					if len(os.Args) == 2 {
						fmt.Fprintf(color.Output, "%s: No hostname was specified. Use 'terracreds forget -h' for help info\n", color.RedString("ERROR"))
						return nil
					}

					vaultProvider := returnVaultProvider(&cfg, os.Args[2])
					hostname := helpers.GetSecretName(&cfg, os.Args[2])

					user, err := user.Current()
					helpers.CheckError(err)
					err = Terracreds.Delete(provider, cfg, os.Args[1], hostname, user, vaultProvider)
					if err != nil {
						helpers.CheckError(err)
					}

					return nil
				},
			},
			{
				Name:  "generate",
				Usage: "Generate the folders and plugin binary required to leverage terracreds as a Terraform credential helper",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "create-cli-config",
						Value: false,
						Usage: "Creates the Terraform CLI config with a terracreds credential helper block. This will overwrite the existing file if it already exists.",
					},
				},
				Action: func(c *cli.Context) error {
					helpers.GenerateTerracreds(c)
					return nil
				},
			},
			{
				Name:  "get",
				Usage: "Get the credential object value by passing the hostname of the Terraform Cloud/Enterprise server as an argument or the name of the secret. The credential is returned as a JSON object and formatted for consumption by Terraform",
				Action: func(c *cli.Context) error {
					if len(os.Args) > 2 {
						user, err := user.Current()
						helpers.CheckError(err)

						vaultProvider := returnVaultProvider(&cfg, os.Args[2])
						hostname := helpers.GetSecretName(&cfg, os.Args[2])

						token, err := Terracreds.Get(provider, cfg, hostname, user, vaultProvider)
						if err != nil {
							helpers.CheckError(err)
						}
						fmt.Println(string(token))
						return nil
					}

					msg := "A hostname was expected after the 'get' command but no argument was provided"
					helpers.Logging(cfg, msg, "ERROR")
					fmt.Fprintf(color.Output, "%s: %s\n", color.RedString("ERROR"), msg)
					return nil
				},
			},
			{
				Name:  "store",
				Usage: "(Terraform Only) Store or update a credential object in your vault provider of choice when 'terraform login' has been called",
				Action: func(c *cli.Context) error {
					if len(os.Args) == 2 {
						fmt.Fprintf(color.Output, "%s: No hostname or token was specified. Use 'terracreds store -h' to print help info\n", color.RedString("ERROR"))
						return nil
					}

					vaultProvider := returnVaultProvider(&cfg, os.Args[2])
					hostname := helpers.GetSecretName(&cfg, os.Args[2])

					user, err := user.Current()
					helpers.CheckError(err)
					err = Terracreds.Create(provider, cfg, hostname, nil, user, vaultProvider)
					if err != nil {
						helpers.CheckError(err)
					}

					return nil
				},
			},
		},
	}

	err = app.Run(os.Args)
	helpers.CheckError(err)
}
