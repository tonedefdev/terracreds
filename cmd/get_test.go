package cmd

import (
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestNewCommandActionGet(t *testing.T) {
	terracreds := config()
	app := app()
	app.Commands = []*cli.Command{
		terracreds.NewCommandCreate(),
		terracreds.NewCommandGet(),
	}

	args := os.Args[0:1]
	args = append(args, "create", "--name=test", "--secret=password")
	app.Run(args)

	args = os.Args[0:1]
	args = append(args, "get", "test")
	app.Run(args)
}
