package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"

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

// WriteToLog will create a log if it doesn't exist and then append
// messages to the log
func WriteToLog(path string, data string, level string) error {
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	checkError(err)
	defer f.Close()

	logger := log.New(f, level, log.LstdFlags)
	logger.Println(data)
	return nil
}

func main() {
	app := &cli.App{
		Name:  "terracreds",
		Usage: "a credential helper for Terraform Cloud/Enterprise that leverages the local operating system's credential manager for securely storing your API tokens",
		Commands: []*cli.Command{
			&cli.Command{
				Name:  "create",
				Usage: "Create a new credential object in the local operating sytem's credential manager that contains the Terraform Cloud/Enterprise authorization token",
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
						Usage:   "The Terraform Cloud/Enterprise authorization token to be securely stored in the local operating system's credential manager",
					},
				},
				Action: func(c *cli.Context) error {
					user, err := user.Current()
					checkError(err)

					if runtime.GOOS == "windows" {
						cred := wincred.NewGenericCredential(c.String("hostname"))
						cred.CredentialBlob = []byte(c.String("apiToken"))
						cred.UserName = string(user.Username)
						err = cred.Write()

						if err == nil {
							fmt.Println("Successfully created the credential object")
						} else {
							log.Fatal(err)
						}
					}
					return nil
				},
			},
			&cli.Command{
				Name:  "delete",
				Usage: "Delete a credential stored in the local operating system's credential manager",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "hostname",
						Aliases: []string{"n"},
						Value:   "",
						Usage:   "The name of the Terraform Cloud/Enterprise server's hostname. This is also the display name of the credential object",
					},
				},
				Action: func(c *cli.Context) error {
					user, err := user.Current()
					checkError(err)

					cred, err := wincred.GetGenericCredential(c.String("hostname"))
					if err == nil && cred.UserName == user.Username {
						cred.Delete()
						fmt.Println("The credential object '" + c.String("hostname") + "' has been removed")
					} else {
						log.Fatal("You do not have permission to access this credential object")
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
						Usage:   "The path of the Terraform plugins-dir. If not specified the default is based on Terraform's default plugin directory for the operating system.",
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
			&cli.Command{
				Name:  "get",
				Usage: "Get the credential object value by passing the hostname of the Terraform Cloud/Enterprise server as an argument",
				Action: func(c *cli.Context) error {
					type credentialResponse struct {
						Token string `json:"token"`
					}

					user, err := user.Current()
					checkError(err)

					if len(os.Args[2]) > 0 {
						hostname := os.Args[2]
						var path string
						if strings.Contains(os.Args[0], "terraform-credentials-terracreds.exe") {
							path = strings.Replace(os.Args[0], "terraform-credentials-terracreds.exe", "", -1)
						} else {
							path = strings.Replace(os.Args[0], "terracreds.exe", "", -1)
						}

						logPath := path + "\\terracreds.log"
						cred, err := wincred.GetGenericCredential(hostname)

						WriteToLog(logPath, "- terraform server: "+hostname, "INFO: ")
						WriteToLog(logPath, "- user requesting access: "+string(user.Username), "INFO: ")

						if err == nil && cred.UserName == user.Username {
							response := &credentialResponse{
								Token: string(cred.CredentialBlob),
							}
							responseA, _ := json.Marshal(response)
							fmt.Println(string(responseA))
							WriteToLog(logPath, "- token was retrieved for: "+hostname, "INFO: ")
						} else {
							WriteToLog(logPath, "- access was denied to: "+string(user.Username), "ERROR: ")
							log.Fatal("You do not have permission to view this credential")
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
