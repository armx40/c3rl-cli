package main

import (
	"encoding/json"
	"fmt"
	pb "main/c3rl-iot-reverse-roxy"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

var command_roxy_startpoint_config_file string
var command_roxy_startpoint_endpoint_uid string
var command_roxy_credentials_file string
var command_roxy_expose_port uint64
var command_roxy_expose_type string
var command_roxy_expose_domain string

func command_roxy_subcommands() (commands cli.Commands) {

	commands = cli.Commands{{
		Name:    "endpoint",
		Aliases: []string{"e"},
		Usage:   "Use this device as an endpoint",
		Action:  command_roxy_endpoint,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "credentials",
				Aliases:     []string{"cr"},
				Value:       "",
				Usage:       "credentials file for authentication",
				Destination: &command_roxy_credentials_file,
				Required:    false,
			},
		},
	}, {
		Name:    "startpoint",
		Aliases: []string{"s"},
		Usage:   "Use this device as a startpoint and expose ports through which to forward data",
		Action:  command_roxy_startpoint,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Value:       "",
				Usage:       "config file for configuring startpoint",
				Destination: &command_roxy_startpoint_config_file,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "credentials",
				Aliases:     []string{"cr"},
				Value:       "",
				Usage:       "credentials file for authentication",
				Destination: &command_roxy_credentials_file,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "uid",
				Aliases:     []string{"u"},
				Value:       "",
				Usage:       "endpoint uid",
				Destination: &command_roxy_startpoint_endpoint_uid,
				Required:    true,
			},
		},
	}, {
		Name:    "expose",
		Aliases: []string{"ex"},
		Usage:   "Expose",
		Action:  command_roxy_expose,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "credentials",
				Aliases:     []string{"cr"},
				Value:       "",
				Usage:       "Credentials file for authentication",
				Destination: &command_roxy_credentials_file,
				Required:    false,
			},
			&cli.Uint64Flag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       80,
				Usage:       "Port to expose",
				Destination: &command_roxy_expose_port,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "type",
				Aliases:     []string{"t"},
				Value:       "tcp",
				Usage:       "TCP or HTTP",
				Destination: &command_roxy_expose_type,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "domain",
				Aliases:     []string{"d"},
				Value:       "",
				Usage:       "Domain to select",
				Destination: &command_roxy_expose_domain,
				Required:    false,
			},
		},
	}, {
		Name:    "install",
		Aliases: []string{"i"},
		Usage:   "Install endpoint service/job",
		Action:  command_roxy_install_endpoint,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "credentials",
				Aliases:     []string{"cr"},
				Value:       "",
				Usage:       "credentials file for authentication",
				Destination: &command_roxy_credentials_file,
				Required:    false,
			},
		},
	}, {
		Name:    "uninstall",
		Aliases: []string{"u"},
		Usage:   "Uninstall endpoint service/job",
		Action:  command_roxy_uninstall_endpoint,
	},
	}

	return commands
}

func command_roxy_endpoint(cCtx *cli.Context) (err error) {

	/* get auth data */
	auth_data, err := command_auth_functions_get_auth_data()
	if err != nil {
		return
	}
	/* */

	/* get credentials */
	credentials, err := credentials_load_credentials(command_roxy_credentials_file)
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

	auth_data_roxy := pb.Proxy_auth_data_t{
		Token: auth_data.Token,
	}

	exposed_data := pb.Exposed_data_t{
		ExposedEnable: false,
	}

	return pb.StartApp("endpoint", "", "", credentials, &auth_data_roxy, machine_data_bytes, &exposed_data, BuildType == "production")
}

func command_roxy_startpoint(cCtx *cli.Context) (err error) {

	/* get auth data */
	auth_data, err := command_auth_functions_get_auth_data()
	if err != nil {
		return
	}
	/* */

	auth_data_roxy := pb.Proxy_auth_data_t{
		Token: auth_data.Token,
	}

	/* get credentials */
	credentials, err := credentials_load_credentials(command_roxy_credentials_file)
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

	exposed_data := pb.Exposed_data_t{
		ExposedEnable: false,
	}

	return pb.StartApp("startpoint", command_roxy_startpoint_config_file, command_roxy_startpoint_endpoint_uid, credentials, &auth_data_roxy, machine_data_bytes, &exposed_data, BuildType == "production")
}

func command_roxy_expose(cCtx *cli.Context) (err error) {

	/* get auth data */
	auth_data, err := command_auth_functions_get_auth_data()
	if err != nil {
		return
	}
	/* */

	/* get credentials */
	credentials, err := credentials_load_credentials(command_roxy_credentials_file)
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

	auth_data_roxy := pb.Proxy_auth_data_t{
		Token: auth_data.Token,
	}

	exposed_ports := make(map[uint16]pb.Exposed_data_port_t)
	exposed_ports[uint16(command_roxy_expose_port)] = pb.Exposed_data_port_t{
		Type:   command_roxy_expose_type,
		Domain: command_roxy_expose_domain,
	}
	exposed_data := pb.Exposed_data_t{
		ExposedPorts:  exposed_ports,
		ExposedEnable: true,
	}

	return pb.StartApp("endpoint", "", "", credentials, &auth_data_roxy, machine_data_bytes, &exposed_data, BuildType == "production")
}

func command_roxy_install_endpoint(cCtx *cli.Context) (err error) {

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
				Default: current_user_home_dir + "/.config/c3rl/credentials.yaml",
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

	err = command_roxy_functions_endpoint_install_routine(data.CLILocation, data.Credentials)
	if err != nil {
		return
	}
	fmt.Println("endpoint service installed")
	return
}

func command_roxy_uninstall_endpoint(cCtx *cli.Context) (err error) {

	err = command_roxy_functions_endpoint_uninstall_routine()
	if err != nil {
		return
	}
	fmt.Println("endpoint service uninstalled")
	return
}
