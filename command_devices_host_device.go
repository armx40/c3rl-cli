package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

var command_devices_subcommands_host_device_subcommands_generate_credentials_file string
var command_devices_subcommands_host_device_subcommands_generate_credentials_output_to_stdout bool

func command_devices_subcommands_host_device_subcommands() (commands cli.Commands) {

	commands = cli.Commands{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add this host device to your account",
			Action:  command_devices_subcommands_host_device_subcommands_add,
		},
		{
			Name:    "remove",
			Aliases: []string{"r"},
			Usage:   "remove this device from your account",
			Action:  command_devices_subcommands_host_device_subcommands_remove,
		},
		{
			Name:    "credentials",
			Aliases: []string{"c"},
			Usage:   "generate credentials for this device",
			Action:  command_devices_subcommands_host_device_subcommands_generate_credentials,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "out",
					Aliases:     []string{"o"},
					Value:       "",
					Usage:       "output file",
					Destination: &command_devices_subcommands_host_device_subcommands_generate_credentials_file,
				},
				&cli.BoolFlag{
					Name:        "print",
					Aliases:     []string{"p"},
					Value:       false,
					Usage:       "print the credentials on stdout",
					Destination: &command_devices_subcommands_host_device_subcommands_generate_credentials_output_to_stdout,
				},
			},
		},
	}

	return commands
}

func command_devices_subcommands_host_device_subcommands_add(cCtx *cli.Context) error {

	var qs = []*survey.Question{
		{
			Name: "Name",
			Prompt: &survey.Input{
				Message: "Name:",
				Default: "",
			},
		},
		{
			Name: "Description",
			Prompt: &survey.Input{
				Message: "Description:",
			},
		},
	}

	var data HostDeviceAddSurveyPayload

	err := survey.Ask(qs, &data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = command_devices_host_device_functions_add_device(data.Name, data.Description)
	return err
}

func command_devices_subcommands_host_device_subcommands_remove(cCtx *cli.Context) error {

	var qs = []*survey.Question{
		{
			Name: "Ok",
			Prompt: &survey.Confirm{
				Message: "Are you sure you want to remove this host-device from your account?",
			},
		},
	}

	var data GenericYesNoSurveyPayload

	err := survey.Ask(qs, &data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if !data.Ok {
		return nil
	}

	err = command_devices_host_device_functions_remove_device()
	return err
}

func command_devices_subcommands_host_device_subcommands_generate_credentials(cCtx *cli.Context) error {

	err := command_devices_host_device_functions_generate_credentials(nil, command_devices_subcommands_host_device_subcommands_generate_credentials_file, command_devices_subcommands_host_device_subcommands_generate_credentials_output_to_stdout)
	return err
}
