package main

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	pb "main/protobuf"

	"github.com/golang/protobuf/proto"
	_ "github.com/mattn/go-sqlite3"
)

var command_data_functions_current_data_point = 0

/* csv file */
func command_data_functions_dump_to_csv(csv_filename string, csv_delimiter rune, data *LogLinePayload) error {

	var err error

	if csv_file == nil {
		csv_file, err = os.Create(csv_filename)
		if err != nil {
			log.Fatalf("failed creating f	ile: %s", err)
			return err
		}

		csv_writer = csv.NewWriter(csv_file)
		csv_writer.Comma = csv_delimiter
		_ = csv_writer.Write(data.csv_headers())
	}

	_ = csv_writer.Write(data.csv())
	csv_writer.Flush()

	return nil
}
func command_data_functions_close_csv() error {
	defer csv_file.Close()
	return nil
}

/* */

/* sqlite */
func command_data_functions_dump_to_sqlite(sqlite_filename string, sqlite_table string, data *LogLinePayload) error {

	var err error

	if sqlite3_db == nil {
		sqlite3_db, err = sql.Open("sqlite3", sqlite_filename)
		if err != nil {
			log.Fatal(err)

			return err
		}

		err = command_data_functions_create_table(sqlite_table)
		if err != nil {
			log.Fatal(err)

			return err
		}

		err = command_data_functions_prepare_sqlite_tx(sqlite_table)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	/* add data */
	_, err = sqlite3_db_statement.Exec(data.csv()[0], data.csv()[1], data.csv()[2], data.csv()[3])
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil

}
func command_data_functions_close_sqlite() error {

	/* commit data */
	if sqlite3_db_tx != nil {
		err := sqlite3_db_tx.Commit()
		if err != nil {
			log.Fatal(err)
			return err
		}
		sqlite3_db_statement.Close()
		sqlite3_db.Close()
	}

	return nil
}

func command_data_functions_create_table(table_name string) error {
	sqlStmt := fmt.Sprintf(`
	create table  %s (id integer not null primary key, Time INTEGER,Tag TEXT,Code INTEGER,Log TEXT);
	`, table_name)
	_, err := sqlite3_db.Exec(sqlStmt)
	if err != nil {
		log.Printf("create table failed %q: %s\n", err, sqlStmt)
		return err
	}

	return nil
}

func command_data_functions_prepare_sqlite_tx(table_name string) error {
	var err error
	sqlite3_db_tx, err = sqlite3_db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}

	sqlite3_db_statement, err = sqlite3_db_tx.Prepare(fmt.Sprintf("insert into %s(Time,Tag,Code,Log) values(?, ?, ?, ?)", table_name))
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

/* */

/* main log file */

func command_data_functions_read_log_main_file(filename string) (main_data pb.DataLogMainFile, err error) {

	data, err := os.ReadFile(filename)

	if err != nil {
		return
	}

	err = proto.Unmarshal(data, &main_data)
	if err != nil {
		return
	}

	return
}

/* */

func command_data_process_data_from_file(log_filename string) (data []LogLinePayload, err error) {

	// cCtx.Command.Flags
	/* filename */

	/* open file */
	log_file, err := os.Open(log_filename)
	if err != nil {
		return
	}

	defer log_file.Close()

	/* prepare tmp buffer */

	buffer := make([]byte, len(command_data_log_file_header))

	/* check header */

	n, err := log_file.Read(buffer)
	if err != nil {
		return
	}

	if n != len(command_data_log_file_header) {
		return data, fmt.Errorf("failed read enough bytes")
	}

	/* check if header is valid */

	if string(buffer) != command_data_log_file_header {
		return data, fmt.Errorf("invalid file header")
	}

	/* start processing data */

	for {

		line_len_buffer := make([]byte, 4)
		n, err := log_file.Read(line_len_buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return data, fmt.Errorf("failed to read line len")
		}
		if n != 4 {
			return data, fmt.Errorf("failed read enough bytes")
		}

		/* calculate line line */

		line_len := binary.LittleEndian.Uint32(line_len_buffer)

		if line_len > 1024 {
			return data, fmt.Errorf("invalid line len %d", line_len)
		}

		if line_len == 0 {
			break
		}

		/* read entire line */
		line_buffer := make([]byte, line_len-4)
		n, err = log_file.Read(line_buffer)
		if err != nil {
			return data, fmt.Errorf("failed to read line")

		}

		if n != int(line_len)-4 {
			return data, fmt.Errorf("failed read enough bytes")
		}

		/* process line options */
		line_options := binary.LittleEndian.Uint32(line_buffer[0:4])

		/* process line options and check if decryptiong or sign check or hmac is required */
		is_encrypted := (line_options & DATA_LOG_IS_ENCRYPTED) == 1
		is_signed := (line_options & DATA_LOG_SIGNED) == 1
		is_hmaced := (line_options & DATA_LOG_HMACED) == 1
		/**/

		/* process line time */
		line_time := binary.LittleEndian.Uint32(line_buffer[4:8])

		/* process line tag */
		tag_start_index := 8

		if is_encrypted {

			if len(command_data_subcommands_process_device_uid) == 0 {
				return data, fmt.Errorf("for encrypted data device UID is required")
			}

			decrypt_key, err := command_devices_functions_request_device_symmetric_key(command_data_subcommands_process_device_uid)
			if err != nil {
				return data, err

			}

			iv := line_buffer[8 : 8+16]

			decrypted_line, err := crypto_aes_decrypt(line_buffer[8+16:], decrypt_key[:16], iv)

			if err != nil {
				return data, err

			}

			line_buffer = append(line_buffer[:8], decrypted_line...)

		}

		if is_hmaced {

		}

		if is_signed {

		}

		tag_index := tag_start_index + bytes.IndexByte(line_buffer[tag_start_index:], 0)
		line_tag := string(line_buffer[tag_start_index:tag_index])

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

		data = append(data, curr_data)

	}

	return
}

func command_data_process_dump_data(curr_data *[]LogLinePayload) (dont_look_more bool, err error) {
	/* write to local db?? */

	/* data points to check */
	var data_points_to_check int
	if len(command_data_subcommands_process_sqlite_database_file) > 0 {
		data_points_to_check = command_data_subcommands_process_sqlite_database_count
	} else {
		data_points_to_check = command_data_subcommands_process_csv_count
	}

	for idx := 0; idx < len(*curr_data); idx++ {

		i := idx
		if command_data_subcommands_process_direction_new {
			i = len(*curr_data) - 1 - idx
		}

		// for i := range *curr_data {
		if data_points_to_check > 0 {
			command_data_functions_current_data_point += 1
			if command_data_functions_current_data_point > data_points_to_check {
				return true, nil
			}
		}
		if len(command_data_subcommands_process_sqlite_database_file) > 0 {
			err = command_data_functions_dump_to_sqlite(command_data_subcommands_process_sqlite_database_file, command_data_subcommands_process_sqlite_database_tablename, &(*curr_data)[i])
			if err != nil {
				return
			}
		}
		/* */

		/* save to csv file? */
		if len(command_data_subcommands_process_csv_out_csv_file) > 0 {
			err = command_data_functions_dump_to_csv(command_data_subcommands_process_csv_out_csv_file, []rune(command_data_subcommands_process_csv_out_csv_file_delimiter)[0], &(*curr_data)[i])
			if err != nil {
				return
			}
		}
		/* */
	}

	return
}
