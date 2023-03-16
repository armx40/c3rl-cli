package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

/* csv file */
func command_data_functions_dump_to_csv(csv_filename string, csv_delimiter rune, data *LogLinePayload) error {

	var err error

	if csv_file == nil {
		csv_file, err = os.Create(csv_filename)
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
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
	err := sqlite3_db_tx.Commit()
	if err != nil {
		log.Fatal(err)
		return err
	}

	sqlite3_db_statement.Close()
	sqlite3_db.Close()
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
