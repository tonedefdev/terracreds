package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestNewCommandActionGenerate(t *testing.T) {
	terracreds := config()
	app := app()
	app.Commands = []*cli.Command{
		terracreds.NewCommandGenerate(),
	}

	args := os.Args[0:1]
	args = append(args, "generate")
	app.Run(args)
}

func TestCreateConfig(t *testing.T) {
	terracreds := config()
	app := app()
	app.Commands = []*cli.Command{
		terracreds.NewCommandGenerate(),
	}

	args := os.Args[0:1]
	args = append(args, "generate", "--create-cli-config", "--force")
	app.Run(args)

	var configFileName string
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		configFileName = ".terraformrc"
	} else {
		configFileName = "terraform.rc"
	}

	userProfile := os.Getenv("HOME")
	cliConfig := filepath.Join(userProfile, configFileName)

	_, err := os.ReadFile(cliConfig)
	if err != nil {
		t.FailNow()
	}
}
