package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	err := crypto_init()
	if err != nil {
		fmt.Println("failed initialize crypto")
		return
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:        "auth",
				Aliases:     []string{"a"},
				Subcommands: command_auth_subcommands(),
			},
			{
				Name:        "devices",
				Aliases:     []string{"d"},
				Subcommands: command_devices_subcommands(),
			},
			{
				Name:        "data",
				Aliases:     []string{"dt"},
				Subcommands: command_data_subcommands(),
			},
		},
		Name:  "c3rl-cli",
		Usage: "cli application for c3rl services",
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}

}
