package cmd

import (
	"os"
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
