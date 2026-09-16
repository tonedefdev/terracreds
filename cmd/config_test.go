package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tonedefdev/terracreds/api"
	"github.com/urfave/cli/v2"
	"github.com/zalando/go-keyring"
)

func app() *cli.App {
	app := cli.NewApp()
	return app
}

func config() Config {
	dir, _ := os.MkdirTemp("", "terracreds-test-")
	path := filepath.Join(dir, "config.yaml")
	config := Config{
		Cfg: &api.Config{},
		ConfigFile: ConfigFile{
			Path: path,
		},
		TerraCreds: NewTerraCreds(),
	}

	keyring.MockInit()

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
	app := app()
	terracreds := config()
	app.Commands = []*cli.Command{terracreds.NewCommandConfig()}
	args := append(os.Args[0:1], "config", "azure", "--secret-name=test", "--subscription-id=test", "--vault-uri=https://test.com")
	if err := app.Run(args); err != nil {
		t.Fatal(err)
	}
	if err := terracreds.LoadConfig(terracreds.ConfigFile.Path); err != nil {
		t.Fatal(err)
	}

	if got := terracreds.Cfg.Azure; got != (api.Azure{SecretName: "test", SubscriptionId: "test", VaultUri: "https://test.com"}) {
		t.Fatalf("Azure config = %#v", got)
	}
}

func TestActionReset(t *testing.T) {
	app := app()
	terracreds := config()
	app.Commands = []*cli.Command{
		terracreds.NewCommandConfig(),
	}

	args := os.Args[0:1]
	args = append(args, "config", "--use-local-vault-only", "--force")
	app.Run(args)
}
