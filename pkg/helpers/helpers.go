package helpers

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/tonedefdev/terracreds/api"
	"gopkg.in/yaml.v2"
)

// CheckError processes the error
func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
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

// Logging forms the path and writes to log if enabled
func Logging(cfg *api.Config, msg string, level string) {
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
