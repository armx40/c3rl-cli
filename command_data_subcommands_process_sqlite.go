package main

import "github.com/urfave/cli/v2"

var command_data_subcommands_process_sqlite_database_file string
var command_data_subcommands_process_sqlite_database_tablename string

func command_data_subcommands_process_sqlite_command() (command *cli.Command) {
	command = &cli.Command{
		Name:    "sqlite",
		Aliases: []string{"s"},
		Usage:   "process and dump data to a sqlite database",
		Action:  command_data_subcommands_process_sqlite,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "sqlite",
				Value:       "",
				Usage:       "create a sqlite database file",
				Destination: &command_data_subcommands_process_sqlite_database_file,
			},
			&cli.StringFlag{
				Name:        "sqlite-table",
				Value:       "data",
				Usage:       "table name for sqlite database file",
				Destination: &command_data_subcommands_process_sqlite_database_tablename,
			},
		},
	}
	return command
}

func command_data_subcommands_process_sqlite(cCtx *cli.Context) error {
	return nil
}
