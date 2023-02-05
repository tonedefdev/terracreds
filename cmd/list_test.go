package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

func createCases(app *cli.App) {
	var args []string
	cases := []string{"test", "test2"}

	for _, v := range cases {
		args = os.Args[0:1]
		name := fmt.Sprintf("--name=%s", v)
		args = append(args, "create", name, "--secret=password")
		app.Run(args)
	}
}

func deleteCases(app *cli.App) {
	var args []string
	cases := []string{"test", "test2"}

	for _, v := range cases {
		args = os.Args[0:1]
		name := fmt.Sprintf("--name=%s", v)
		args = append(args, "delete", name)
		app.Run(args)
	}
}

func TestNewCommandActionList(t *testing.T) {
	terracreds := config()
	app := app()
	app.Commands = []*cli.Command{
		terracreds.NewCommandCreate(),
		terracreds.NewCommandList(),
		terracreds.NewCommandDelete(),
	}

	createCases(app)

	args := os.Args[0:1]
	args = append(args, "list", "-l=test,test2")
	app.Run(args)

	deleteCases(app)
}

func TestNewCommandActionListAsJson(t *testing.T) {
	terracreds := config()
	app := app()
	app.Commands = []*cli.Command{
		terracreds.NewCommandCreate(),
		terracreds.NewCommandList(),
		terracreds.NewCommandDelete(),
	}

	createCases(app)

	args := os.Args[0:1]
	args = append(args, "list", "-l=test,test2", "--as-json")
	app.Run(args)

	deleteCases(app)
}

func TestNewCommandActionListAsTFVars(t *testing.T) {
	terracreds := config()
	app := app()
	app.Commands = []*cli.Command{
		terracreds.NewCommandCreate(),
		terracreds.NewCommandList(),
		terracreds.NewCommandDelete(),
	}

	createCases(app)

	args := os.Args[0:1]
	args = append(args, "list", "-l=test,test2", "--as-tfvars")
	app.Run(args)

	deleteCases(app)
}
