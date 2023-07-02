package main

import (
	"encoding/json"
	"fmt"
	pb "main/c3rl-iot-reverse-roxy"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var command_roxy_startpoint_config_file string
var command_roxy_startpoint_endpoint_uid string

var command_roxy_credentials_file string

var command_roxy_expose_port uint64
var command_roxy_expose_type string
var command_roxy_expose_domain string

var command_roxy_ssh_endpoint_uid string
var command_roxy_ssh_endpoint_remote_ssh_port uint64
var command_roxy_ssh_endpoint_local_port uint64
var command_roxy_ssh_endpoint_user string
var command_roxy_ssh_endpoint_identity_file string

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
				Required:    false,
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
		{
			Name:    "ssh",
			Aliases: []string{"sh"},
			Usage:   "SSH into remote endpoint",
			Action:  command_roxy_ssh,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "uid",
					Aliases:     []string{"u"},
					Value:       "",
					Usage:       "Endpoint uid",
					Destination: &command_roxy_ssh_endpoint_uid,
					Required:    false,
				},
				&cli.Uint64Flag{
					Name:        "remote-port",
					Aliases:     []string{"p"},
					Value:       0,
					Usage:       "SSH port on the endpoint",
					Destination: &command_roxy_ssh_endpoint_remote_ssh_port,
					Required:    false,
				},
				&cli.Uint64Flag{
					Name:        "local-port",
					Aliases:     []string{"l"},
					Value:       0,
					Usage:       "Local port for the startpoint",
					Destination: &command_roxy_ssh_endpoint_local_port,
					Required:    false,
				},
				&cli.StringFlag{
					Name:        "user",
					Aliases:     []string{"U"},
					Value:       "",
					Usage:       "Endpoint SSH user",
					Destination: &command_roxy_ssh_endpoint_user,
					Required:    false,
				},

				&cli.StringFlag{
					Name:        "identity-file",
					Aliases:     []string{"i"},
					Value:       "",
					Usage:       "Identity for SSH user",
					Destination: &command_roxy_ssh_endpoint_identity_file,
					Required:    false,
				},
			},
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

	return pb.StartApp("endpoint", "", "", credentials, &auth_data_roxy, machine_data_bytes, &exposed_data, BuildType == "production", false, false, nil)
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

	/* exposed data */
	exposed_data := pb.Exposed_data_t{
		ExposedEnable: false,
	}
	/**/

	return pb.StartApp("startpoint", command_roxy_startpoint_config_file, command_roxy_startpoint_endpoint_uid, credentials, &auth_data_roxy, machine_data_bytes, &exposed_data, BuildType == "production", false, false, nil)
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

	return pb.StartApp("endpoint", "", "", credentials, &auth_data_roxy, machine_data_bytes, &exposed_data, BuildType == "production", false, false, nil)
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

func command_roxy_ssh(cCtx *cli.Context) (err error) {

	var endpoint_uid string
	var remote_ssh_port uint16
	var local_startpoint_port uint16
	var user_message string

	/* ask whatever that has not been given by the user */
	/* check if endpoint uid is provided */
	if len(command_roxy_ssh_endpoint_uid) > 0 {
		endpoint_uid = command_roxy_ssh_endpoint_uid
		user_message = command_roxy_ssh_endpoint_uid
	} else {

		/* get connected endpoints */
		endpoints, errd := command_roxy_functions_get_connected_endpoints()
		if err != nil {
			/* check if any uid is provided by the user */
			if len(command_roxy_ssh_endpoint_uid) == 0 {
				err = errd
				return
			}
		}
		/**/

		/* prepare options */
		endpoints_options := []string{}
		for i := range endpoints {

			device_name := endpoints[i].DeviceName
			if len(device_name) == 0 {
				device_name = "<no_name>"
			}

			display_string := fmt.Sprintf("%s, [%s]", device_name, endpoints[i].UID)

			endpoints_options = append(endpoints_options, display_string)

		}
		var qs = []*survey.Question{

			{
				Name: "Endpoint",
				Prompt: &survey.Select{
					Message: "Select Endpoint:",
					Options: endpoints_options,
				},
			},
		}

		type answer_t struct {
			Endpoint int
		}

		var answers answer_t

		err = survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		/* get endpoint to connect to */
		// endpoint_to_connect := endpoints[answers.Endpoint]
		// endpoint_option_to_connect := endpoints_options[answers.Endpoint]

		user_message = endpoints_options[answers.Endpoint]
		endpoint_uid = endpoints[answers.Endpoint].UID
		/**/
	}

	/* ask for username if not given */
	if len(command_roxy_ssh_endpoint_user) == 0 {
		var askUsernamePrompt = &survey.Input{
			Message: "Username:",
		}
		err = survey.AskOne(askUsernamePrompt, &command_roxy_ssh_endpoint_user)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if len(command_roxy_ssh_endpoint_user) == 0 {
			err = fmt.Errorf("Please provide a username.")
			return
		}
	}

	if command_roxy_ssh_endpoint_local_port > 0 {
		local_startpoint_port = uint16(command_roxy_ssh_endpoint_local_port)
	} else {
		local_startpoint_port = 1794
	}

	if command_roxy_ssh_endpoint_remote_ssh_port > 0 {
		remote_ssh_port = uint16(command_roxy_ssh_endpoint_remote_ssh_port)
	} else {
		remote_ssh_port = 22
	}

	print_text := color.New(color.FgWhite)
	print_text.Printf("Connecting to: ")
	print_text.Add(color.FgGreen)
	print_text.Printf("%s on port: %d\n", user_message, remote_ssh_port)
	/* run ssh routine */

	err = command_roxy_functions_ssh_routine(endpoint_uid, command_roxy_ssh_endpoint_user, command_roxy_ssh_endpoint_identity_file, local_startpoint_port, remote_ssh_port)

	return
}
