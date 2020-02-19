package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/rbwsam/semrel/internal"
)

const (
	app = "semrel"
)

func main() {
	app := &cli.App{
		Name:     app,
		Usage:    "Generate Semantic Versions for your git commits",
		Commands: []*cli.Command{get()},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func get() *cli.Command {
	return &cli.Command{
		Name:        "get",
		Usage:       "Get something",
		Subcommands: []*cli.Command{version()},
	}
}

func version() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Prints a Semantic Version for the head commit of the current git repo",
		Action: func(c *cli.Context) error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			ver, err := internal.CurrentVersion(wd)
			if err != nil {
				return errors.WithStack(err)
			}
			fmt.Println(ver)
			return nil
		},
	}
}
