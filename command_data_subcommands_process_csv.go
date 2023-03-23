package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
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

func command_data_process_data_from_file(log_filename string) error {

	// cCtx.Command.Flags
	/* filename */

	/* open file */

	log_file, err := os.Open(log_filename)
	if err != nil {
		return err
	}

	defer log_file.Close()

	/* prepare tmp buffer */

	buffer := make([]byte, len(command_data_log_file_header))

	/* check header */

	n, err := log_file.Read(buffer)
	if err != nil {
		return err
	}

	if n != len(command_data_log_file_header) {
		return fmt.Errorf("failed read enough bytes")
	}

	/* check if header is valid */

	if string(buffer) != command_data_log_file_header {
		return fmt.Errorf("invalid file header")
	}

	/* start processing data */

	log_lines := []LogLinePayload{}

	for {
		line_len_buffer := make([]byte, 4)
		n, err := log_file.Read(line_len_buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read line len")
		}
		if n != 4 {
			return fmt.Errorf("failed read enough bytes")
		}

		/* calculate line line */

		line_len := binary.LittleEndian.Uint32(line_len_buffer)

		if line_len > 1024 {
			return fmt.Errorf("invalid line len %d", line_len)
		}

		if line_len == 0 {
			break
		}

		/* read entire line */
		line_buffer := make([]byte, line_len-4)
		n, err = log_file.Read(line_buffer)
		if err != nil {
			return fmt.Errorf("failed to read line")

		}

		if n != int(line_len)-4 {
			return fmt.Errorf("failed read enough bytes")
		}

		/* process line options */
		line_options := binary.LittleEndian.Uint32(line_buffer[0:4])

		/* process line options and check if decryptiong or sign check or hmac is required */
		is_encrypted := (line_options & DATA_LOG_IS_ENCRYPTED) == 1
		is_signed := (line_options & DATA_LOG_SIGNED) == 1
		is_hmaced := (line_options & DATA_LOG_HMACED) == 1
		/**/

		if is_encrypted {

		}

		if is_hmaced {

		}

		if is_signed {

		}

		/* process line time */
		line_time := binary.LittleEndian.Uint32(line_buffer[4:8])

		/* process line tag */

		tag_index := 8 + bytes.IndexByte(line_buffer[8:], 0)

		line_tag := string(line_buffer[8:tag_index])

		/* process line code */
		tag_index += 1

		line_code := binary.LittleEndian.Uint32(line_buffer[tag_index : tag_index+4])

		/* process line */

		line_line := line_buffer[tag_index+4:]

		curr_data := LogLinePayload{
			LineOptions: line_options,
			LineTime:    line_time,
			LineTag:     line_tag,
			LineLine:    line_line,
			LineCode:    line_code,
			LineLength:  line_len,
		}

		log_lines = append(log_lines, curr_data)

		/* write to local db?? */

		if len(command_data_subcommands_process_sqlite_database_file) > 0 {
			err = command_data_functions_dump_to_sqlite(command_data_subcommands_process_sqlite_database_file, command_data_subcommands_process_sqlite_database_tablename, &curr_data)
			if err != nil {
				return err
			}
		}

		/* save to csv file? */

		if len(command_data_subcommands_process_csv_out_csv_file) > 0 {
			err = command_data_functions_dump_to_csv(command_data_subcommands_process_csv_out_csv_file, []rune(command_data_subcommands_process_csv_out_csv_file_delimiter)[0], &curr_data)
			if err != nil {
				return err
			}
		}

	}

	if len(command_data_subcommands_process_csv_out_csv_file) > 0 {
		// close the file
		command_data_functions_close_csv()
	}

	if len(command_data_subcommands_process_sqlite_database_file) > 0 {
		// close the file
		command_data_functions_close_sqlite()
	}

	return err
}
