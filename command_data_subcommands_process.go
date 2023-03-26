package main

import "github.com/urfave/cli/v2"

var command_data_subcommands_process_direction_new bool

func command_data_subcommands_process() (commands cli.Commands) {
	commands = cli.Commands{
		command_data_subcommands_process_csv_command(),
		command_data_subcommands_process_sqlite_command(),
	}
	return commands
}
