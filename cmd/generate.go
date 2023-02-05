package cmd

import (
	"github.com/tonedefdev/terracreds/pkg/helpers"
	"github.com/urfave/cli/v2"
)

// NewCommandGenerate instantiates the command used to generate the Terracreds binary and terraform.rc file
func (cmd *Config) NewCommandGenerate() *cli.Command {
	cmdGenerate := &cli.Command{
		Name:  "generate",
		Usage: "Generate the folders and plugin binary required to leverage terracreds as a Terraform credential helper",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "create-cli-config",
				Value: false,
				Usage: "Creates the Terraform CLI config with a terracreds credential helper block. This will overwrite the existing file if it already exists.",
			},
			&cli.BoolFlag{
				Name:  "force",
				Value: false,
				Usage: "Force creation of the CLI config without user input.",
			},
		},
		Action: func(c *cli.Context) error {
			err := GenerateTerraCreds(c, cmd.Version, cmd.Confirm)
			if err != nil {
				helpers.CheckError(err)
			}
			return err
		},
	}

	return cmdGenerate
}
