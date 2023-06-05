package main

import (
	"flag"
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	debug = flag.Int("debug", 0, "libusb debug level (0..3)")
)

func command_devices_subcommands() (commands cli.Commands) {

	commands = cli.Commands{{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "list connected devices",
		Action:  command_devices_list,
	}, {
		Name:        "settings",
		Aliases:     []string{"s"},
		Usage:       "Manage connected device settings",
		Subcommands: command_devices_subcommands_settings_subcommands(),
	}, {
		Name:        "data",
		Aliases:     []string{"d"},
		Usage:       "Get connected device data",
		Subcommands: command_devices_subcommands_data_subcommands(),
	}, {
		Name:        "host-device",
		Aliases:     []string{"hd"},
		Usage:       "Manage/configure host device",
		Subcommands: command_devices_subcommands_host_device_subcommands(),
	}}

	return commands
}

func command_devices_list(cCtx *cli.Context) error {
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
