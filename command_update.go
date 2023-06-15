package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/kardianos/osext"
	"github.com/urfave/cli/v2"
)

var command_update_location string

func command_update_subcommands() (commands cli.Commands) {

	/* get the current executable path */

	current_executable_path, err := osext.Executable()
	if err != nil {
		return
	}

	commands = cli.Commands{
		{
			Name:    "check",
			Aliases: []string{"c"},
			Usage:   "Check if any newer version is availables",
			Action:  command_update_check,
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Perform update",
			Action:  command_update_update,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "location",
					Aliases:     []string{"l"},
					Value:       current_executable_path,
					Usage:       "Location of currently installed binary",
					Destination: &command_update_location,
					Required:    false,
				},
			},
		},
	}

	return commands
}

func command_update_check(cCtx *cli.Context) (err error) {

	current_version, err := command_update_functions_get_current_version()
	if err != nil {
		return
	}

	update_required, latest_version, err := command_update_functions_check_if_update_is_required()

	if err != nil {
		return
	}

	if update_required {
		print_text := color.New(color.FgYellow)
		print_text.Printf("New update available!\n")
		print_text.Add(color.FgWhite)
		print_text.Printf("Current Version: %d.%d.%d\n", current_version.Major, current_version.Minor, current_version.Patch)
		print_text.Printf("New Version: ")
		print_text.Add(color.FgGreen)
		print_text.Printf(" %d.%d.%d\n", latest_version.Major, latest_version.Minor, latest_version.Patch)
		return
	}

	print_text := color.New(color.FgGreen)
	print_text.Printf("You are running the latest version!\n")
	return
}

func command_update_update(cCtx *cli.Context) (err error) {
	update_required, latest_version, err := command_update_functions_check_if_update_is_required()

	if err != nil {
		return
	}

	if !update_required {
		print_text := color.New(color.FgGreen)
		print_text.Printf("You are running the latest version!\n")
		return
	}

	/* perform update */

	var qs = []*survey.Question{
		{
			Name: "CLILocation",
			Prompt: &survey.Input{
				Message: "c3rl-cli location:",
				Default: command_update_location,
			},
		},
	}

	type answers struct {
		CLILocation string
	}

	var data answers

	err = survey.Ask(qs, &data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = command_update_functions_update_routine(latest_version, data.CLILocation)
	if err != nil {
		return
	}

	print_text := color.New(color.FgGreen)
	print_text.Printf("Successfully updated to version %d.%d.%d!\n", latest_version.Major, latest_version.Minor, latest_version.Patch)

	return
}
