package config

import (
	"github.com/tonedefdev/terracreds/api"
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/urfave/cli/v2"
)

func NewConfig() *cli.Command {
	config := &cli.Command{
		Name:  "config",
		Usage: "View or modify the Terracreds configuration file",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "use-local-vault-only",
				Usage:    "Resets configuration to only use the local operating system's credential vault. This will delete all configuration values for cloud provider vaults from the config file",
				Required: false,
			},
		},
	}

	return config
}

func NewAwsConfig(cfg api.Config, configFilePath string) *cli.Command {
	awsConfig := &cli.Command{
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
				Value:    "",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			cfg.Aws.Description = c.String("description")
			cfg.Aws.Region = c.String("region")
			cfg.Aws.SecretName = c.String("secret-name")

			// Set all other config values to empty
			cfg.Azure.SecretName = ""
			cfg.Azure.SubscriptionId = ""
			cfg.Azure.VaultUri = ""

			cfg.GCP.ProjectId = ""
			cfg.GCP.SecretId = ""

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
	}

	return awsConfig
}
