package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/tonedefdev/terracreds/api"
	"github.com/urfave/cli/v2"
)

func app() *cli.App {
	app := cli.NewApp()
	return app
}

func config() Config {
	dir, _ := os.UserHomeDir()
	path := filepath.Join(dir, "config.yaml")
	config := Config{
		Cfg: &api.Config{},
		ConfigFile: ConfigFile{
			Path: path,
		},
		TerraCreds: NewTerraCreds(runtime.GOOS),
	}

	return config
}

func contains[T comparable](elements []T, val T) bool {
	for _, v := range elements {
		if v == val {
			return true
		}
	}

	return false
}

func TestNewCommandConfig(t *testing.T) {
	app := app()
	terracreds := config()
	app.Commands = []*cli.Command{
		terracreds.NewCommandConfig(),
	}

	args := os.Args[0:1]
	args = append(args, "config")
	app.Run(args)
}

func TestNewCommandAws(t *testing.T) {
	app := app()
	terracreds := config()
	app.Commands = []*cli.Command{
		terracreds.NewCommandConfig(),
	}

	args := os.Args[0:1]
	args = append(args, "config", "aws")
	app.Run(args)
}

func TestNewCommandActionAws(t *testing.T) {
	app := app()
	terracreds := config()
	app.Commands = []*cli.Command{
		terracreds.NewCommandConfig(),
	}

	args := os.Args[0:1]
	args = append(args, "config", "aws", "--description=test", "--region=test", "--secret-name=test")
	app.Run(args)
}

func TestNewCommandAzure(t *testing.T) {
	app := app()
	terracreds := config()
	app.Commands = []*cli.Command{
		terracreds.NewCommandConfig(),
	}

	args := os.Args[0:1]
	args = append(args, "config", "azure")
	app.Run(args)
}

func TestNewCommandActionAzure(t *testing.T) {
	app := app()
	terracreds := config()
	app.Commands = []*cli.Command{
		terracreds.NewCommandConfig(),
	}

	args := os.Args[0:1]
	args = append(args, "config", "azure", "--secret-name=test", "--subscription-id=test", "--vault-uri=https://test.com")
	app.Run(args)
}

func TestActionAzureResult(t *testing.T) {
	terracreds := config()
	terracreds.LoadConfig(terracreds.ConfigFile.Path)

	failures := make([]bool, 3)

	if terracreds.Cfg.Azure.SecretName != "test" {
		t.Logf("Azure.SecretName is '%s' expected 'test'", terracreds.Cfg.Azure.SecretName)
		failures = append(failures, true)
	}

	if terracreds.Cfg.Azure.SubscriptionId != "test" {
		t.Logf("Azure.SecretName is '%s' expected 'test'", terracreds.Cfg.Azure.SecretName)
		failures = append(failures, true)
	}

	if terracreds.Cfg.Azure.VaultUri != "https://test.com" {
		t.Logf("Azure.SecretName is '%s' expected 'test'", terracreds.Cfg.Azure.SecretName)
		failures = append(failures, true)
	}

	if contains(failures, true) {
		t.Fail()
	}
}

func TestActionReset(t *testing.T) {
	app := app()
	terracreds := config()
	app.Commands = []*cli.Command{
		terracreds.NewCommandConfig(),
	}

	args := os.Args[0:1]
	args = append(args, "config", "--use-local-vault-only")
	app.Run(args)
}
