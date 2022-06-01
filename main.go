package main

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"

	"github.com/tonedefdev/terracreds/api"
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/tonedefdev/terracreds/pkg/platform"
	"github.com/tonedefdev/terracreds/pkg/vault"
)

var (
	cfg            api.Config
	configFilePath string
)

// TerraCreds interface implements these methods for a credential's lifecycle
type TerraCreds interface {
	// Create or store a secret in a vault
	Create(cfg api.Config, hostname string, token interface{}, user *user.User, vault vault.TerraVault) error
	// Delete or forget a secret in a vault
	Delete(cfg api.Config, command string, hostname string, user *user.User, vault vault.TerraVault) error
	// Get or retrieve a secret in a vault
	Get(cfg api.Config, hostname string, user *user.User, vault vault.TerraVault) ([]byte, error)
	// List the secrets from within a vault
	List(c *cli.Context, cfg api.Config, secretNames []string, user *user.User, vault vault.TerraVault) ([]string, error)
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

// NewTerrVault is the constructor to create a TerraVault interface
// for the vault provider defined in the cfg
func NewTerraVault(cfg *api.Config, hostname string) vault.TerraVault {
	if cfg.Aws.Region != "" {
		vault := &vault.AwsSecretsManager{
			Description: cfg.Aws.Description,
			Region:      cfg.Aws.Region,
			SecretName:  hostname,
		}

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
	version := "2.1.0"

	fileEnvVar := os.Getenv("TC_CONFIG_PATH")
	if fileEnvVar != "" {
		configFilePath = fileEnvVar + "config.yaml"
	} else {
		binPath := helpers.GetBinaryPath(os.Args[0], runtime.GOOS)
		configFilePath = binPath + "config.yaml"
	}

	err := helpers.CreateConfigFile(configFilePath)
	if err != nil {
		helpers.CheckError(err)
	}

	err = helpers.LoadConfig(configFilePath, &cfg)
	if err != nil {
		helpers.CheckError(err)
	}

	terraCreds := NewTerraCreds(runtime.GOOS)
	if terraCreds == nil {
		fmt.Fprintf(color.Output, "%s: terracreds cannot run on this platform: '%s'\n", color.RedString("ERROR"), runtime.GOOS)
		return
	}

	app := &cli.App{
		Name:      "terracreds",
		Usage:     "a credential helper for Terraform Cloud/Enterprise that leverages your vault provider of choice for securely storing your API tokens or other secrets.\n\n   Visit https://github.com/tonedefdev/terracreds for more information",
		UsageText: "Directly store credentials from Terraform using 'terraform login' or manually store them using 'terracreds create -n app.terraform.io -t myAPItoken'",
		Version:   version,
		Commands: []*cli.Command{
			{
				Name:  "config",
				Usage: "View or modify the Terracreds configuration file",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "use-local-vault",
						Usage:    "WARNING: Resets configuration to only use the local operating system's credential vault. This will delete all configuration values for cloud provider vaults from the config file",
						Required: false,
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "aws",
						Usage: "AWS Secret Managers provider configuration settings",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "description",
								Usage:    "A description to provide to the secret",
								Required: false,
							},
							&cli.StringFlag{
								Name:     "region",
								Usage:    "The region where AWS Secrets Manager is hosted",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "secret-name",
								Usage:    "The friendly name of the secret stored in AWS Secrets Manager. If omitted Terracreds will use the hostname value instead",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							cfg.Aws.Description = c.String("description")
							cfg.Aws.Region = c.String("region")
							cfg.Aws.SecretName = c.String("secret-name")

							// Set all other config values to empty
							cfg.Azure.SecretName = ""
							cfg.Azure.UseMSI = false
							cfg.Azure.VaultUri = ""

							cfg.HashiVault.EnvironmentTokenName = ""
							cfg.HashiVault.KeyVaultPath = ""
							cfg.HashiVault.SecretName = ""
							cfg.HashiVault.SecretPath = ""
							cfg.HashiVault.VaultUri = ""

							err := helpers.WriteConfig(configFilePath, &cfg)
							if err != nil {
								helpers.CheckError(err)
							}

							return nil
						},
					},
					{
						Name:  "azure",
						Usage: "Azure Key Vault provider configuration settings",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "secret-name",
								Usage:    "The name of the secret stored in Azure Key Vault. If omitted Terracreds will use the hostname value instead",
								Required: false,
							},
							&cli.BoolFlag{
								Name:     "use-msi",
								Usage:    "A flag to indicate if the Managed Identity of the Azure VM should be used for authentication",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "vault-uri",
								Usage:    "The FQDN of the Azure Key Vault resource",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							cfg.Azure.SecretName = c.String("secret-name")
							cfg.Azure.UseMSI = c.Bool("use-msi")
							cfg.Azure.VaultUri = c.String("vault-uri")

							// Set all other config values to empty
							cfg.Aws.Description = ""
							cfg.Aws.Region = ""
							cfg.Aws.SecretName = ""

							cfg.HashiVault.EnvironmentTokenName = ""
							cfg.HashiVault.KeyVaultPath = ""
							cfg.HashiVault.SecretName = ""
							cfg.HashiVault.SecretPath = ""
							cfg.HashiVault.VaultUri = ""

							err := helpers.WriteConfig(configFilePath, &cfg)
							if err != nil {
								helpers.CheckError(err)
							}

							return nil
						},
					},
					{
						Name:  "hashicorp",
						Usage: "HashiCorp Vault provider configuration settings",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "environment-token-name",
								Usage:    "The name of the environment variable that currently holds the Vault token",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "key-vault-path",
								Usage:    "The name of the Key Vault store inside of Vault",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "secret-name",
								Usage:    "The name of the secret stored inside of Vault. If omitted Terracreds will use the hostname value instead",
								Required: false,
							},
							&cli.StringFlag{
								Name:     "secret-path",
								Usage:    "The path of the secret itself inside of the vault",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "vault-uri",
								Usage:    "The URL of the Vault instance including its port",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							cfg.HashiVault.EnvironmentTokenName = c.String("environment-token-name")
							cfg.HashiVault.KeyVaultPath = c.String("key-vault-path")
							cfg.HashiVault.SecretName = c.String("secret-name")
							cfg.HashiVault.SecretPath = c.String("secret-path")
							cfg.HashiVault.VaultUri = c.String("vault-uri")

							// Set all other config values to empty
							cfg.Aws.Description = ""
							cfg.Aws.Region = ""
							cfg.Aws.SecretName = ""

							cfg.Azure.SecretName = ""
							cfg.Azure.UseMSI = false
							cfg.Azure.VaultUri = ""

							err := helpers.WriteConfig(configFilePath, &cfg)
							if err != nil {
								helpers.CheckError(err)
							}

							return nil
						},
					},
					{
						Name:  "logging",
						Usage: "Configure the Terracreds logging settings",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:     "enable",
								Usage:    "Enable logging",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "path",
								Aliases:  []string{"p"},
								Usage:    "The path on the file system where the log file is stored",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							cfg.Logging.Enabled = c.Bool("enable")
							cfg.Logging.Path = c.String("path")

							err := helpers.WriteConfig(configFilePath, &cfg)
							if err != nil {
								helpers.CheckError(err)
							}

							return nil
						},
					},
					{
						Name:  "view",
						Usage: "Print the current configuration to the screen",
						Action: func(c *cli.Context) error {
							bytes, err := yaml.Marshal(&cfg)
							if err != nil {
								helpers.CheckError(err)
							}

							print(string(bytes))
							return nil
						},
					},
				},
				Action: func(c *cli.Context) error {
					if c.Bool("use-local-vault") == true {
						cfg.Aws.Description = ""
						cfg.Aws.Region = ""
						cfg.Aws.SecretName = ""

						cfg.Azure.SecretName = ""
						cfg.Azure.UseMSI = false
						cfg.Azure.VaultUri = ""

						cfg.HashiVault.EnvironmentTokenName = ""
						cfg.HashiVault.KeyVaultPath = ""
						cfg.HashiVault.SecretName = ""
						cfg.HashiVault.SecretPath = ""
						cfg.HashiVault.VaultUri = ""

						err := helpers.WriteConfig(configFilePath, &cfg)
						if err != nil {
							helpers.CheckError(err)
						}
					}

					return nil
				},
			},
			{
				Name:  "create",
				Usage: "Manually create or update a credential object in the vault provider of your choice that contains the Terraform Cloud/Enterprise authorization token or another secret",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "place_holder",
						Usage:   "The name of the Terraform Cloud/Enterprise server's hostname or the name of the secret. This is also the display name of the credential object",
					},
					&cli.StringFlag{
						Name:    "secret",
						Aliases: []string{"s"},
						Value:   "",
						Usage:   "The Terraform Cloud/Enterprise API authorization token or other secret value to be securely stored in your vault provider of choice",
					},
				},
				Action: func(c *cli.Context) error {
					if len(os.Args) == 2 {
						fmt.Fprintf(color.Output, "%s: No secret name or secret was specified. Use 'terracreds create -h' to print help info\n", color.RedString("ERROR"))
						return nil
					}

					terraVault := NewTerraVault(&cfg, c.String("name"))
					name := helpers.GetSecretName(&cfg, c.String("name"))

					user, err := user.Current()
					helpers.CheckError(err)

					err = terraCreds.Create(cfg, name, c.String("secret"), user, terraVault)
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
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "place_holder",
						Usage:   "The name of the Terraform Cloud/Enterprise server's hostname or the name of the secret. This is also the display name of the credential object.",
					},
				},
				Action: func(c *cli.Context) error {
					if len(os.Args) == 2 {
						fmt.Fprintf(color.Output, "%s: No secret name was specified. Use 'terracreds delete -h' for help info\n", color.RedString("ERROR"))
						return nil
					}

					if !strings.Contains(os.Args[2], "-n") && !strings.Contains(os.Args[2], "--name") {
						msg := fmt.Sprintf("A secret name was not expected here: '%s'", os.Args[2])
						helpers.Logging(cfg, msg, "WARNING")
						fmt.Fprintf(color.Output, "%s: %s Did you mean `terracreds delete --name/-n %s'?\n", color.YellowString("WARNING"), msg, os.Args[2])
						return nil
					}

					terraVault := NewTerraVault(&cfg, c.String("name"))
					name := helpers.GetSecretName(&cfg, c.String("name"))
					method := os.Args[1]

					user, err := user.Current()
					helpers.CheckError(err)

					err = terraCreds.Delete(cfg, method, name, user, terraVault)
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
						fmt.Fprintf(color.Output, "%s: No secret name was specified. Use 'terracreds forget -h' for help info\n", color.RedString("ERROR"))
						return nil
					}

					terraVault := NewTerraVault(&cfg, os.Args[2])
					name := helpers.GetSecretName(&cfg, os.Args[2])
					method := os.Args[1]

					user, err := user.Current()
					helpers.CheckError(err)

					err = terraCreds.Delete(cfg, method, name, user, terraVault)
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
					helpers.GenerateTerraCreds(c, version)
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

						terraVault := NewTerraVault(&cfg, os.Args[2])
						name := helpers.GetSecretName(&cfg, os.Args[2])

						token, err := terraCreds.Get(cfg, name, user, terraVault)
						if err != nil {
							helpers.CheckError(err)
						}

						fmt.Println(string(token))
						return nil
					}

					msg := "A secret name was expected after the 'get' command but no argument was provided"
					helpers.Logging(cfg, msg, "ERROR")
					fmt.Fprintf(color.Output, "%s: %s\n", color.RedString("ERROR"), msg)
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "List the credentials stored in a vault using a provided set of secret names",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "secret-names",
						Aliases: []string{"s"},
						Value:   "",
						Usage:   "A comma separated list of secret names to be retrieved",
					},
					&cli.StringFlag{
						Name:    "input-file",
						Aliases: []string{"f"},
						Value:   "",
						Usage:   "The path to the file that provides the list of secrets to be retrieved. Each secret name should be on its own line",
					},
					&cli.BoolFlag{
						Name:  "export-as-tfvars",
						Value: false,
						Usage: "Exports the secret keys and values as 'TF_VARS_secret_key=secret_value' for the given operating system",
					},
				},
				Action: func(c *cli.Context) error {
					if len(os.Args) > 1 {
						terraVault := NewTerraVault(&cfg, os.Args[2])
						secretNames := strings.Split(c.String("secret-names"), ",")

						user, err := user.Current()
						if err != nil {
							helpers.CheckError(err)
						}

						list, err := terraCreds.List(c, cfg, secretNames, user, terraVault)
						if err != nil {
							helpers.CheckError(err)
						}

						for _, secret := range list {
							value := fmt.Sprintf("%s\n", secret)
							print(value)
						}

						return nil
					}

					msg := "A hostname or secret name was expected after the 'get' command but no argument was provided"
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
						fmt.Fprintf(color.Output, "%s: No hostname or secret name was specified. Use 'terracreds store -h' to print help info\n", color.RedString("ERROR"))
						return nil
					}

					terraVault := NewTerraVault(&cfg, os.Args[2])
					name := helpers.GetSecretName(&cfg, os.Args[2])

					user, err := user.Current()
					helpers.CheckError(err)

					err = terraCreds.Create(cfg, name, nil, user, terraVault)
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
