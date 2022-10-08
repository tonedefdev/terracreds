package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/fatih/color"
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/urfave/cli/v2"
)

// NewCommandCreate instantiates the command to create secrets
func (cmd *Config) NewCommandCreate() *cli.Command {
	cmdCreate := &cli.Command{
		Name:    "create",
		Aliases: []string{"update"},
		Usage:   "Manually create or update a credential object in the vault provider of your choice that contains either the Terraform Enterprise or Cloud API token or any other type of secret",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Value:   "place_holder",
				Usage:   "The name of the Terraform Automation and Collaboration Software server's hostname or the name of the secret. This is also the display name of the credential object",
			},
			&cli.StringFlag{
				Name:    "secret",
				Aliases: []string{"s", "t", "v"},
				Value:   "",
				Usage:   "The Terraform Automation and Collaboration Software API authorization token or other secret value to be securely stored in your vault provider of choice",
			},
		},
		Action: func(c *cli.Context) error {
			err := cmd.newCommandActionCreate(c)
			return err
		},
	}

	return cmdCreate
}

// newCommandActionCreate creates the secret based on the OS and type of vault
func (cmd *Config) newCommandActionCreate(c *cli.Context) error {
	if len(os.Args) == 2 {
		fmt.Fprintf(color.Output, "%s: No secret name or secret was specified. Use 'terracreds create -h' to print help info\n", color.RedString("ERROR"))
		return nil
	}

	terraVault := cmd.NewTerraVault(c.String("name"))
	name := GetSecretName(cmd.Cfg, c.String("name"))

	user, err := user.Current()
	helpers.CheckError(err)

	err = cmd.TerraCreds.Create(cmd.Cfg, name, c.String("secret"), user, terraVault)
	if err != nil {
		helpers.CheckError(err)
	}

	return nil
}
