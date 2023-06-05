package main

import (
	"github.com/urfave/cli/v2"
)

const command_data_log_file_header = "c3rl - data log - format 1\n"

func command_data_subcommands() (commands cli.Commands) {

	commands = cli.Commands{{
		Name:        "process",
		Aliases:     []string{"p"},
		Usage:       "Process data from file",
		Subcommands: command_data_subcommands_process(),
	},
	}

	return commands
}
