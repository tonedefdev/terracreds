package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/fatih/color"
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/urfave/cli/v2"
)

// NewCommandStore instantiates the command used to create credentials when 'terraform login' is called
func (cmd *Config) NewCommandStore() *cli.Command {
	cmdStore := &cli.Command{
		Name:  "store",
		Usage: "(Terraform Only) Store or update a Terraform Enterprise or Cloud API token in your vault provider of choice when 'terraform login' has been called",
		Action: func(c *cli.Context) error {
			err := cmd.newCommandActionStore(c)
			return err
		},
	}

	return cmdStore
}

// newCommandActionStore creates the secret in the vault when 'terraform login' is called
func (cmd *Config) newCommandActionStore(c *cli.Context) error {
	if len(os.Args) == 2 {
		fmt.Fprintf(color.Output, "%s: No hostname was specified. Use 'terracreds store -h' to print help info\n", color.RedString("ERROR"))
		return nil
	}

	terraVault := cmd.NewTerraVault(os.Args[2])
	name := GetSecretName(cmd.Cfg, os.Args[2])

	user, err := user.Current()
	helpers.CheckError(err)

	err = cmd.TerraCreds.Create(cmd.Cfg, name, nil, user, terraVault)
	if err != nil {
		helpers.CheckError(err)
	}

	return err
}
