package main

import (
	"encoding/json"
	"fmt"
	pb "main/c3rl-iot-reverse-proxy"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

var command_proxy_startpoint_config_file string
var command_proxy_startpoint_endpoint_uid string
var command_proxy_credentials_file string

func command_proxy_subcommands() (commands cli.Commands) {

	commands = cli.Commands{{
		Name:    "endpoint",
		Aliases: []string{"e"},
		Usage:   "Use this device as an endpoint",
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
		Usage:   "Use this device as a startpoint and expose ports through which to forward data",
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
			Usage:   "Install endpoint service/job",
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
			Usage:   "Uninstall endpoint service/job",
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

	current_user_home_dir, err := os.UserHomeDir()

	if err != nil {
		return
	}

	var qs = []*survey.Question{
		{
			Name: "CLILocation",
			Prompt: &survey.Input{
				Message: "c3rl-cli location:",
				Default: "/usr/local/bin/c3rl-cli",
			},
		},
		{
			Name: "Credentials",
			Prompt: &survey.Input{
				Message: "credentials:",
				Default: current_user_home_dir + "/.config/c3rl/credentials.json",
			},
		},
	}

	type answers struct {
		CLILocation string
		Credentials string
	}

	var data answers

	err = survey.Ask(qs, &data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = command_proxy_functions_linux_endpoint_install_routine(data.CLILocation, data.Credentials)
	if err != nil {
		return
	}
	fmt.Println("endpoint service installed")
	return
}

func command_proxy_uninstall_endpoint(cCtx *cli.Context) (err error) {

	err = command_proxy_functions_linux_endpoint_uninstall_routine()
	if err != nil {
		return
	}
	fmt.Println("endpoint service uninstalled")
	return
}
