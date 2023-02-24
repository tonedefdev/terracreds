package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/tonedefdev/terracreds/pkg/errors"
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/urfave/cli/v2"
)

// NewCommandList instantiates the command used to list secrets from the vault
func (cmd *Config) NewCommandList() *cli.Command {
	cmdList := &cli.Command{
		Name:  "list",
		Usage: "List the credentials stored in a vault using a provided set of secret names",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "secret-names",
				Aliases:  []string{"s", "l"},
				Value:    "",
				Usage:    "A comma separated list of secret names to be retrieved",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "from-config",
				Aliases:  []string{"f"},
				Value:    false,
				Usage:    "Get the secrets from the 'secrets' list in the configuration file",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "as-tfvars",
				Value:    false,
				Usage:    "Prints the secret keys and values as 'TF_VARS_secret_key=secret_value'",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "as-json",
				Value:    false,
				Usage:    "Prints the secret keys and values as a JSON string",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "override-replace-string",
				Value:    "",
				Usage:    "When running '--as-tfvars' the default is to replace any dashes [-] in the secret name with underscores [_]. This flag overrides that behavior and will instead replace dashes with this value",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			err := cmd.newCommandActionList(c)
			return err
		},
	}

	return cmdList
}

// newCommandActionList returns the secret names from the vault either as a string, TF_VARs, or JSON
func (cmd *Config) newCommandActionList(c *cli.Context) error {
	if len(os.Args) == 2 {
		err := &errors.CustomError{
			Message: "No list command was specified. Use 'terracreds list -h' to print help info",
			Level:   "ERROR",
		}

		helpers.Logging(cmd.Cfg, err.Message, err.Level)
		return err
	}

	if len(os.Args) > 1 {
		terraVault := cmd.NewTerraVault(os.Args[2])

		if len(cmd.Cfg.Secrets) > 0 {
			cmd.SecretNames = cmd.Cfg.Secrets
		}

		if c.String("secret-names") != "" {
			cmd.SecretNames = strings.Split(c.String("secret-names"), ",")
		}

		if len(cmd.Cfg.Secrets) < 1 && c.String("secret-names") == "" {
			err := &errors.CustomError{
				Message: "A list of secrets must be provided. Use '--secret-names' and pass it a comma separated list of secrets, or setup the 'secrets' block in the terracreds config file to use this command",
				Level:   "ERROR",
			}

			helpers.Logging(cmd.Cfg, err.Message, err.Level)
			return err
		}

		user, err := user.Current()
		if err != nil {
			helpers.CheckError(err)
		}

		list, err := cmd.TerraCreds.List(c, cmd.Cfg, cmd.SecretNames, user, terraVault)
		if err != nil {
			helpers.CheckError(err)
		}

		if c.Bool("as-json") {
			body := make(map[string]string, len(cmd.SecretNames))
			for i, name := range cmd.SecretNames {
				body[name] = list[i]
			}

			json, err := json.Marshal(body)
			if err != nil {
				helpers.CheckError(err)
			}

			fmt.Println(string(json))
			return nil
		}

		if c.Bool("as-tfvars") {
			for i, name := range cmd.SecretNames {
				if c.String("override-replace-string") != "" {
					cmd.DefaultReplaceString = c.String("override-replace-string")
				}

				formatSecretName := strings.Replace(name, "-", cmd.DefaultReplaceString, -1)
				fmt.Printf("TF_VAR_%s=%s\n", formatSecretName, list[i])
			}

			return nil
		}

		for _, secret := range list {
			value := fmt.Sprintf("%s\n", secret)
			fmt.Print(value)
		}

		return err
	}

	return nil
}
