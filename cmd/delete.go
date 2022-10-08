package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/fatih/color"
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/urfave/cli/v2"
)

// NewCommandDelete instantiates the command to delete secrets
func (cmd *Config) NewCommandDelete() *cli.Command {
	cmdDelete := &cli.Command{
		Name:  "delete",
		Usage: "Delete a stored credential in the vault",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Value:   "place_holder",
				Usage:   "The name of the Terraform Automation and Collaboration Software server's hostname or the name of the secret. This is also the display name of the credential object",
			},
		},
		Action: func(c *cli.Context) error {
			err := cmd.newCommandActionDelete(c)
			return err
		},
	}

	return cmdDelete
}

// newCommandActionDelete deletes the secret based on the type of vault
func (cmd *Config) newCommandActionDelete(c *cli.Context) error {
	if len(os.Args) == 2 {
		fmt.Fprintf(color.Output, "%s: No secret name was specified. Use 'terracreds delete -h' for help info\n", color.RedString("ERROR"))
		return nil
	}

	if !strings.Contains(os.Args[2], "-n") && !strings.Contains(os.Args[2], "--name") {
		msg := fmt.Sprintf("A secret name was not expected here: '%s'", os.Args[2])
		helpers.Logging(cmd.Cfg, msg, "WARNING")
		fmt.Fprintf(color.Output, "%s: %s Did you mean `terracreds delete --name/-n %s'?\n", color.YellowString("WARNING"), msg, os.Args[2])
		return nil
	}

	terraVault := cmd.NewTerraVault(c.String("name"))
	name := GetSecretName(cmd.Cfg, c.String("name"))
	method := os.Args[1]

	user, err := user.Current()
	helpers.CheckError(err)

	err = cmd.TerraCreds.Delete(cmd.Cfg, method, name, user, terraVault)
	if err != nil {
		helpers.CheckError(err)
	}

	return err
}
