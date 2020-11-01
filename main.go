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
	"github.com/zalando/go-keyring"
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
	fmt.Fprintf(color.Output, "%s: Copied binary '%s' to '%s'\n", color.CyanString("INFO"), string(os.Args[0]), dest)
	return nil
}

// GetBinaryPath returns the directory of the binary path
func GetBinaryPath(binary string) string {
	var path string
	if runtime.GOOS == "windows" {
		if strings.Contains(binary, "terraform-credentials-terracreds.exe") {
			path = strings.Replace(binary, "terraform-credentials-terracreds.exe", "", -1)
		} else if strings.Contains(binary, "terracreds.test.exe") {
			path = strings.Replace(binary, "terracreds.test.exe", "", -1)
		} else {
			path = strings.Replace(binary, "terracreds.exe", "", -1)
		}
	}

	if runtime.GOOS == "darwin" {
		if strings.Contains(binary, "terraform-credentials-terracreds") {
			path = strings.Replace(binary, "terraform-credentials-terracreds", "", -1)
		} else if strings.Contains(binary, "terracreds.test") {
			path = strings.Replace(binary, "terracreds.test", "", -1)
		} else {
			path = strings.Replace(binary, "terracreds", "", -1)
		}
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
			fmt.Fprintf(color.Output, "%s: Created directory '%s'\n", color.CyanString("INFO"), path)
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
	fmt.Fprintf(color.Output, "%s: Created file '%s'\n", color.CyanString("INFO"), filename)
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
	bin := GetBinaryPath(os.Args[0])
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

	bin := GetBinaryPath(os.Args[0])
	path := bin + "config.yaml"
	f, err := os.Open(string(path))
	CheckError(err)
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	CheckError(err)
	return err
}

// CreateCredential checks the current os
// then creates a credential object in its vault
func CreateCredential(c *cli.Context, hostname string, apiToken string) {
	user, err := user.Current()
	CheckError(err)

	if runtime.GOOS == "windows" {
		cred := wincred.NewGenericCredential(hostname)
		cred.CredentialBlob = []byte(apiToken)
		cred.UserName = string(user.Username)
		err = cred.Write()

		if err == nil {
			fmt.Fprintf(color.Output, "%s: Created\\updated the credential object '%s'\n", color.GreenString("SUCCESS"), hostname)
		} else {
			fmt.Fprintf(color.Output, "%s: You do not have permission to create this credential\n", color.RedString("ERROR"))
		}
	}

	if runtime.GOOS == "darwin" {
		err := keyring.Set(hostname, string(user.Username), apiToken)
		if err != nil {
			fmt.Fprintf(color.Output, "%s: Created\\updated the credential object '%s'\n", color.GreenString("SUCCESS"), hostname)
		} else {
			fmt.Fprintf(color.Output, "%s: You do not have permission to create this credential\n", color.RedString("ERROR"))
		}
	}
}

// DeleteCredential removes the identified credential
// from the vault
func DeleteCredential(c *cli.Context, cfg Config, hostname string) {
	user, err := user.Current()
	CheckError(err)

	if runtime.GOOS == "windows" {
		cred, err := wincred.GetGenericCredential(hostname)
		if err == nil && cred.UserName == user.Username {
			cred.Delete()

			msg := "The credential object '" + hostname + "' has been removed"
			fmt.Fprintf(color.Output, "%s: %s\n", color.GreenString("SUCCESS"), msg)
			if cfg.Logging.Enabled == true {
				logPath := cfg.Logging.Path + "\\terracreds.log"
				WriteToLog(logPath, msg, "INFO: ")
			}
		} else {
			fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
		}
	}

	if runtime.GOOS == "darwin" {
		err := keyring.Delete(hostname, string(user.Username))
		fmt.Println(string(user.Username))
		fmt.Println(err)
		if err == nil {
			msg := "The credential object '" + hostname + "' has been removed"
			fmt.Fprintf(color.Output, "%s: %s\n", color.GreenString("SUCCESS"), msg)
		} else {
			fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
		}
	}
}

// GenerateTerracreds creates the binary to use this package as a credential helper
// and optionally the terraform.rc file
func GenerateTerracreds(c *cli.Context, path string, tfUser string) {
	var cliConfig string
	var tfPlugins string
	var binary string

	if runtime.GOOS == "windows" {
		cliConfig = tfUser + "\\terraform.rc"
		tfPlugins = tfUser + "\\plugins"
		binary = tfPlugins + "\\terraform-credentials-terracreds.exe"
	}

	if runtime.GOOS == "darwin" {
		cliConfig = os.Getenv("HOME") + "/.terraformrc"
		tfPlugins = tfUser + "/plugins"
		binary = tfPlugins + "/terraform-credentials-terracreds"
	}

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

// GetCredential returns the stored credential as a JSON
// object as required to be consumed by Terraform Cloud/Enterprise
func GetCredential(c *cli.Context, cfg Config, hostname string) {
	user, err := user.Current()
	var logPath string
	CheckError(err)

	if runtime.GOOS == "windows" {
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
			fmt.Fprintf(color.Output, "%s: You do not have permission to view this credential\n", color.RedString("ERROR"))
		}
	}

	if runtime.GOOS == "darwin" {
		secret, err := keyring.Get(hostname, string(user.Username))
		if err != nil {
			fmt.Fprintf(color.Output, "%s: You do not have permission to view this credential\n", color.RedString("ERROR"))
		} else {
			response := &CredentialResponse{
				Token: secret,
			}
			responseA, _ := json.Marshal(response)
			fmt.Println(string(responseA))
		}
	}
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
	var logPath string
	LoadConfig(&cfg)
	app := &cli.App{
		Name:      "terracreds",
		Usage:     "a credential helper for Terraform Cloud/Enterprise that leverages the local operating system's credential manager for securely storing your API tokens.\n\n   Visit https://github.com/tonedefdev/terracreds for more information",
		UsageText: "terracreds create -n api.terraform.com -t sampleApiTokenString",
		Version:   "1.0.1",
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
					CreateCredential(c, c.String("hostname"), c.String("apiToken"))
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
					hostname := c.String("hostname")
					DeleteCredential(c, cfg, hostname)
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
					var userProfile string
					var tfUser string

					if runtime.GOOS == "windows" {
						userProfile = os.Getenv("USERPROFILE")
						tfUser = userProfile + "\\AppData\\Roaming\\terraform.d"
					}

					if runtime.GOOS == "darwin" {
						userProfile = os.Getenv("HOME")
						tfUser = userProfile + "/.terraform.d"
					}
					GenerateTerracreds(c, userProfile, tfUser)
					return nil
				},
			},
			&cli.Command{
				Name:  "get",
				Usage: "Get the credential object value by passing the hostname of the Terraform Cloud/Enterprise server as an argument. The credential is returned as a JSON object and formatted for consumption by Terraform",
				Action: func(c *cli.Context) error {
					if len(os.Args) > 2 {
						hostname := os.Args[2]
						GetCredential(c, cfg, hostname)
					} else {
						msg := "A hostname was expected after the 'get' command but no argument was provided"
						if cfg.Logging.Enabled == true {
							logPath = cfg.Logging.Path + "\\terracreds.log"
							WriteToLog(logPath, msg, "ERROR: ")
						}
						fmt.Fprintf(color.Output, "%s: %s\n", color.RedString("ERROR"), msg)
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	CheckError(err)
}
