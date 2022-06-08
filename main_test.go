package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"runtime"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"

	helpers "github.com/tonedefdev/terracreds/pkg/helpers"
)

func TestWriteToFile(t *testing.T) {
	const fileName = "test.txt"
	const testText = "terracreds test sample text"
	filePath := t.TempDir() + "\\" + fileName

	test := helpers.WriteToFile(filePath, testText)
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

	test := helpers.NewDirectory(filePath)
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
		binaryPath := helpers.GetBinaryPath(path, runtime.GOOS)
		if binaryPath != argPath {
			t.Errorf("Expected '%s' got '%s'", path, binaryPath)
		}

		t.Logf("The binary path sent was '%s' and it correctly returned '%s'", path, binaryPath)
	}
}

func TestCreateConfigFile(t *testing.T) {
	cfgPath := fmt.Sprintf("%s\\config.yaml", t.TempDir())

	err := helpers.CreateConfigFile(cfgPath)
	if err != nil {
		helpers.CheckError(err)
	}

	if cfg.Logging.Enabled != false {
		t.Errorf("Expected logging enabled 'false' got 'true'")
	} else {
		t.Logf("Correctly created the config file")
	}
}

func TestLoadConfig(t *testing.T) {
	cfgPath := fmt.Sprintf("%s\\config.yaml", t.TempDir())
	err := helpers.CreateConfigFile(cfgPath)
	if err != nil {
		helpers.CheckError(err)
	}

	err = helpers.LoadConfig(cfgPath, &cfg)
	if err != nil {
		helpers.CheckError(err)
	}

	if cfg.Logging.Enabled != false {
		t.Errorf("Expected logging enabled 'false' got 'true'")
	} else {
		t.Logf("Correctly loaded the config file")
	}
}

func TestGenerateTerracreds(t *testing.T) {
	var c *cli.Context
	path := t.TempDir()
	tfUser := path + "\\terraform.d"
	helpers.NewDirectory(tfUser)
	helpers.GenerateTerraCreds(c, version, "yes")
}

func TestTerracreds(t *testing.T) {
	const hostname = "terracreds.test.io"
	const apiToken = "9ZWRa0Ge0iQCtA.atlasv1.HpZAd8426rHFskeEFo3AzimnkfR1ldYy69zz0op0NJZ79et8nrgjw3lQfi0FyJ1o8iw"
	const command = "delete"

	terraCreds := NewTerraCreds(runtime.GOOS)
	terraVault := NewTerraVault(&cfg, hostname)

	user, err := user.Current()
	helpers.CheckError(err)

	terraCreds.Create(cfg, hostname, apiToken, user, terraVault)
	token, err := terraCreds.Get(cfg, hostname, user, terraVault)
	if err != nil {
		helpers.CheckError(err)
	}
	fmt.Println(string(token))
	terraCreds.Delete(cfg, command, hostname, user, terraVault)
}
