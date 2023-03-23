package main

import (
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var command_data_subcommands_process_csv_out_csv_file string
var command_data_subcommands_process_csv_out_csv_file_delimiter string
var command_data_subcommands_process_csv_decrypt bool
var command_data_subcommands_process_csv_tail int

func command_data_subcommands_process_csv_command() (command *cli.Command) {
	command = &cli.Command{

		Name:    "csv",
		Aliases: []string{"c"},
		Usage:   "process and dump data to a csv file",
		Action:  command_data_subcommands_process_csv,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "decrypt",
				Aliases:     []string{"dc"},
				Value:       false,
				Usage:       "use this to decrypt encrypted data",
				Destination: &command_data_subcommands_process_csv_decrypt,
			},
			&cli.IntFlag{
				Name:        "tail",
				Value:       1,
				Usage:       "number of latest data file to process data from",
				Destination: &command_data_subcommands_process_csv_tail,
			},
			&cli.StringFlag{
				Name:        "csv",
				Value:       "",
				Usage:       "csv file to output data to",
				Destination: &command_data_subcommands_process_csv_out_csv_file,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "delimiter",
				Value:       ",",
				Usage:       "delimiter for csv file",
				Destination: &command_data_subcommands_process_csv_out_csv_file_delimiter,
			},
		},
	}
	return command
}

func command_data_subcommands_process_csv(cCtx *cli.Context) (err error) {

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
	/* find how many file have to be read */

	/* */

	return nil
}
