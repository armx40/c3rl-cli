package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func runtime_generate_run_file(location string) (err error) {

	if location == "" {
		location = "/tmp/c3rl-cli.pid"
	}

	pid := os.Getpid()

	pid_s := fmt.Sprint(pid)

	err = os.WriteFile(location, []byte(pid_s), 0644)

	return

}

func runtime_check_run_file(location string) (err error) {

	if location == "" {
		location = "/tmp/c3rl-cli.pid"
	}

	_, err = os.Stat(location)
	if os.IsNotExist(err) {
		return nil
	} else if err == nil {
		print_text := color.New(color.FgRed)
		print_text.Println("Another instance of c3rl-cli is already running.")
		return fmt.Errorf("Another instance of c3rl-cli is already running.")
	} else {
		return err
	}
}

func runtime_purge_run_file(location string) (err error) {

	if location == "" {
		location = "/tmp/c3rl-cli.pid"
	}

	err = os.Remove(location)

	return

}
