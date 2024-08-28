package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/tonedefdev/terracreds/api"
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

// Config struct defines the configuration for the Terracreds CLI application
type Config struct {
	Cfg                  *api.Config
	ConfigFile           ConfigFile
	DefaultReplaceString string
	TerraCreds           TerraCreds
	SecretNames          []string
	Version              string

	Confirm string
}

// ConfigFile defines values for the configuration file
type ConfigFile struct {
	EnvironmentValue string
	Path             string
	Name             string
}

// NewCommandConfig instantiates the config command
func (cmd *Config) NewCommandConfig() *cli.Command {
	config := &cli.Command{
		Name:  "config",
		Usage: "View or modify the Terracreds configuration file",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "use-local-vault-only",
				Usage:    "Resets configuration to only use the local operating system's credential vault. This will delete all configuration values for cloud provider vaults from the config file",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "force",
				Usage:    "Force resetting the configuration without asking for user input",
				Required: false,
				Value:    false,
			},
		},
		Subcommands: []*cli.Command{
			cmd.newCommandAws(),
			cmd.newCommandAzure(),
			cmd.newCommandGcp(),
			cmd.newCommandHashi(),
			cmd.newCommandLogging(),
			cmd.newCommandSecrets(),
			cmd.newCommandView(),
		},
		Action: func(c *cli.Context) error {
			err := cmd.newCommandActionReset(c)
			return err
		},
	}

	return config
}

// newCommandActionReset resets the configuration file to only leverage the local vault
func (cmd *Config) newCommandActionReset(c *cli.Context) error {
	if c.Bool("use-local-vault-only") {
		newCfg := api.Config{
			Logging: cmd.Cfg.Logging,
			Secrets: cmd.Cfg.Secrets,
		}

		if !c.Bool("force") {
			const verbiage = "This will reset the configuration to only use the local operating system's credential vault. Any configuration values for a cloud provider vault will be permanently lost!"
			fmt.Fprintf(color.Output, "%s: %s\n\n    Enter 'yes' to continue or press 'enter' or 'return' to cancel: ", color.YellowString("WARNING"), verbiage)
			fmt.Scanln(&cmd.Confirm)
			fmt.Print("\n")

			if cmd.Confirm == "yes" {
				err := helpers.WriteConfig(cmd.ConfigFile.Path, &newCfg)
				if err != nil {
					helpers.CheckError(err)
				}

				return err
			}
		}

		err := helpers.WriteConfig(cmd.ConfigFile.Path, &newCfg)
		if err != nil {
			helpers.CheckError(err)
		}

		return err
	}

	err := c.Command.Run(c)
	return err
}

// newCommandAws instantiates the command used to setup the AWS configuration
func (cmd *Config) newCommandAws() *cli.Command {
	awsConfig := &cli.Command{
		Name:  "aws",
		Usage: "AWS Secrets Manager provider configuration settings",
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
				Value:    "",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			err := cmd.newCommandActionAws(c)
			return err
		},
	}

	return awsConfig
}

// newCommandActionAws sets the AWS configuration and writes it to the config file
func (cmd *Config) newCommandActionAws(c *cli.Context) error {
	newCfg := api.Config{
		Aws: api.Aws{
			Description: c.String("description"),
			Region:      c.String("region"),
			SecretName:  c.String("secret-name"),
		},
		Logging: cmd.Cfg.Logging,
		Secrets: cmd.Cfg.Secrets,
	}

	err := helpers.WriteConfig(cmd.ConfigFile.Path, &newCfg)
	if err != nil {
		helpers.CheckError(err)
	}

	return err
}

// newCommandAzure instantiates the command used to setup the Azure configuration
func (cmd *Config) newCommandAzure() *cli.Command {
	azureConfig := &cli.Command{
		Name:  "azure",
		Usage: "Azure Key Vault provider configuration settings",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "secret-name",
				Usage:    "The name of the secret stored in Azure Key Vault. If omitted Terracreds will use the hostname value instead",
				Value:    "",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "subscription-id",
				Aliases:  []string{"id"},
				Usage:    "The subscription ID where the Key Vault instance has been created",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "vault-uri",
				Usage:    "The FQDN of the Azure Key Vault resource",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			err := cmd.newCommandActionAzure(c)
			return err
		},
	}

	return azureConfig
}

// newCommandActionAzure sets the Azure configuration and writes it to the config file
func (cmd *Config) newCommandActionAzure(c *cli.Context) error {
	newCfg := api.Config{
		Azure: api.Azure{
			SecretName:     c.String("secret-name"),
			SubscriptionId: c.String("subscription-id"),
			VaultUri:       c.String("vault-uri"),
		},
		Logging: cmd.Cfg.Logging,
		Secrets: cmd.Cfg.Secrets,
	}

	err := helpers.WriteConfig(cmd.ConfigFile.Path, &newCfg)
	if err != nil {
		helpers.CheckError(err)
	}

	return err
}

