package main

import (
	pb "main/c3rl-iot-reverse-proxy"

	"github.com/urfave/cli/v2"
)

var command_proxy_startpoint_config_file string

func command_proxy_subcommands() (commands cli.Commands) {

	commands = cli.Commands{{
		Name:    "endpoint",
		Aliases: []string{"e"},
		Usage:   "configure this device as endpoint",
		Action:  command_proxy_endpoint,
	}, {
		Name:    "startpoint",
		Aliases: []string{"s"},
		Usage:   "use this device as startpoint and expose ports on both ends",
		Action:  command_proxy_startpoint,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Value:       "",
				Usage:       "config file for configuring startpoint",
				Destination: &command_proxy_startpoint_config_file,
				Required:    true,
			},
		},
	},
	}

	return commands
}

func command_proxy_endpoint(cCtx *cli.Context) error {
	return pb.StartApp("endpoint", "")
}

func command_proxy_startpoint(cCtx *cli.Context) error {

	return pb.StartApp("startpoint", command_proxy_startpoint_config_file)
}
