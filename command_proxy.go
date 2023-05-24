package main

import "github.com/urfave/cli/v2"

func command_proxy_subcommands() (commands cli.Commands) {

	commands = cli.Commands{{

		Name:    "login",
		Aliases: []string{"l"},
		Usage:   "log in to your account",
		Action:  command_auth_login,
	}, {

		Name:    "verify",
		Aliases: []string{"v"},
		Usage:   "verify the auth status",
		Action:  command_auth_verify,
	},
	}

	return commands
}
