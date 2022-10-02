package helpers

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/tonedefdev/terracreds/api"
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
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}

	fmt.Fprintf(color.Output, "%s: Copied binary '%s' to '%s'\n", color.CyanString("INFO"), string(os.Args[0]), dest)
	return err
}

// GetBinaryPath returns the directory of the binary path
func GetBinaryPath(binary string, os string) string {
	var path string

	switch os {
	case "darwin":
		if strings.Contains(binary, "terraform-credentials-terracreds") {
			path = strings.Replace(binary, "terraform-credentials-terracreds", "", -1)
			return path
		}

		if strings.Contains(binary, "terracreds.test") {
			path = strings.Replace(binary, "terracreds.test", "", -1)
			return path
		}

		path = strings.Replace(binary, "terracreds", "", -1)
		return path
	case "linux":
		if strings.Contains(binary, "terraform-credentials-terracreds") {
			path = strings.Replace(binary, "terraform-credentials-terracreds", "", -1)
			return path
		}

		if strings.Contains(binary, "terracreds.test") {
			path = strings.Replace(binary, "terracreds.test", "", -1)
			return path
		}

		path = strings.Replace(binary, "terracreds", "", -1)
		return path
	case "windows":
		if strings.Contains(binary, "terraform-credentials-terracreds.exe") {
			path = strings.Replace(binary, "terraform-credentials-terracreds.exe", "", -1)
			return path
		}

		if strings.Contains(binary, "terracreds.test.exe") {
			path = strings.Replace(binary, "terracreds.test.exe", "", -1)
			return path
		}

		path = strings.Replace(binary, "terracreds.exe", "", -1)
		return path
	default:
		return "Unsupported platform"
	}
}

// NewDirectory checks for the existence of a directory
// if it doesn't exist it creates it and checks for errors
func NewDirectory(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(path, 0644)
			if err != nil {
				return err
			}

			fmt.Fprintf(color.Output, "%s: Created directory '%s'\n", color.CyanString("INFO"), path)
		}
	}
	return nil
}

// WriteToFile will print any string of text to a file safely by
// checking for errors and syncing at the end.
func WriteToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}

	fmt.Fprintf(color.Output, "%s: Created file '%s'\n", color.CyanString("INFO"), filename)
	return file.Sync()
}

// WriteToLog will create a log if it doesn't exist and then append
// messages to the log
func WriteToLog(path string, data string, level string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	logger := log.New(f, level, log.LstdFlags)
	logger.Println(data)
	return nil
}

// CreateConfigFile creates a default terracreds config file if one does not
// exist in the specified file path
func CreateConfigFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			cfgFile := api.Config{
				Logging: api.Logging{
					Enabled: false,
				},
			}

			bytes, err := yaml.Marshal(&cfgFile)
			err = ioutil.WriteFile(path, bytes, 0644)
			if err != nil {
				return err
			}

			fmt.Fprintf(color.Output, "%s: Created file '%s'\n", color.CyanString("INFO"), path)
			return err
		}
	}

	return nil
}

// LoadConfig loads the config file if it exists
func LoadConfig(path string, cfg *api.Config) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, &cfg)
	return err
}

// WriteConfig makes requested changes to config file
func WriteConfig(path string, cfg *api.Config) error {
	bytes, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}

	fmt.Fprintf(color.Output, "%s: Modified config file '%s'\n", color.GreenString("SUCCESS"), path)
	return err
}

// Logging forms the path and writes to log if enabled
func Logging(cfg api.Config, msg string, level string) {
	if cfg.Logging.Enabled == true {
		absolutePath, err := homedir.Expand(cfg.Logging.Path)
		if err != nil {
			log.Fatal(err)
		}

		logPath := filepath.Join(absolutePath, "terracreds.log")
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

// GenerateTerracreds creates the binary to use this package as a credential helper
// and optionally the terraform.rc file
func GenerateTerraCreds(c *cli.Context, version string, confirm string) error {
	var cliConfig string
	var tfPlugins string
	var binary string

	if runtime.GOOS == "windows" {
		userProfile := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		cliConfig = filepath.Join(userProfile, "terraform.rc")
		tfPlugins = filepath.Join(userProfile, "terraform.d", "plugins")
		binary = filepath.Join(tfPlugins, "terraform-credentials-terracreds.exe")
	}

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		userProfile := os.Getenv("HOME")
		cliConfig = filepath.Join(userProfile, ".terraform.d", ".terraformrc")
		tfPlugins = filepath.Join(userProfile, ".terraform.d", "plugins")
		binary = filepath.Join(tfPlugins, "terraform-credentials-terracreds")
	}

	err := NewDirectory(tfPlugins)
	if err != nil {
		return err
	}

	err = CopyTerraCreds(binary)
	if err != nil {
		return err
	}

	if c.Bool("create-cli-config") == true {
		const verbiage = "This command will delete any settings in your .terraformrc file\n\n    Enter 'yes' to coninue or press 'enter' or 'return' to cancel: "
		fmt.Fprintf(color.Output, "%s: %s", color.YellowString("WARNING"), verbiage)
		fmt.Scanln(&confirm)
		fmt.Print("\n")

		if confirm == "yes" {
			doc := heredoc.Doc(`
			credentials_helper "terracreds" {
				args = []
			}`)

			err := WriteToFile(cliConfig, doc)
			return err
		}
	}

	return nil
}

// GetSecretName returns the name of the secret from the config
// or returns the hostname value from the CLI
func GetSecretName(cfg *api.Config, hostname string) string {
	if cfg.Aws.SecretName != "" {
		return cfg.Aws.SecretName
	}
	if cfg.Azure.SecretName != "" {
		return cfg.Azure.SecretName
	}
	if cfg.HashiVault.SecretName != "" {
		return cfg.HashiVault.SecretName
	}
	return hostname
}
