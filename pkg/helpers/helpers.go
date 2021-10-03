package helpers

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/fatih/color"
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
	bin := GetBinaryPath(os.Args[0], runtime.GOOS)
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
func LoadConfig(cfg *api.Config) error {
	CreateConfigFile()

	bin := GetBinaryPath(os.Args[0], runtime.GOOS)
	path := bin + "config.yaml"
	f, err := os.Open(string(path))
	CheckError(err)
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	CheckError(err)
	return err
}

// Logging forms the path and writes to log if enabled
func Logging(cfg api.Config, msg string, level string) {
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
