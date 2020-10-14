package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"runtime"

	"github.com/MakeNowJust/heredoc"
	"github.com/danieljoos/wincred"
	"github.com/urfave/cli/v2"
)

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// CopyTerraCreds will create a copy of the binary to the
// destination path.
func CopyTerraCreds(dest string) error {
	from, err := os.Open(string(os.Args[0]))
	checkError(err)
	defer from.Close()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0755)
	checkError(err)
	defer to.Close()

	_, err = io.Copy(to, from)
	checkError(err)
	fmt.Println("Successfully copied binary: " + dest)
	return nil
}

// NewDirectory checks for the existence of a directory
// if it doesn't exist it creates it and checks for errors
func NewDirectory(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(path, 0755)
			checkError(err)
			fmt.Println("Successfully created directory: " + path)
		}
	}
	return nil
}

// WriteToFile will print any string of text to a file safely by
// checking for errors and syncing at the end.
func WriteToFile(filename string, data string) error {
	file, err := os.Create(filename)
	checkError(err)
	defer file.Close()

	_, err = io.WriteString(file, data)
	checkError(err)
	fmt.Println("Successfully created file: " + filename)
	return file.Sync()
}

func main() {
	app := &cli.App{
		Name:  "terracreds",
		Usage: "a credential helper for Terraform Cloud/Enterprise that leverages the local OS credential manager for storing API tokens",
		Commands: []*cli.Command{
			&cli.Command{
				Name:  "new",
				Usage: "Create a new Windows Credential object that contains the Terraform Cloud/Enterprise authorization token",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "hostname",
						Aliases: []string{"n"},
						Value:   "",
						Usage:   "The name of the Terraform Cloud/Enterprise server's hostname. This is also the display name of the Windows Credential",
					},
					&cli.StringFlag{
						Name:    "apiToken",
						Aliases: []string{"t"},
						Value:   "",
						Usage:   "The Terraform Cloud/Enterprise authorization token to store as a Windows Credential object",
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
							log.Fatal("You do not have permission to view this Windows Credential")
						}
					}
					return nil
				},
			},
			&cli.Command{
				Name:  "generate",
				Usage: "Generate the folders and plugin binary required to leverage terracreds as a Terraform credential helper",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "plugins-dir",
						Aliases: []string{"p"},
						Value:   "",
						Usage:   "The path of the Terraform plugins-dir. If not specified the default is based on Terraform's default plugin directory.",
					},
					&cli.BoolFlag{
						Name:  "create-cli-config",
						Value: false,
						Usage: "Creates the Terraform CLI config with a terracreds credential helper block",
					},
				},
				Action: func(c *cli.Context) error {
					if runtime.GOOS == "windows" {
						userProfile := os.Getenv("USERPROFILE")
						tfUser := userProfile + "\\AppData\\Roaming\\terraform.d"
						cliConfig := tfUser + "\\terraform.rc"
						tfPlugins := tfUser + "\\plugins"
						binary := tfPlugins + "\\terraform-credentials-terracreds.exe"

						NewDirectory(tfPlugins)
						CopyTerraCreds(binary)
						if c.Bool("create-cli-config") == true {
							doc := heredoc.Doc(`
							credentials_helper "terracreds" {
								args = []
							}
							`)
							WriteToFile(cliConfig, doc)
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
//
