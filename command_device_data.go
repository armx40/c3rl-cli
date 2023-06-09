package main

import (
	"github.com/urfave/cli/v2"
)

var command_devices_subcommands_data_subcommands_device_uid string

func command_devices_subcommands_data_subcommands() (commands cli.Commands) {

	commands = cli.Commands{{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get device settings",
		Action:  command_devices_subcommands_data_subcommands_get,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "uid",
				Aliases:     []string{"u"},
				Value:       "",
				Usage:       "UID of the device",
				Destination: &command_devices_subcommands_data_subcommands_device_uid,
				Required:    true,
			},
		},
	},
	}

	return commands
}

func command_devices_subcommands_data_subcommands_get(cCtx *cli.Context) error {

	return nil
}
