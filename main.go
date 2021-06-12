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

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
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
func CreateCredential(c *cli.Context, hostname string, token interface{}, cfg Config) {
	user, err := user.Current()
	var apiToken CredentialResponse
	CheckError(err)

	if runtime.GOOS == "windows" {
		var method string
		_, err := wincred.GetGenericCredential(hostname)
		if err != nil {
			method = "Created"
		} else {
			method = "Updated"
		}

		cred := wincred.NewGenericCredential(hostname)
		if token == nil {
			err = json.NewDecoder(os.Stdin).Decode(&apiToken)
			if err != nil {
				fmt.Print(err.Error())
			}
			cred.CredentialBlob = []byte(apiToken.Token)
		} else {
			str := fmt.Sprintf("%v", token)
			cred.CredentialBlob = []byte(str)
		}

		cred.UserName = string(user.Username)
		err = cred.Write()

		if err == nil {
			msg := fmt.Sprintf("- %s the credential object %s", strings.ToLower(method), hostname)
			Logging(cfg, msg, "SUCCESS")

			if token != nil {
				fmt.Fprintf(color.Output, "%s: %s the credential object '%s'\n", color.GreenString("SUCCESS"), method, hostname)
			}
		} else {
			Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")

			if token != nil {
				fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
			}
		}
	}

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		var method string
		_, err := keyring.Get(hostname, string(user.Username))
		if err != nil {
			method = "Created"
		} else {
			method = "Updated"
		}

		if token == nil {
			err = json.NewDecoder(os.Stdin).Decode(&apiToken)
			if err != nil {
				fmt.Print(err.Error())
			}
			err = keyring.Set(hostname, string(user.Username), apiToken.Token)
		} else {
			str := fmt.Sprintf("%v", token)
			err = keyring.Set(hostname, string(user.Username), str)
		}

		if err == nil {
			msg := fmt.Sprintf("- %s the credential object %s", strings.ToLower(method), hostname)
			Logging(cfg, msg, "SUCCESS")

			if token != nil {
				fmt.Fprintf(color.Output, "%s: %s the credential object '%s'\n", color.GreenString("SUCCESS"), method, hostname)
			}
		} else {
			Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")

			if token != nil {
				fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
			}
		}
	}
}

// Logging forms the path and writes to log if enabled
func Logging(cfg Config, msg string, level string) {
	if cfg.Logging.Enabled == true {
		logPath := cfg.Logging.Path + "terracreds.log"
		WriteToLog(logPath, msg, LogLevel(level))
	}
}

// LogLevel returns properly formatted log level string
func LogLevel(level string) string {
	switch level {
	case "INFO":
		return "INFO: "
	case "ERROR":
		return "ERROR: "
	case "SUCCESS":
		return "SUCCESS: "
	case "WARNING":
		return "WARNING: "
	default:
		return ""
	}
}

// DeleteCredential removes the identified credential
// from the vault
func DeleteCredential(c *cli.Context, cfg Config, hostname string, command string) {
	user, err := user.Current()
	CheckError(err)

	if runtime.GOOS == "windows" {
		cred, err := wincred.GetGenericCredential(hostname)
		if err == nil && cred.UserName == user.Username {
			cred.Delete()

			msg := fmt.Sprintf("- the credential object '%s' has been removed", hostname)
			Logging(cfg, msg, "INFO")

			if command == "delete" {
				msg := fmt.Sprintf("The credential object '%s' has been removed", hostname)
				fmt.Fprintf(color.Output, "%s: %s\n", color.GreenString("SUCCESS"), msg)
			}
		} else {
			Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")

			if command == "delete" {
				fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
			}
		}
	}

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		err := keyring.Delete(hostname, string(user.Username))
		if err == nil {
			msg := fmt.Sprintf("- the credential object '%s' has been removed", hostname)
			Logging(cfg, msg, "INFO")

			if command == "delete" {
				msg := fmt.Sprintf("The credential object '%s' has been removed", hostname)
				fmt.Fprintf(color.Output, "%s: %s\n", color.GreenString("SUCCESS"), msg)
			}
		} else {
			Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")

			if command == "delete" {
				fmt.Fprintf(color.Output, "%s: You do not have permission to modify this credential\n", color.RedString("ERROR"))
			}
		}
	}
}

