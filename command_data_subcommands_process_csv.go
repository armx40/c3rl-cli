package main

import (
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var command_data_subcommands_process_csv_out_csv_file string
var command_data_subcommands_process_csv_device_uid string
var command_data_subcommands_process_csv_out_csv_file_delimiter string
var command_data_subcommands_process_csv_decrypt bool
var command_data_subcommands_process_csv_tail int
var command_data_subcommands_process_csv_count int

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
			&cli.IntFlag{
				Name:        "count",
				Value:       -1,
				Usage:       "number of data points to extract",
				Destination: &command_data_subcommands_process_csv_count,
			},
			&cli.BoolFlag{
				Name:        "sample-new",
				Usage:       "sample data starting from newest datapoint",
				Destination: &command_data_subcommands_process_direction_new,
			},
			&cli.StringFlag{
				Name:        "uid",
				Aliases:     []string{"u"},
				Value:       "",
				Usage:       "UID of the device",
				Destination: &command_data_subcommands_process_csv_device_uid,
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

	/* close csv file */
	defer func() {
		if len(command_data_subcommands_process_csv_out_csv_file) > 0 {
			// close the file
			command_data_functions_close_csv()
		}
	}()
	/**/

	/* process each file */

	/* */
	// for i := range log_files {
	read_errors := []string{}
	data_process_errors := []string{}
	data_process_success := []string{}
	number_of_files_read := 0
	for idx := 0; idx < len(log_files); idx++ {
		i := idx
		number_of_files_read += 1

		if !command_data_subcommands_process_direction_new {
			/* because log files are alrady sorted by new */
			i = len(log_files) - 1 - idx
		}

		data, err := command_data_process_data_from_file(filepath.Join(user_device.MountPoint, log_files[i].Name()))

		if err != nil {
			// read_errors = append(read_errors, fmt.Sprintf("File: %s READ FAILED!", log_files[i].Name()))
			read_errors = append(read_errors, fmt.Sprintf("File: %s %s", log_files[i].Name(), err.Error()))
			continue
		}

		/* dump data */

		dont_look_more, err := command_data_process_dump_data(&data)
		if err != nil {
			data_process_errors = append(data_process_errors, fmt.Sprintf("File: %s DATA PROCESSING FAILED!", log_files[i].Name()))
			continue
		}

		data_process_success = append(data_process_success, fmt.Sprintf("File: %s data processing success!", log_files[i].Name()))

		if dont_look_more {
			break
		}

	}

	/* print errors if any */

	if len(read_errors) > 0 {
		for i := range read_errors {
			fmt.Println(read_errors[i])
		}
	}

	if len(data_process_errors) > 0 {
		for i := range data_process_errors {
			fmt.Println(data_process_errors[i])
		}
	}

	if len(read_errors) > 0 || len(data_process_errors) > 0 {
		fmt.Printf("processed data from %d files out of %d\n", len(data_process_success), number_of_files_read)
	}

	/**/
	return nil
}
