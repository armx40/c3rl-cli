package main

import (
	"github.com/urfave/cli/v2"
)

var command_install_location string

func command_install_subcommands() (commands cli.Commands) {

	commands = cli.Commands{
		{
			Name:    "print",
			Aliases: []string{"p"},
			Usage:   "Print installation instructions and commands",
			Action:  command_install_print,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "location",
					Aliases:     []string{"cr"},
					Value:       "",
					Usage:       "install location of the binary",
					Destination: &command_install_location,
					Required:    false,
				},
			},
		},
		{
			Name:    "perform",
			Aliases: []string{"pr"},
			Usage:   "Perform the installation of the binary",
			Action:  command_install_print,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "location",
					Aliases:     []string{"cr"},
					Value:       "",
					Usage:       "install location of the binary",
					Destination: &command_install_location,
					Required:    false,
				},
			},
		},
		{
			Name:    "uninstall",
			Aliases: []string{"u"},
			Usage:   "Perform uninstallation of the binary",
			Action:  command_install_print,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "location",
					Aliases:     []string{"cr"},
					Value:       "",
					Usage:       "install location of the binary",
					Destination: &command_install_location,
					Required:    false,
				},
			},
		},
	}

	return commands
}

func command_install_print(cCtx *cli.Context) (err error) {

	return
}

func command_install_perform(cCtx *cli.Context) (err error) {

	return
}

func command_install_uninstall(cCtx *cli.Context) (err error) {

	return
}
