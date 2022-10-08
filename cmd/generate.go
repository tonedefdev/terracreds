package cmd

import (
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
		},
		Action: func(c *cli.Context) error {
			GenerateTerraCreds(c, cmd.Version, cmd.confirm)
			return nil
		},
	}

	return cmdGenerate
}