// GenerateTerracreds creates the binary to use this package as a credential helper
// and optionally the terraform.rc file
func GenerateTerracreds(c *cli.Context) {
	var cliConfig string
	var tfPlugins string
	var binary string

	if runtime.GOOS == "windows" {
		userProfile := os.Getenv("USERPROFILE")
		cliConfig = userProfile + "\\AppData\\Roaming\\terraform.rc"
		tfPlugins = userProfile + "\\AppData\\Roaming\\terraform.d\\plugins"
		binary = tfPlugins + "\\terraform-credentials-terracreds.exe"
	}

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		userProfile := os.Getenv("HOME")
		cliConfig = userProfile + "/.terraform.d/.terraformrc"
		tfPlugins = userProfile + "/.terraform.d/plugins"
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
	CheckError(err)

	if runtime.GOOS == "windows" {
		if cfg.Logging.Enabled == true {
			msg := fmt.Sprintf("- terraform server: %s", hostname)
			Logging(cfg, msg, "INFO")
			msg = fmt.Sprintf("- user requesting access: %s", string(user.Username))
			Logging(cfg, msg, "INFO")
		}

		cred, err := wincred.GetGenericCredential(hostname)
		if err == nil && cred.UserName == user.Username {
			response := &CredentialResponse{
				Token: string(cred.CredentialBlob),
			}
			responseA, _ := json.Marshal(response)
			fmt.Println(string(responseA))

			if cfg.Logging.Enabled == true {
				msg := fmt.Sprintf("- token was retrieved for: %s", hostname)
				Logging(cfg, msg, "INFO")
			}
		} else {
			if cfg.Logging.Enabled == true {
				Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")
			}
			fmt.Fprintf(color.Output, "%s: You do not have permission to view this credential\n", color.RedString("ERROR"))
		}
	}

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		if cfg.Logging.Enabled == true {
			msg := fmt.Sprintf("- terraform server: %s", hostname)
			Logging(cfg, msg, "INFO")
			msg = fmt.Sprintf("- user requesting access: %s", string(user.Username))
			Logging(cfg, msg, "INFO")
		}

		secret, err := keyring.Get(hostname, string(user.Username))
		if err != nil {
			if cfg.Logging.Enabled == true {
				Logging(cfg, fmt.Sprintf("- %s", err), "ERROR")
			}
			fmt.Fprintf(color.Output, "%s: You do not have permission to view this credential\n", color.RedString("ERROR"))
		} else {
			response := &CredentialResponse{
				Token: secret,
			}
			responseA, _ := json.Marshal(response)
			fmt.Println(string(responseA))

			if cfg.Logging.Enabled == true {
				msg := fmt.Sprintf("- token was retrieved for: %s", hostname)
				Logging(cfg, msg, "INFO")
			}
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
	version := "1.1.1"
	LoadConfig(&cfg)
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
						CreateCredential(c, c.String("hostname"), c.String("apiToken"), cfg)
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
							Logging(cfg, msg, "WARNING")
						}
						fmt.Fprintf(color.Output, "%s: %s Did you mean `terracreds delete --hostname/-n %s'?\n", color.YellowString("WARNING"), msg, os.Args[2])
					} else {
						hostname := c.String("hostname")
						command := os.Args[1]
						DeleteCredential(c, cfg, hostname, command)
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
						hostname := os.Args[2]
						command := os.Args[1]
						DeleteCredential(c, cfg, hostname, command)
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
					GenerateTerracreds(c)
					return nil
				},
			},
			{
				Name:  "get",
				Usage: "Get the credential object value by passing the hostname of the Terraform Cloud/Enterprise server as an argument. The credential is returned as a JSON object and formatted for consumption by Terraform",
				Action: func(c *cli.Context) error {
					if len(os.Args) > 2 {
						hostname := os.Args[2]
						GetCredential(c, cfg, hostname)
					} else {
						msg := "- hostname was expected after the 'get' command but no argument was provided"
						Logging(cfg, msg, "ERROR")
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
						hostname := os.Args[2]
						CreateCredential(c, hostname, nil, cfg)
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	CheckError(err)
}
