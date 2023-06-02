package main

import (
	"encoding/json"
	pb "main/c3rl-iot-reverse-proxy"

	"github.com/urfave/cli/v2"
)

var command_proxy_startpoint_config_file string
var command_proxy_startpoint_endpoint_uid string
var command_proxy_credentials_file string

func command_proxy_subcommands() (commands cli.Commands) {

	commands = cli.Commands{{
		Name:    "endpoint",
		Aliases: []string{"e"},
		Usage:   "configure this device as endpoint",
		Action:  command_proxy_endpoint,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "credentials",
				Aliases:     []string{"cr"},
				Value:       "",
				Usage:       "credentials file for authentication",
				Destination: &command_proxy_credentials_file,
				Required:    false,
			},
		},
	}, {
		Name:    "startpoint",
		Aliases: []string{"s"},
		Usage:   "use this device as startpoint and expose ports on both ends",
		Action:  command_proxy_startpoint,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Value:       "",
				Usage:       "config file for configuring startpoint",
				Destination: &command_proxy_startpoint_config_file,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "credentials",
				Aliases:     []string{"cr"},
				Value:       "",
				Usage:       "credentials file for authentication",
				Destination: &command_proxy_credentials_file,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "uid",
				Aliases:     []string{"u"},
				Value:       "",
				Usage:       "endpoint uid",
				Destination: &command_proxy_startpoint_endpoint_uid,
				Required:    true,
			},
		},
	},
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "install endpoint service/job",
			Action:  command_proxy_install_endpoint,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "credentials",
					Aliases:     []string{"cr"},
					Value:       "",
					Usage:       "credentials file for authentication",
					Destination: &command_proxy_credentials_file,
					Required:    false,
				},
			},
		},
		{
			Name:    "uninstall",
			Aliases: []string{"u"},
			Usage:   "uninstall endpoint service/job",
			Action:  command_proxy_uninstall_endpoint,
		},
	}

	return commands
}

func command_proxy_endpoint(cCtx *cli.Context) (err error) {

	/* get auth data */
	auth_data, err := command_auth_functions_get_auth_data()
	if err != nil {
		return
	}
	/* */

	/* get credentials */
	credentials, err := credentials_load_credentials(command_proxy_credentials_file)
	if err != nil {
		return
	}
	/**/

	/* get machine data */

	machine_data := Host_device_payloads_information_data_t{}

	err = machine_data.Get()
	if err != nil {
		return
	}

	machine_data_bytes, err := json.Marshal(machine_data)
	if err != nil {
		return
	}
	/**/

	auth_data_proxy := pb.Proxy_auth_data_t{
		Token: auth_data.Token,
	}

	return pb.StartApp("endpoint", "", "", credentials, &auth_data_proxy, machine_data_bytes)
}

func command_proxy_startpoint(cCtx *cli.Context) (err error) {

	/* get auth data */
	auth_data, err := command_auth_functions_get_auth_data()
	if err != nil {
		return
	}
	/* */

	auth_data_proxy := pb.Proxy_auth_data_t{
		Token: auth_data.Token,
	}

	/* get credentials */
	credentials, err := credentials_load_credentials(command_proxy_credentials_file)
	if err != nil {
		return
	}
	/**/
	/* get machine data */

	machine_data := Host_device_payloads_information_data_t{}

	err = machine_data.Get()
	if err != nil {
		return
	}

	machine_data_bytes, err := json.Marshal(machine_data)
	if err != nil {
		return
	}
	/**/

	return pb.StartApp("startpoint", command_proxy_startpoint_config_file, command_proxy_startpoint_endpoint_uid, credentials, &auth_data_proxy, machine_data_bytes)
}

func command_proxy_install_endpoint(cCtx *cli.Context) (err error) {

	return
}

func command_proxy_uninstall_endpoint(cCtx *cli.Context) (err error) {

	return
}
