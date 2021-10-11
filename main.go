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

// Terracreds interface implements these methods for a credential's lifecycle
type Terracreds interface {
	// Create or store an API token in a vault
	Create(cfg api.Config, hostname string, token interface{}, user *user.User)
	// Delete or forget an API token in a vault
	Delete(cfg api.Config, command string, hostname string, user *user.User)
	// Get or retrieve an API token in a vault
	Get(cfg api.Config, hostname string, user *user.User) ([]byte, error)
}

// returnProvider returns the correct struct for the specific operating system
func returnProvider(os string) Terracreds {
	switch os {
	case "darwin":
		return platform.Mac{}
	case "linux":
		return platform.Linux{}
	case "windows":
		return platform.Windows{}
	default:
		return nil
	}
}

func main() {
	var cfg api.Config
	version := "1.1.2"
	helpers.LoadConfig(&cfg)

	provider := returnProvider(runtime.GOOS)
	if provider == nil {
		fmt.Fprintf(color.Output, "%s: Terracreds cannot run on this platform: '%s'\n", color.RedString("ERROR"), runtime.GOOS)
	}

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
						return nil
					}

					user, err := user.Current()
					helpers.CheckError(err)
					Terracreds.Create(provider, cfg, c.String("hostname"), c.String("apiToken"), user)
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
						return nil
					}

					if !strings.Contains(os.Args[2], "-n") && !strings.Contains(os.Args[2], "--hostname") {
						msg := fmt.Sprintf("A hostname was not expected here: '%s'", os.Args[2])
						helpers.Logging(cfg, msg, "WARNING")
						fmt.Fprintf(color.Output, "%s: %s Did you mean `terracreds delete --hostname/-n %s'?\n", color.YellowString("WARNING"), msg, os.Args[2])
						return nil
					}

					user, err := user.Current()
					helpers.CheckError(err)
					Terracreds.Delete(provider, cfg, os.Args[1], c.String("hostname"), user)
					return nil
				},
			},
			{
				Name:  "forget",
				Usage: "(Terraform Only) Forget a stored credential when 'terraform logout' has been called",
				Action: func(c *cli.Context) error {
					if len(os.Args) == 2 {
						fmt.Fprintf(color.Output, "%s: No hostname was specified. Use 'terracreds forget -h' for help info\n", color.RedString("ERROR"))
						return nil
					}

					user, err := user.Current()
					helpers.CheckError(err)
					Terracreds.Delete(provider, cfg, os.Args[1], os.Args[2], user)
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

						token, err := Terracreds.Get(provider, cfg, os.Args[2], user)
						if err != nil {
							helpers.CheckError(err)
						}
						fmt.Println(string(token))
						return nil
					}

					msg := "A hostname was expected after the 'get' command but no argument was provided"
					helpers.Logging(cfg, msg, "ERROR")
					fmt.Fprintf(color.Output, "%s: %s\n", color.RedString("ERROR"), msg)
					return nil
				},
			},
			{
				Name:  "store",
				Usage: "(Terraform Only) Store or update a credential object in the local operating sytem's credential manager when 'terraform login' has been called",
				Action: func(c *cli.Context) error {
					if len(os.Args) == 2 {
						fmt.Fprintf(color.Output, "%s: No hostname or token was specified. Use 'terracreds store -h' to print help info\n", color.RedString("ERROR"))
						return nil
					}

					user, err := user.Current()
					helpers.CheckError(err)
					Terracreds.Create(provider, cfg, os.Args[2], nil, user)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	helpers.CheckError(err)
}
