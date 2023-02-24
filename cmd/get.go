package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/tonedefdev/terracreds/pkg/errors"
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/urfave/cli/v2"
)

// NewCommandGet instantiates the command used to get secrets
func (cmd *Config) NewCommandGet() *cli.Command {
	cmdGet := &cli.Command{
		Name:  "get",
		Usage: "Get the credential object value by passing the server's hostname (Terraform backend default behavior) or the name of the secret as an argument. The credential is returned as a JSON object and formatted for consumption by Terraform",
		Action: func(c *cli.Context) error {
			err := cmd.newCommandActionGet()
			return err
		},
	}

	return cmdGet
}

// newCommandActionGet returns a JSON string representing the secret value stored in the vault
func (cmd *Config) newCommandActionGet() error {
	if len(os.Args) > 2 {
		user, err := user.Current()
		helpers.CheckError(err)

		terraVault := cmd.NewTerraVault(os.Args[2])
		name := GetSecretName(cmd.Cfg, os.Args[2])

		token, err := cmd.TerraCreds.Get(cmd.Cfg, name, user, terraVault)
		if err != nil {
			helpers.CheckError(err)
		}

		fmt.Println(string(token))
		return err
	}

	err := &errors.CustomError{
		Message: "A secret name was expected after the 'get' command but no argument was provided",
		Level:   "ERROR",
	}

	helpers.Logging(cmd.Cfg, err.Message, err.Level)
	return err
}