// newCommandGcp instantiates the command used to setup the GCP configuration
func (cmd *Config) newCommandGcp() *cli.Command {
	gcpConfig := &cli.Command{
		Name:  "gcp",
		Usage: "Google Cloud Provider Secrets Manager configuration settings",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project-id",
				Usage:    "The name of the GCP project where the Secrets Manager has been created",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "secret-id",
				Usage:    "The name of the secret identifier in GCP Secrets Manager. If omitted Terracreds will use the hostname value instead",
				Value:    "",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			err := cmd.newCommandActionGcp(c)
			return err
		},
	}

	return gcpConfig
}

// newCommandActionGcp sets the GCP configuration and writes it to the config file
func (cmd *Config) newCommandActionGcp(c *cli.Context) error {
	newCfg := api.Config{
		GCP: api.GCP{
			ProjectId: c.String("project-id"),
			SecretId:  c.String("secret-id"),
		},
		Logging: cmd.Cfg.Logging,
		Secrets: cmd.Cfg.Secrets,
	}

	err := helpers.WriteConfig(cmd.ConfigFile.Path, &newCfg)
	if err != nil {
		helpers.CheckError(err)
	}

	return err
}

// newCommandHashi instantiates the command to setup the Hashi Vault configuration
func (cmd *Config) newCommandHashi() *cli.Command {
	hashiConfig := &cli.Command{
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
				Value:    "",
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
			err := cmd.newCommandActionHashi(c)
			return err
		},
	}

	return hashiConfig
}

// newCommandActionHashi sets the Hashi Vault configuration and writes it to the config file
func (cmd *Config) newCommandActionHashi(c *cli.Context) error {
	newCfg := api.Config{
		HashiVault: api.HCVault{
			EnvironmentTokenName: c.String("environment-token-name"),
			KeyVaultPath:         c.String("key-vault-path"),
			SecretName:           c.String("secret-name"),
			SecretPath:           c.String("secret-path"),
			VaultUri:             c.String("vault-uri"),
		},
		Logging: cmd.Cfg.Logging,
		Secrets: cmd.Cfg.Secrets,
	}

	err := helpers.WriteConfig(cmd.ConfigFile.Path, &newCfg)
	if err != nil {
		helpers.CheckError(err)
	}

	return err
}

// newCommandLogging instantiates the command to manage the Terracreds logging configuration
func (cmd *Config) newCommandLogging() *cli.Command {
	loggingConfig := &cli.Command{
		Name:  "logging",
		Usage: "Configure the Terracreds logging settings",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "enabled",
				Usage:    "Enable logging",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "The path on the file system where the log file is stored",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			err := cmd.newCommandActionLogging(c)
			return err
		},
	}

	return loggingConfig
}

// newCommandActionLogging sets the logging configuration and writes it to file
func (cmd *Config) newCommandActionLogging(c *cli.Context) error {
	if c.Bool("enabled") {
		cmd.Cfg.Logging.Enabled = c.Bool("enabled")
	}

	if c.String("path") != "" {
		cmd.Cfg.Logging.Path = c.String("path")
	}

	err := helpers.WriteConfig(cmd.ConfigFile.Path, cmd.Cfg)
	if err != nil {
		helpers.CheckError(err)
	}

	return err
}

// newCommandSecrets instantiates the command to manage secrets in the Terracreds configuration file
func (cmd *Config) newCommandSecrets() *cli.Command {
	secretsConfig := &cli.Command{
		Name:  "secrets",
		Usage: "Add a list of secret names to the configuration file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "secret-list",
				Aliases:  []string{"l"},
				Usage:    "Add a comma separated list of secret names to be stored in the configuration file to use with the 'list' command",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			err := cmd.newCommandActionSecrets(c)
			return err
		},
	}

	return secretsConfig
}

// newCommandActionSecrets writes the list of secret names to the configuration file
func (cmd *Config) newCommandActionSecrets(c *cli.Context) error {
	secretValues := strings.Split(c.String("secret-list"), ",")
	cmd.Cfg.Secrets = secretValues

	err := helpers.WriteConfig(cmd.ConfigFile.Path, cmd.Cfg)
	if err != nil {
		helpers.CheckError(err)
	}

	return err
}

// newCommandView instantiates the command to view the configuration file
func (cmd *Config) newCommandView() *cli.Command {
	viewConfig := &cli.Command{
		Name:  "view",
		Usage: "Print the current configuration to the screen",
		Action: func(c *cli.Context) error {
			err := cmd.newCommandActionView(c)
			return err
		},
	}

	return viewConfig
}

// newCommandActionView reads the configuration file and prints it to the screen
func (cmd *Config) newCommandActionView(c *cli.Context) error {
	bytes, err := yaml.Marshal(&cmd.Cfg)
	if err != nil {
		helpers.CheckError(err)
	}

	print(string(bytes))
	return err
}
