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
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

// CheckError processes the error
func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// CopyTerraCreds will create a copy of the binary to the
// destination path.
func CopyTerraCreds(dest string) error {
	from, err := os.Open(string(os.Args[0]))
	CheckError(err)
	defer from.Close()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0755)
	CheckError(err)
	defer to.Close()

	_, err = io.Copy(to, from)
	CheckError(err)
	fmt.Println("Successfully copied binary: " + dest)
	return nil
}

// GetBinaryPath returns the directory of the binary path
func GetBinaryPath() string {
	var path string
	if strings.Contains(os.Args[0], "terraform-credentials-terracreds.exe") {
		path = strings.Replace(os.Args[0], "terraform-credentials-terracreds.exe", "", -1)
	} else {
		path = strings.Replace(os.Args[0], "terracreds.exe", "", -1)
	}
	return path
}

// NewDirectory checks for the existence of a directory
// if it doesn't exist it creates it and checks for errors
func NewDirectory(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(path, 0755)
			CheckError(err)
			fmt.Println("Successfully created directory: " + path)
		}
	}
	return nil
}

// WriteToFile will print any string of text to a file safely by
// checking for errors and syncing at the end.
func WriteToFile(filename string, data string) error {
	file, err := os.Create(filename)
	CheckError(err)
	defer file.Close()

	_, err = io.WriteString(file, data)
	CheckError(err)
	fmt.Println("Successfully created file: " + filename)
	return file.Sync()
}

// WriteToLog will create a log if it doesn't exist and then append
// messages to the log
func WriteToLog(path string, data string, level string) error {
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	CheckError(err)
	defer f.Close()

	logger := log.New(f, level, log.LstdFlags)
	logger.Println(data)
	return nil
}

// CreateConfigFile creates a default terracreds config file if one does not
// exist in the same path as the binary
func CreateConfigFile() error {
	bin := GetBinaryPath()
	path := bin + "config.yaml"

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			doc := heredoc.Doc(`
logging:
  enabled: false
  path:`)
			WriteToFile(path, doc)
		}
	}
	return nil
}

// LoadConfig loads the config file if it exists
// and creates the config file if doesn't exist
func LoadConfig(cfg *Config) error {
	CreateConfigFile()

	bin := GetBinaryPath()
	path := bin + "config.yaml"
	f, err := os.Open(string(path))
	CheckError(err)
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	CheckError(err)
	return err
}

// Config struct for terracreds custom configuration
type Config struct {
	Logging struct {
		Enabled bool   `yaml:"enabled"`
		Path    string `yaml:"path"`
	} `yaml:"logging"`
}

// CredentialResponse formatted for consumption by Terraform
type CredentialResponse struct {
	Token string `json:"token"`
}

func main() {
	var cfg Config
	LoadConfig(&cfg)
	app := &cli.App{
		Name:      "terracreds",
		Usage:     "a credential helper for Terraform Cloud/Enterprise that leverages the local operating system's credential manager for securely storing your API tokens.\n\n   Visit https://github.com/tonedefdev/terracreds for more information",
		UsageText: "terracreds create -n api.terraform.com -t sampleApiTokenString",
		Version:   "1.0.0",
		Commands: []*cli.Command{
			&cli.Command{
				Name:  "create",
				Usage: "Create or update a credential object in the local operating sytem's credential manager that contains the Terraform Cloud/Enterprise authorization token",
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
					user, err := user.Current()
					CheckError(err)

					if runtime.GOOS == "windows" {
						cred := wincred.NewGenericCredential(c.String("hostname"))
						cred.CredentialBlob = []byte(c.String("apiToken"))
						cred.UserName = string(user.Username)
						err = cred.Write()

						if err == nil {
							fmt.Fprintf(color.Output, "%s: Created\\updated the credential object '%s'", color.GreenString("SUCCESS"), c.String("hostname"))
						} else {
							fmt.Fprintf(color.Output, "%s: You do not have permission to view this credential", color.RedString("ERROR"))
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
					CheckError(err)

					cred, err := wincred.GetGenericCredential(c.String("hostname"))
					if err == nil && cred.UserName == user.Username {
						cred.Delete()

						msg := "The credential object '" + c.String("hostname") + "' has been removed"
						fmt.Fprintf(color.Output, "%s: %s", color.GreenString("SUCCESS"), msg)
						if cfg.Logging.Enabled == true {
							logPath := cfg.Logging.Path + "\\terracreds.log"
							WriteToLog(logPath, msg, "INFO: ")
						}
					} else {
						fmt.Fprintf(color.Output, "%s: You do not have permission to view this credential", color.RedString("ERROR"))
					}
					return nil
				},
			},
			&cli.Command{
				Name:  "generate",
				Usage: "Generate the folders and plugin binary required to leverage terracreds as a Terraform credential helper",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "create-cli-config",
						Value: false,
						Usage: "Creates the Terraform CLI config with a terracreds credential helper block. This will overwrite the existing file if it already exists",
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
							}`)
							WriteToFile(cliConfig, doc)
						}
					}
					return nil
				},
			},
			&cli.Command{
				Name:  "get",
				Usage: "Get the credential object value by passing the hostname of the Terraform Cloud/Enterprise server as an argument. The credential is returned as a JSON object and formatted for consumption by Terraform",
				Action: func(c *cli.Context) error {
					user, err := user.Current()
					CheckError(err)

					if len(os.Args) > 2 {
						var logPath string
						hostname := os.Args[2]
						if cfg.Logging.Enabled == true {
							logPath = cfg.Logging.Path + "\\terracreds.log"
							WriteToLog(logPath, "- terraform server: "+hostname, "INFO: ")
							WriteToLog(logPath, "- user requesting access: "+string(user.Username), "INFO: ")
						}
						cred, err := wincred.GetGenericCredential(hostname)
						if err == nil && cred.UserName == user.Username {
							response := &CredentialResponse{
								Token: string(cred.CredentialBlob),
							}
							responseA, _ := json.Marshal(response)
							fmt.Println(string(responseA))

							if cfg.Logging.Enabled == true {
								WriteToLog(logPath, "- token was retrieved for: "+hostname, "INFO: ")
							}
						} else {
							if cfg.Logging.Enabled == true {
								WriteToLog(logPath, "- access was denied for user: "+string(user.Username), "ERROR: ")
							}
							fmt.Fprintf(color.Output, "%s: You do not have permission to view this credential", color.RedString("ERROR"))
						}
					} else {
						var logPath string
						msg := "A hostname was expected after the 'get' command but no argument was provided"
						if cfg.Logging.Enabled == true {
							logPath = cfg.Logging.Path + "\\terracreds.log"
							WriteToLog(logPath, msg, "ERROR: ")
						}
						fmt.Fprintf(color.Output, "%s: %s", color.RedString("ERROR"), msg)
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	CheckError(err)
}
