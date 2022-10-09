package cmd

import (
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestNewCommandActionCreate(t *testing.T) {
	terracreds := config()
	app := app()
	app.Commands = []*cli.Command{
		terracreds.NewCommandCreate(),
	}

	args := os.Args[0:1]
	args = append(args, "create", "--name=test", "--secret=password")
	app.Run(args)
}
