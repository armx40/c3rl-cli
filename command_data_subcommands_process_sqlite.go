package main

import (
	"path/filepath"

	"github.com/urfave/cli/v2"
)

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
				Name:        "sqlite-database",
				Value:       "",
				Usage:       "create a sqlite database file",
				Destination: &command_data_subcommands_process_sqlite_database_file,
				Required:    true,
			},
			&cli.IntFlag{
				Name:        "tail",
				Value:       1,
				Usage:       "number of latest data file to process data from",
				Destination: &command_data_subcommands_process_csv_tail,
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

func command_data_subcommands_process_sqlite(cCtx *cli.Context) (err error) {

	/* ask for storage device */
	user_device, err := command_devices_functions_read_user_storage_device()
	if err != nil {
		return
	}
	/* */

	/* get the main data log file settings */
	log_settings, err := command_data_functions_read_log_main_file(filepath.Join(user_device.MountPoint, "data_log_main.data"))
	if err != nil {
		return
	}
	/* */

	/* get the latest file index */
	latest_file_index := log_settings.FileCounter % 40000
	/* */

	/* calculate how many files needs to be read */

	number_of_files_to_read := command_data_subcommands_process_csv_tail

	/* */

	/* find of many data logs files are present */
	log_files, err := command_devices_functions_get_all_log_files_sorted(user_device, int(latest_file_index), number_of_files_to_read)
	if err != nil {
		return
	}

	/* */

	/* process each file */

	/* */
	for i := range log_files {
		err = command_data_process_data_from_file(filepath.Join(user_device.MountPoint, log_files[i].Name()))
		if err != nil {
			return err
		}
	}

	/* close sqlite here */
	if len(command_data_subcommands_process_sqlite_database_file) > 0 {
		// close the file
		command_data_functions_close_sqlite()
	}
	/**/
	return nil
}
