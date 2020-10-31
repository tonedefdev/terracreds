package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/danieljoos/wincred"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func TestWriteToFile(t *testing.T) {
	const fileName = "test.txt"
	const testText = "terracreds test sample text"
	filePath := t.TempDir() + "\\" + fileName

	test := WriteToFile(filePath, testText)
	if test != nil {
		t.Errorf("Unable to write the test file at '%s'", filePath)
	} else {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Errorf("Unable to read the test file at '%s'", filePath)
		} else {
			if string(data) != testText {
				t.Errorf("Expected string '%s' got '%s'", testText, data)
			}
			t.Logf("Read test text '%s' from file '%s'", string(data), filePath)
		}
	}
}

func TestNewDirectory(t *testing.T) {
	const dirName = "terracreds"
	filePath := t.TempDir() + "\\" + dirName

	test := NewDirectory(filePath)
	if test != nil {
		t.Errorf("Unable to create the test directory at '%s'", filePath)
	}
	t.Logf("Created test directory at '%s'", filePath)
}

func TestGetBinaryPath(t *testing.T) {
	var paths [3]string
	argPath := strings.Replace(os.Args[0], "terracreds.test.exe", "", -1)
	paths[0] = argPath + "terraform-credentials-terracreds.exe"
	paths[1] = argPath + "terracreds.test.exe"
	paths[2] = argPath + "terracreds.exe"

	for _, path := range paths {
		binaryPath := GetBinaryPath(path)
		if binaryPath != argPath {
			t.Errorf("Expected '%s' got '%s'", path, binaryPath)
		}

		t.Logf("The binary path sent was '%s' and it correctly returned '%s'", path, binaryPath)
	}
}

func TestCreateConfigFile(t *testing.T) {
	bin := GetBinaryPath(os.Args[0])
	path := bin + "config.yaml"
	var doc string

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			doc = heredoc.Doc(`
logging:
  enabled: false
  path:`)
			WriteToFile(path, doc)
		}
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("Unable to read the config file at '%s'", path)
	} else {
		if string(data) != doc {
			t.Errorf("Expected config '%s' got '%s'", doc, string(data))
		}
		t.Logf("Read config '%s' from file '%s'", string(data), path)
	}
}

func TestLoadConfig(t *testing.T) {
	var cfg Config
	CreateConfigFile()

	bin := GetBinaryPath(os.Args[0])
	path := bin + "config.yaml"
	file, err := os.Open(string(path))
	CheckError(err)
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if cfg.Logging.Enabled != false {
		t.Errorf("Expected logging enabled 'false' got 'true'")
	}
	t.Logf("Correctly loaded the config file")
}

func TestGenerateTerracreds(t *testing.T) {
	var c *cli.Context
	path := t.TempDir()
	tfUser := path + "\\terraform.d"
	NewDirectory(tfUser)
	GenerateTerracreds(c, path, tfUser)
}

func TestCreateCredential(t *testing.T) {
	var c *cli.Context
	const hostname = "terracreds.test.io"
	const apiToken = "9ZWRa0Ge0iQCtA.atlasv1.HpZAd8426rHFskeEFo3AzimnkfR1ldYy69zz0op0NJZ79et8nrgjw3lQfi0FyJ1o8iw"

	CreateCredential(c, hostname, apiToken)
	cred, err := wincred.GetGenericCredential(hostname)
	if err != nil {
		t.Errorf("Expected credential object '%s' got '%s'", hostname, string(cred.CredentialBlob))
	}
	cred.Delete()
}
