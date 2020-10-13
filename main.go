package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/danieljoos/wincred"
	"github.com/urfave/cli/v2"
)

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
					cred := wincred.NewGenericCredential(c.String("hostname"))
					cred.CredentialBlob = []byte(c.String("apiToken"))
					err := cred.Write()

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

					if len(os.Args[2]) > 0 {
						hostname := os.Args[2]
						cred, err := wincred.GetGenericCredential(hostname)
						if err == nil {
							response := &credentialResponse{
								Token: string(cred.CredentialBlob),
							}
							responseA, _ := json.Marshal(response)
							fmt.Println(string(responseA))
						} else {
							log.Fatal(err)
						}
					} else {
						log.Fatal("The name of the Terraform Enterprise server hostname was expected")
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
