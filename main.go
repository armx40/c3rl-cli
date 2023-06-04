package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var BuildType = ""

func main() {

	if BuildType == "production" {
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)
		os.Setenv("GHW_DISABLE_WARNINGS", "1")
	} else {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

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
			{
				Name:    "version",
				Aliases: []string{"v"},
				Action:  command_version_action,
			},
			{
				Name:        "proxy",
				Aliases:     []string{"p"},
				Subcommands: command_proxy_subcommands(),
			},
		},
		Name:  "c3rl-cli",
		Usage: "cli application for c3rl services",
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}

}
