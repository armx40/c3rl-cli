package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"github.com/urfave/cli/v2"

	pb "main/protobuf"
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
	},
	// {
	// 	Name:    "gen",
	// 	Aliases: []string{"g"},
	// 	Usage:   "generate device settings",
	// 	Action:  command_devices_subcommands_settings_subcommands_generate,
	// }
	}

	return commands
}

func command_devices_subcommands_settings_subcommands_write(cCtx *cli.Context) error {

	mountpoint := ""

	if !command_devices_subcommands_settings_subcommands_output_to_stdout {
		user_device, err := command_devices_functions_read_user_storage_device()
		if err != nil {
			return err
		}

		/* get selected device and prepare filename */
		mountpoint = user_device.MountPoint
		/**/
	}

	/* start asking for options according to the device selected */
	switch command_devices_subcommands_settings_subcommands_device_type {
	case "atnode":
		return command_devices_subcommands_settings_subcommands_write_atnode(cCtx, mountpoint, command_devices_subcommands_settings_subcommands_output_to_stdout)
	case "at":
		return command_devices_subcommands_settings_subcommands_write_atnode(cCtx, mountpoint, command_devices_subcommands_settings_subcommands_output_to_stdout)
	default:
		return fmt.Errorf("invalid device selected")
	}
	/**/

}

func command_devices_subcommands_settings_subcommands_read(cCtx *cli.Context) (err error) {

	/* ask for storage device */
	user_device, err := command_devices_functions_read_user_storage_device()
	if err != nil {
		return
	}
	/* */

	/* read settings files */
	settings_data_bytes, err := os.ReadFile(filepath.Join(user_device.MountPoint, "settings.data"))
	if err != nil {
		return
	}
	/* */

	/* decoode settings data */
	var settings_data pb.SDCardSettings
	err = proto.Unmarshal(settings_data_bytes, &settings_data)
	if err != nil {
		return
	}
	/**/

	/* decode settings data to survey answer data */
	var data DeviceSettingsSurveyAnswerPayload
	err = data.parse_sdcard_settings(&settings_data)
	if err != nil {
		return
	}
	/* */

	/* pretty print data */
	data.pretty_print()
	/* */
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
