package cmd

import (
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestNewCommandActionForget(t *testing.T) {
	terracreds := config()
	app := app()
	app.Commands = []*cli.Command{
		terracreds.NewCommandCreate(),
		terracreds.NewCommandForget(),
	}

	args := os.Args[0:1]
	args = append(args, "create", "--name=test", "--secret=password")
	app.Run(args)

	args = os.Args[0:1]
	args = append(args, "forget", "test")
	app.Run(args)
}
