package main

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	api "github.com/tonedefdev/terracreds/api"
	helpers "github.com/tonedefdev/terracreds/pkg/helpers"
	platform "github.com/tonedefdev/terracreds/pkg/platform"
)

func main() {
	var cfg api.Config
	version := "1.1.2"
	helpers.LoadConfig(&cfg)
	app := &cli.App{
		Name:      "terracreds",
		Usage:     "a credential helper for Terraform Cloud/Enterprise that leverages the local operating system's credential manager for securely storing your API tokens.\n\n   Visit https://github.com/tonedefdev/terracreds for more information",
		UsageText: "Directly store credentials from Terraform using 'terraform login' or manually store them using 'terracreds create -n app.terraform.io -t myAPItoken'",
		Version:   version,
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Manually create or update a credential object in the local operating sytem's credential manager that contains the Terraform Cloud/Enterprise authorization token",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "hostname",
						Aliases: []string{"n"},
						Value:   "",
						Usage:   "The name of the Terraform Cloud/Enterprise server's hostname. This is also the display name of the credential object",
					},
					&cli.StringFlag{
						Name:    "apiToken",
						Aliases: []string{"t"},
						Value:   "",
						Usage:   "The Terraform Cloud/Enterprise API authorization token to be securely stored in the local operating system's credential manager",
					},
				},
				Action: func(c *cli.Context) error {
					if len(os.Args) == 2 {
						fmt.Fprintf(color.Output, "%s: No hostname or token was specified. Use 'terracreds create -h' to print help info\n", color.RedString("ERROR"))
					} else {
						user, err := user.Current()
						helpers.CheckError(err)

						if runtime.GOOS == "windows" {
							os := platform.Windows{
								ApiToken: api.CredentialResponse{},
								Config:   cfg,
								Context:  c,
								Hostname: c.String("hostname"),
								Token:    c.String("apiToken"),
								User:     user,
							}

							api.Terracreds.Create(os)
						}

						if runtime.GOOS == "dawrin" {
							os := platform.Mac{
								ApiToken: api.CredentialResponse{},
								Config:   cfg,
								Context:  c,
								Hostname: c.String("hostname"),
								Token:    c.String("apiToken"),
								User:     user,
							}

							api.Terracreds.Create(os)
						}

						if runtime.GOOS == "linux" {
							os := platform.Linux{
								ApiToken: api.CredentialResponse{},
								Config:   cfg,
								Context:  c,
								Hostname: c.String("hostname"),
								Token:    c.String("apiToken"),
								User:     user,
							}

							api.Terracreds.Create(os)
						}
					}
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "Delete a stored credential in the local operating system's credential manager",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "hostname",
						Aliases: []string{"n"},
						Value:   "",
						Usage:   "The name of the Terraform Cloud/Enterprise server's hostname. This is also the display name of the credential object",
					},
				},
				Action: func(c *cli.Context) error {
					if len(os.Args) == 2 {
						fmt.Fprintf(color.Output, "%s: No hostname was specified. Use 'terracreds delete -h' for help info\n", color.RedString("ERROR"))
					} else if !strings.Contains(os.Args[2], "-n") && !strings.Contains(os.Args[2], "--hostname") {
						msg := fmt.Sprintf("A hostname was not expected here: %s", os.Args[2])
						if cfg.Logging.Enabled == true {
							helpers.Logging(cfg, msg, "WARNING")
						}
						fmt.Fprintf(color.Output, "%s: %s Did you mean `terracreds delete --hostname/-n %s'?\n", color.YellowString("WARNING"), msg, os.Args[2])
					} else {
						user, err := user.Current()
						helpers.CheckError(err)

						if runtime.GOOS == "windows" {
							os := platform.Windows{
								ApiToken: api.CredentialResponse{},
								Command:  os.Args[1],
								Config:   cfg,
								Context:  c,
								Hostname: c.String("hostname"),
								Token:    c.String("apiToken"),
								User:     user,
							}

							api.Terracreds.Delete(os)
						}

						if runtime.GOOS == "dawrin" {
							os := platform.Mac{
								ApiToken: api.CredentialResponse{},
								Command:  os.Args[1],
								Config:   cfg,
								Context:  c,
								Hostname: c.String("hostname"),
								Token:    c.String("apiToken"),
								User:     user,
							}

							api.Terracreds.Delete(os)
						}

						if runtime.GOOS == "linux" {
							os := platform.Linux{
								ApiToken: api.CredentialResponse{},
								Command:  os.Args[1],
								Config:   cfg,
								Context:  c,
								Hostname: c.String("hostname"),
								Token:    c.String("apiToken"),
								User:     user,
							}

							api.Terracreds.Delete(os)
						}
					}
					return nil
				},
			},
			{
				Name:  "forget",
				Usage: "(Terraform Only) Forget a stored credential when 'terraform logout' has been called",
				Action: func(c *cli.Context) error {
					if len(os.Args) == 2 {
						fmt.Fprintf(color.Output, "%s: No hostname was specified. Use 'terracreds forget -h' for help info\n", color.RedString("ERROR"))
					} else {
						user, err := user.Current()
						helpers.CheckError(err)

						if runtime.GOOS == "windows" {
							os := platform.Windows{
								ApiToken: api.CredentialResponse{},
								Command:  os.Args[1],
								Config:   cfg,
								Context:  c,
								Hostname: os.Args[2],
								Token:    c.String("apiToken"),
								User:     user,
							}

							api.Terracreds.Delete(os)
						}

						if runtime.GOOS == "dawrin" {
							os := platform.Mac{
								ApiToken: api.CredentialResponse{},
								Command:  os.Args[1],
								Config:   cfg,
								Context:  c,
								Hostname: os.Args[2],
								Token:    c.String("apiToken"),
								User:     user,
							}

							api.Terracreds.Delete(os)
						}

						if runtime.GOOS == "linux" {
							os := platform.Linux{
								ApiToken: api.CredentialResponse{},
								Command:  os.Args[1],
								Config:   cfg,
								Context:  c,
								Hostname: os.Args[2],
								Token:    c.String("apiToken"),
								User:     user,
							}

							api.Terracreds.Delete(os)
						}
					}
					return nil
				},
			},
			{
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
					helpers.GenerateTerracreds(c)
					return nil
				},
			},
			{
				Name:  "get",
				Usage: "Get the credential object value by passing the hostname of the Terraform Cloud/Enterprise server as an argument. The credential is returned as a JSON object and formatted for consumption by Terraform",
				Action: func(c *cli.Context) error {
					if len(os.Args) > 2 {
						user, err := user.Current()
						helpers.CheckError(err)

						if runtime.GOOS == "windows" {
							os := platform.Windows{
								ApiToken: api.CredentialResponse{},
								Command:  os.Args[1],
								Config:   cfg,
								Context:  c,
								Hostname: os.Args[2],
								Token:    c.String("apiToken"),
								User:     user,
							}

							api.Terracreds.Get(os)
						}

						if runtime.GOOS == "dawrin" {
							os := platform.Mac{
								ApiToken: api.CredentialResponse{},
								Command:  os.Args[1],
								Config:   cfg,
								Context:  c,
								Hostname: os.Args[2],
								Token:    c.String("apiToken"),
								User:     user,
							}

							api.Terracreds.Get(os)
						}

						if runtime.GOOS == "linux" {
							os := platform.Linux{
								ApiToken: api.CredentialResponse{},
								Command:  os.Args[1],
								Config:   cfg,
								Context:  c,
								Hostname: os.Args[2],
								Token:    c.String("apiToken"),
								User:     user,
							}

							api.Terracreds.Get(os)
						}
					} else {
						msg := "- hostname was expected after the 'get' command but no argument was provided"
						helpers.Logging(cfg, msg, "ERROR")
						fmt.Fprintf(color.Output, "%s: %s\n", color.RedString("ERROR"), msg)
					}
					return nil
				},
			},
			{
				Name:  "store",
				Usage: "(Terraform Only) Store or update a credential object in the local operating sytem's credential manager when 'terraform login' has been called",
				Action: func(c *cli.Context) error {
					if len(os.Args) == 2 {
						fmt.Fprintf(color.Output, "%s: No hostname or token was specified. Use 'terracreds store -h' to print help info\n", color.RedString("ERROR"))
					} else {
						user, err := user.Current()
						helpers.CheckError(err)

						if runtime.GOOS == "windows" {
							os := platform.Windows{
								ApiToken: api.CredentialResponse{},
								Command:  os.Args[1],
								Config:   cfg,
								Context:  c,
								Hostname: os.Args[2],
								Token:    nil,
								User:     user,
							}

							api.Terracreds.Create(os)
						}

						if runtime.GOOS == "dawrin" {
							os := platform.Mac{
								ApiToken: api.CredentialResponse{},
								Command:  os.Args[1],
								Config:   cfg,
								Context:  c,
								Hostname: os.Args[2],
								Token:    nil,
								User:     user,
							}

							api.Terracreds.Create(os)
						}

						if runtime.GOOS == "linux" {
							os := platform.Linux{
								ApiToken: api.CredentialResponse{},
								Command:  os.Args[1],
								Config:   cfg,
								Context:  c,
								Hostname: os.Args[2],
								Token:    nil,
								User:     user,
							}

							api.Terracreds.Create(os)
						}
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	helpers.CheckError(err)
}
