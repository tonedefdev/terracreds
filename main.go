package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/danieljoos/wincred"
	"github.com/urfave/cli/v2"
)

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			&cli.Command{
				Name:  "new",
				Usage: "Create a new Windows Credential object that contains the Terraform Enterprise authorization token",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "hostname",
						Aliases: []string{"n"},
						Value:   "",
						Usage:   "The name of the Terraform Enterprise server's hostname. This is also the display name of the Windows Credential",
					},
					&cli.StringFlag{
						Name:    "apiToken",
						Aliases: []string{"t"},
						Value:   "",
						Usage:   "The Terraform Enterprise authorization token to store as a Windows Credential",
					},
				},
				Action: func(c *cli.Context) error {
					user, err := user.Current()
					checkError(err)

					cred := wincred.NewGenericCredential(c.String("hostname"))
					cred.CredentialBlob = []byte(c.String("apiToken"))
					cred.UserName = string(user.Username)
					err = cred.Write()

					if err == nil {
						fmt.Println("Successfully created Windows Credential")
					} else {
						log.Fatal(err)
					}
					return nil
				},
			},
			&cli.Command{
				Name:  "get",
				Usage: "Get the Windows Credential value by passing the hostname of the Terraform Enterprise server as an argument",
				Action: func(c *cli.Context) error {
					type credentialResponse struct {
						Token string `json:"token"`
					}

					user, err := user.Current()
					checkError(err)

					if len(os.Args[2]) > 0 {
						hostname := os.Args[2]
						cred, err := wincred.GetGenericCredential(hostname)

						if err == nil && cred.UserName == user.Username {
							response := &credentialResponse{
								Token: string(cred.CredentialBlob),
							}
							responseA, _ := json.Marshal(response)
							fmt.Println(string(responseA))
						} else {
							log.Fatal("You do not have permission view this Windows Credential")
						}
					} else {
						log.Fatal("The name of the Terraform Enterprise server hostname was expected")
					}
					return nil
				},
			},
			&cli.Command{
				Name:  "generate",
				Usage: "Generate the folders, credential-helpers file, and plugins required to leverage terracreds as a Terraform credential helper",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "plugins-dir",
						Aliases: []string{"p"},
						Value:   "",
						Usage:   "The path of the Terraform plugins-dir. If not specified the default is based on Terraform's default plugin directory.",
					},
				},
				Action: func(c *cli.Context) error {
					userProfile := os.Getenv("USERPROFILE")
					tfPlugins := userProfile + "\\AppData\\Roaming\\terraform.d\\plugins"

					if _, err := os.Stat(tfPlugins); err != nil {
						if os.IsNotExist(err) {
							err := os.Mkdir(tfPlugins, 0755)
							checkError(err)
						}
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	checkError(err)
}
