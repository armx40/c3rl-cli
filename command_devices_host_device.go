package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var command_devices_subcommands_host_device_subcommands_generate_credentials_file string
var command_devices_subcommands_host_device_subcommands_generate_credentials_output_to_stdout bool

func command_devices_subcommands_host_device_subcommands() (commands cli.Commands) {

	commands = cli.Commands{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add this host device to your account",
			Action:  command_devices_subcommands_host_device_subcommands_add,
		},
		{
			Name:    "remove",
			Aliases: []string{"r"},
			Usage:   "Remove this device from your account",
			Action:  command_devices_subcommands_host_device_subcommands_remove,
		},
		{
			Name:    "info",
			Aliases: []string{"i"},
			Usage:   "Print the host device credentials information",
			Action:  command_devices_subcommands_host_device_subcommands_info,
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

	reg_data, err := command_devices_host_device_functions_add_device(data.Name, data.Description)

	print_text := color.New(color.FgWhite)
	print_text.Printf("Device registered with UID: ")
	print_text.Add(color.FgGreen)
	print_text.Printf("%s\n", reg_data.UID)

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

func command_devices_subcommands_host_device_subcommands_info(cCtx *cli.Context) error {

	credentials_data, err := command_devices_host_device_functions_read_credentials()

	print_text := color.New(color.FgWhite)
	print_text.Printf("Device UID: ")
	print_text.Add(color.FgGreen)
	print_text.Printf("%s\n", credentials_data.UID)
	return err
}
