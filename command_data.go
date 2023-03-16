package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

const command_data_log_file_header = "c3rl - data log - format 1\n"

var command_data_out_csv_file string
var command_data_out_csv_file_delimiter string

var command_data_out_sqlite_database_file string
var command_data_out_sqlite_database_tablename string

func command_data_subcommands() (commands cli.Commands) {

	commands = cli.Commands{{

		Name:    "process",
		Aliases: []string{"p"},
		Usage:   "process data from file",
		Action:  command_data_process_data_from_file,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "csv",
				Value:       "",
				Usage:       "csv file to output data to",
				Destination: &command_data_out_csv_file,
			},
			&cli.StringFlag{
				Name:        "delimiter",
				Value:       ",",
				Usage:       "delimiter for csv file",
				Destination: &command_data_out_csv_file_delimiter,
			},
			&cli.StringFlag{
				Name:        "sqlite",
				Value:       "",
				Usage:       "create a sqlite database file",
				Destination: &command_data_out_sqlite_database_file,
			},
			&cli.StringFlag{
				Name:        "sqlite-table",
				Value:       "data",
				Usage:       "table name for sqlite database file",
				Destination: &command_data_out_sqlite_database_tablename,
			},
		},
	},
	}

	return commands
}

func command_data_process_data_from_file(cCtx *cli.Context) error {

	// cCtx.Command.Flags
	/* filename */

	log_filename := cCtx.Args().Get(0)

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

		if len(command_data_out_sqlite_database_file) > 0 {
			err = command_data_functions_dump_to_sqlite(command_data_out_sqlite_database_file, command_data_out_sqlite_database_tablename, &curr_data)
			if err != nil {
				return err
			}
		}

		/* save to csv file? */

		if len(command_data_out_csv_file) > 0 {
			err = command_data_functions_dump_to_csv(command_data_out_csv_file, []rune(command_data_out_csv_file_delimiter)[0], &curr_data)
			if err != nil {
				return err
			}
		}

	}

	if len(command_data_out_csv_file) > 0 {
		// close the file
		command_data_functions_close_csv()
	}

	if len(command_data_out_sqlite_database_file) > 0 {
		// close the file
		command_data_functions_close_sqlite()
	}

	return err
}
