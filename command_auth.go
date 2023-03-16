package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func command_auth_subcommands() (commands cli.Commands) {

	commands = cli.Commands{{

		Name:    "login",
		Aliases: []string{"l"},
		Usage:   "log in to your account",
		Action:  command_auth_login,
	},
	}

	return commands
}

func command_auth_login(cCtx *cli.Context) error {
	fmt.Print("Username: ")

	var username string

	// Taking input from user
	fmt.Scanln(&username)

	/*  */
	return nil
}
