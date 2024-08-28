package main

import (
	"os"

	"github.com/tonedefdev/terracreds/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	terracreds := &cmd.Config{
		DefaultReplaceString: "_",
		TerraCreds:           cmd.NewTerraCreds(),
		Version:              "2.1.6",

		ConfigFile: cmd.ConfigFile{
			EnvironmentValue: os.Getenv("TC_CONFIG_PATH"),
			Name:             "config.yaml",
		},
	}

	terracreds.InitTerraCreds()

	app := &cli.App{
		Name:                 "terracreds",
		EnableBashCompletion: true,
		Usage:                "a credential helper for Terraform Automation and Collaboration Software (TACOS) that leverages your vault provider of choice for securely storing API tokens or other secrets.\n\n   Visit https://github.com/tonedefdev/terracreds for more information",
		UsageText:            "Store Terraform Enterprise or Cloud API tokens by running 'terraform login' or manually store any secret you choose with 'terracreds create -n mySuperSecret -v mySuperSafePassword'",
		Version:              terracreds.Version,
		Commands: []*cli.Command{
			terracreds.NewCommandConfig(),
			terracreds.NewCommandCreate(),
			terracreds.NewCommandDelete(),
			terracreds.NewCommandForget(),
			terracreds.NewCommandGenerate(),
			terracreds.NewCommandGet(),
			terracreds.NewCommandList(),
			terracreds.NewCommandStore(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}
}
