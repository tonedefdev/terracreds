package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/fatih/color"
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/urfave/cli/v2"
)

// NewCommandForget creates the command to delete secrets called by 'terraform logout'
func (cmd *Config) NewCommandForget() *cli.Command {
	cmdForget := &cli.Command{
		Name:  "forget",
		Usage: "(Terraform Only) Forget a stored credential in your vault when 'terraform logout' has been called",
		Action: func(c *cli.Context) error {
			err := cmd.newCommandActionForget(c)
			return err
		},
	}
	return cmdForget
}

// newCommandActionForget deletes the requested secret in the vault when called by 'terraform logout'
func (cmd *Config) newCommandActionForget(c *cli.Context) error {
	if len(os.Args) == 2 {
		fmt.Fprintf(color.Output, "%s: No secret name or secret was specified. Use 'terracreds create -h' to print help info\n", color.RedString("ERROR"))
		return nil
	}

	terraVault := cmd.NewTerraVault(os.Args[2])
	name := GetSecretName(cmd.Cfg, os.Args[2])

	user, err := user.Current()
	helpers.CheckError(err)

	err = cmd.TerraCreds.Delete(cmd.Cfg, "delete", name, user, terraVault)
	if err != nil {
		helpers.CheckError(err)
	}

	return err
}
