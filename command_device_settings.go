package main

import (
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var command_devices_subcommands_settings_subcommands_device_type string
var command_devices_subcommands_settings_subcommands_output_to_stdout bool

func command_devices_subcommands_settings_subcommands() (commands cli.Commands) {

	commands = cli.Commands{{
		Name:    "write",
		Aliases: []string{"w"},
		Usage:   "write device settings",
		Action:  command_devices_subcommands_settings_subcommands_write,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "device-type",
				Aliases:     []string{"dt"},
				Value:       "",
				Usage:       "device type",
				Destination: &command_devices_subcommands_settings_subcommands_device_type,
				Required:    true,
			},
			&cli.BoolFlag{
				Name:        "out",
				Aliases:     []string{"o"},
				Value:       false,
				Usage:       "output the setting to stdout",
				Destination: &command_devices_subcommands_settings_subcommands_output_to_stdout,
			},
		},
	}, {
		Name:    "read",
		Aliases: []string{"r"},
		Usage:   "read device settings",
		Action:  command_devices_subcommands_settings_subcommands_read,
	}, {
		Name:    "gen",
		Aliases: []string{"g"},
		Usage:   "generate device settings",
		Action:  command_devices_subcommands_settings_subcommands_generate,
	}}

	return commands
}

func command_devices_subcommands_settings_subcommands_write(cCtx *cli.Context) error {

	filename := ""

	if !command_devices_subcommands_settings_subcommands_output_to_stdout {
		devices, err := command_devices_functions_find_sdcard_device()
		if err != nil {
			return err
		}

		if len(devices) == 0 {
			fmt.Println("no storage present")
			return nil
		}
		for i, device := range devices {
			devices = append(devices, device)
			fmt.Printf("[%d]:  %v\n", i, device)
		}

		fmt.Print("Select storage device: ")

		var device_idx int
		fmt.Scan(&device_idx)

		if device_idx > len(devices)-1 {
			return fmt.Errorf("invalid storage device")
		}

		/* get selected device and prepare filename */
		filename = filepath.Join(devices[device_idx].MountPoint, "settings.data")
		/**/
	}

	/* start asking for options according to the device selected */
	switch command_devices_subcommands_settings_subcommands_device_type {
	case "atnode":
		return command_devices_subcommands_settings_subcommands_write_atnode(cCtx, filename, command_devices_subcommands_settings_subcommands_output_to_stdout)
	case "at":
		return command_devices_subcommands_settings_subcommands_write_atnode(cCtx, filename, command_devices_subcommands_settings_subcommands_output_to_stdout)
	default:
		return fmt.Errorf("invalid device selected")
	}
	/**/

}

func command_devices_subcommands_settings_subcommands_read(cCtx *cli.Context) error {
	devices, err := command_devices_functions_find_sdcard_device()
	if err != nil {
		return err
	}

	for _, device := range devices {
		devices = append(devices, device)
		fmt.Printf("  %v\n", device)
	}

	return nil
}

func command_devices_subcommands_settings_subcommands_generate(cCtx *cli.Context) error {
	devices, err := command_devices_functions_find_sdcard_device()
	if err != nil {
		return err
	}

	for _, device := range devices {
		devices = append(devices, device)
		fmt.Printf("  %v\n", device)
	}

	return nil
}
