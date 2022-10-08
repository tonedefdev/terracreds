package main

import (
	"os"
	"runtime"

	"github.com/tonedefdev/terracreds/cmd"
	"github.com/urfave/cli/v2"
)

const (
	cfgName = "config.yaml"
	version = "2.1.2"
)

var (
	defaultReplaceString = "_"
	fileEnvVar           = os.Getenv("TC_CONFIG_PATH")
)

func main() {
	cmdConfig := &cmd.Config{
		DefaultReplaceString: defaultReplaceString,
		TerraCreds:           cmd.NewTerraCreds(runtime.GOOS),
		Version:              version,
	}

	cmdConfig.InitTerraCreds(cfgName, fileEnvVar)

	app := &cli.App{
		Name:                 "terracreds",
		EnableBashCompletion: true,
		Usage:                "a credential helper for Terraform Automation and Collaboration Software (TACOS) that leverages your vault provider of choice for securely storing API tokens or other secrets.\n\n   Visit https://github.com/tonedefdev/terracreds for more information",
		UsageText:            "Store Terraform Enterprise or Cloud API tokens by running 'terraform login' or manually store any secret you choose with 'terracreds create -n mySuperSecret -v mySuperSafePassword'",
		Version:              version,
		Commands: []*cli.Command{
			cmdConfig.NewCommandConfig(),
			cmdConfig.NewCommandCreate(),
			cmdConfig.NewCommandDelete(),
			cmdConfig.NewCommandForget(),
			cmdConfig.NewCommandGenerate(),
			cmdConfig.NewCommandGet(),
			cmdConfig.NewCommandList(),
			cmdConfig.NewCommandStore(),
		},
	}

	app.Run(os.Args)
}
