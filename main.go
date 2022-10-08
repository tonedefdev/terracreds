package main

import (
	"os"
	"runtime"

	"github.com/tonedefdev/terracreds/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	terracreds := &cmd.Config{
		ConfigFileName:       "config.yaml",
		ConfigFileEnvValue:   os.Getenv("TC_CONFIG_PATH"),
		DefaultReplaceString: "_",
		TerraCreds:           cmd.NewTerraCreds(runtime.GOOS),
		Version:              "2.1.2",
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

	app.Run(os.Args)
}
