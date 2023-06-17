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
		API_HOST = "https://es.c3rl.com/api/c3rl-cli/"
	} else {
		API_HOST = "http://127.0.0.1:9000/api/c3rl-cli/"
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
				Usage:       "User authentication related functions",
				Subcommands: command_auth_subcommands(),
			},
			{
				Name:        "devices",
				Aliases:     []string{"d"},
				Usage:       "View and manage local and your c3rl devices",
				Subcommands: command_devices_subcommands(),
			},
			{
				Name:        "data",
				Aliases:     []string{"dt"},
				Usage:       "Process and view data",
				Subcommands: command_data_subcommands(),
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print version information",
				Action:  command_version_action,
			},
			{
				Name:        "roxy",
				Aliases:     []string{"r"},
				Usage:       "Roxy functions",
				Subcommands: command_roxy_subcommands(),
			},
			{
				Name:        "update",
				Aliases:     []string{"u"},
				Usage:       "Binary update related functions",
				Subcommands: command_update_subcommands(),
			},
		},
		Name:  "c3rl-cli",
		Usage: "cli application for c3rl services",
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}

}
