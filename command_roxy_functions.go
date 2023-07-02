package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"

	roxy "main/c3rl-iot-reverse-roxy"
)

func command_roxy_functions_endpoint_linux_generate_service_file_data(cli_location string, credential_file string) (service_file_data string, err error) {

	currentUser, err := user.Current()
	if err != nil {
		return
	}

	current_user_home_dir, err := os.UserHomeDir()

	if err != nil {
		return
	}

	if cli_location == "" {
		cli_location = "/usr/local/bin/c3rl-cli"
	}

	if credential_file == "" {
		credential_file = current_user_home_dir + "/.config/c3rl/credentials.yaml"
	}

	service_file_data = fmt.Sprintf(`[Unit]
Description=c3rl-cli endpoint service
After=network.target

[Service]
User=%s
Group=%s
ExecStart=%s r e --cr=%s

[Install]
WantedBy=multi-user.target
`, currentUser.Username, currentUser.Username, cli_location, credential_file,
	)
	return
}

func command_roxy_functions_endpoint_linux_generate_install_script(service_file_location string) (script_file_data string, err error) {

	script_file_data = fmt.Sprintf(`sudo cp %s /etc/systemd/system/c3rl-cli-endpoint.service
sudo systemctl enable c3rl-cli-endpoint.service
sudo systemctl start c3rl-cli-endpoint.service
`, service_file_location)

	return
}

func command_roxy_functions_endpoint_linux_generate_uninstall_script() (script_file_data string, err error) {

	script_file_data = `sudo systemctl stop c3rl-cli-endpoint.service
sudo systemctl disable c3rl-cli-endpoint.service
sudo rm /etc/systemd/system/c3rl-cli-endpoint.service
sudo systemctl daemon-reload
sudo systemctl reset-failed
`

	return
}

func command_roxy_functions_endpoint_linux_install_routine(cli_location string, credential_file string) (err error) {

	/* generate service file data */
	service_file_data, err := command_roxy_functions_endpoint_linux_generate_service_file_data(cli_location, credential_file)
	if err != nil {
		return
	}
	/**/

	/* save service file in tmp */
	tmp_service_filename := "/tmp/c3rl-cli-endpoint-install.service"
	// tmp_service_filename := fmt.Sprintf("/tmp/c3rl-cli-endpoint-install-%d.service", time.Now().Unix())

	err = os.WriteFile(tmp_service_filename, []byte(service_file_data), 0644)
	if err != nil {
		return
	}
	/**/

	/* generate install script data */
	install_script_data, err := command_roxy_functions_endpoint_linux_generate_install_script(tmp_service_filename)
	if err != nil {
		return
	}
	/**/

	/* save script file in tmp */
	tmp_script_filename := "/tmp/c3rl-cli-endpoint-install-script.sh"
	// tmp_script_filename := fmt.Sprintf("/tmp/c3rl-cli-endpoint-install-script-%d.sh", time.Now().Unix())

	err = os.WriteFile(tmp_script_filename, []byte(install_script_data), 0644)
	if err != nil {
		return
	}
	/**/

	/* exec the script file */
	cmd := exec.Command("bash", tmp_script_filename)
	_, err = cmd.Output()
	if err != nil {
		return
	}
	/**/
	return
}

func command_roxy_functions_endpoint_linux_uninstall_routine() (err error) {

	/* generate install script data */
	uninstall_script_data, err := command_roxy_functions_endpoint_linux_generate_uninstall_script()
	if err != nil {
		return
	}
	/**/

	/* save script file in tmp */
	tmp_script_filename := "/tmp/c3rl-cli-endpoint-uninstall-script.sh"
	// tmp_script_filename := fmt.Sprintf("/tmp/c3rl-cli-endpoint-uninstall-script-%d.sh", time.Now().Unix())

	err = os.WriteFile(tmp_script_filename, []byte(uninstall_script_data), 0644)
	if err != nil {
		return
	}
	/**/

	/* exec the script file */
	cmd := exec.Command("bash", tmp_script_filename)
	_, err = cmd.Output()
	if err != nil {
		return
	}
	/**/
	return
}

/***************************************/

func command_roxy_functions_endpoint_darwin_generate_service_file_data(cli_location string, credential_file string) (service_file_data string, err error) {

	current_user_home_dir, err := os.UserHomeDir()

	if err != nil {
		return
	}

	if cli_location == "" {
		cli_location = "/usr/local/bin/c3rl-cli"
	}

	if credential_file == "" {
		credential_file = current_user_home_dir + "/.config/c3rl/credentials.yaml"
	}

	service_file_data = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN"
	"http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>org.c3rl.cli.endpoint</string>
	<key>ServiceDescription</key>
	<string>c3rl-cli endpoint</string>
	<key>ProgramArguments</key>
	<array>             
		<string>%s</string>
		<string>r</string>
		<string>e</string>
		<string>--cr=%s</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>StandardErrorPath</key>
	<string>/tmp/com.c3rl.cli.endpoint.err</string>
	<key>StandardOutPath</key>
	<string>/tmp/com.c3rl.cli.endpoint.out</string>
</dict>
</plist>
`, cli_location, credential_file,
	)
	return
}

func command_roxy_functions_endpoint_darwin_generate_install_script(service_file_location string) (script_file_data string, err error) {

	script_file_data = fmt.Sprintf(`cp %s ~/Library/LaunchAgents/com.c3rl.cli.endpoint.plist
launchctl load ~/Library/LaunchAgents/com.c3rl.cli.endpoint.plist
launchctl start ~/Library/LaunchAgents/com.c3rl.cli.endpoint.plist
echo installed
`, service_file_location)

	return
}

func command_roxy_functions_endpoint_darwin_generate_uninstall_script() (script_file_data string, err error) {

	script_file_data = `launchctl stop ~/Library/LaunchAgents/com.c3rl.cli.endpoint.plist
launchctl unload ~/Library/LaunchAgents/com.c3rl.cli.endpoint.plist
rm ~/Library/LaunchAgents/com.c3rl.cli.endpoint.plist
`

	return
}

func command_roxy_functions_endpoint_darwin_install_routine(cli_location string, credential_file string) (err error) {

	/* generate service file data */
	service_file_data, err := command_roxy_functions_endpoint_darwin_generate_service_file_data(cli_location, credential_file)
	if err != nil {
		return
	}
	/**/

	/* save service file in tmp */
	tmp_service_filename := "/tmp/com.c3rl.cli.endpoint.plist"
	// tmp_service_filename := fmt.Sprintf("/tmp/c3rl-cli-endpoint-install-%d.service", time.Now().Unix())

	err = os.WriteFile(tmp_service_filename, []byte(service_file_data), 0644)
	if err != nil {
		return
	}
	/**/

	/* generate install script data */
	install_script_data, err := command_roxy_functions_endpoint_darwin_generate_install_script(tmp_service_filename)
	if err != nil {
		return
	}
	/**/

	/* save script file in tmp */
	tmp_script_filename := "/tmp/c3rl-cli-endpoint-install-script.sh"
	// tmp_script_filename := fmt.Sprintf("/tmp/c3rl-cli-endpoint-install-script-%d.sh", time.Now().Unix())

	err = os.WriteFile(tmp_script_filename, []byte(install_script_data), 0644)
	if err != nil {
		return
	}
	/**/

	/* exec the script file */
	cmd := exec.Command("bash", tmp_script_filename)
	_, err = cmd.Output()

	if err != nil {
		return
	}
	/**/
	return
}

func command_roxy_functions_endpoint_darwin_uninstall_routine() (err error) {

	/* generate install script data */
	uninstall_script_data, err := command_roxy_functions_endpoint_darwin_generate_uninstall_script()
	if err != nil {
		return
	}
	/**/

	/* save script file in tmp */
	tmp_script_filename := "/tmp/c3rl-cli-endpoint-uninstall-script.sh"
	// tmp_script_filename := fmt.Sprintf("/tmp/c3rl-cli-endpoint-uninstall-script-%d.sh", time.Now().Unix())

	err = os.WriteFile(tmp_script_filename, []byte(uninstall_script_data), 0644)
	if err != nil {
		return
	}
	/**/

	/* exec the script file */
	cmd := exec.Command("bash", tmp_script_filename)
	_, err = cmd.Output()
	if err != nil {
		return
	}
	/**/
	return
}

/***************************************/

func command_roxy_functions_endpoint_install_routine(cli_location string, credential_file string) (err error) {
	if runtime.GOOS == "linux" {
		err = command_roxy_functions_endpoint_linux_install_routine(cli_location, credential_file)
	} else if runtime.GOOS == "darwin" {
		err = command_roxy_functions_endpoint_darwin_install_routine(cli_location, credential_file)
	} else {
		err = fmt.Errorf("invalid os")
		return
	}
	if err != nil {
		return
	}
	return
}

func command_roxy_functions_endpoint_uninstall_routine() (err error) {
	if runtime.GOOS == "linux" {
		err = command_roxy_functions_endpoint_linux_uninstall_routine()
	} else if runtime.GOOS == "darwin" {
		err = command_roxy_functions_endpoint_darwin_uninstall_routine()
	} else {
		err = fmt.Errorf("invalid os")
		return
	}
	if err != nil {
		return
	}
	return
}

func command_roxy_functions_get_connected_endpoints() (endpoints []endpoint_detailed_response_t, err error) {

	/* get auth data */
	auth_data, err := command_auth_functions_get_auth_data()
	if err != nil {
		return
	}
	/* */

	params := make(map[string]string)
	params["g"] = "gce"

	headers := make(map[string]string)
	headers["Authorization"] = auth_data.Token

	resp, err := network_request(API_HOST+"roxy", params, headers, nil)

	if err != nil {
		log.Println(err)
		return
	}

	var response generalPayloadV2

	err = json.Unmarshal(resp, &response)
	if err != nil {
		return
	}

	if response.Status != "success" {
		err = fmt.Errorf("failed to connected endpoints")
		return
	}

	/* decode data */

	marshaled_bytes, err := json.Marshal(response.Payload)
	if err != nil {
		return
	}

	type response_ struct {
		TotalCount uint64                         `json:"tc" validate:"required"`
		Data       []endpoint_detailed_response_t `json:"data" validate:"required"`
	}

	response_endpoints := response_{}

	err = json.Unmarshal(marshaled_bytes, &response_endpoints)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &response)
	if err != nil {
		return
	}

	endpoints = response_endpoints.Data
	return
}

/***************************************/

func command_roxy_functions_ssh_generate_startpoint_config_data(endpoint_uid string, localport uint16, sshport uint16) (startpoint_config_data string, err error) {

	startpoint_config_data = fmt.Sprintf(`
uid: %s
hostPorts:
- startPointHost: 127.0.0.1
  startPointPort: %d
  endPointHost: 127.0.0.1
  endPointPort: %d
`, endpoint_uid, localport, sshport)

	return
}

func command_roxy_functions_ssh_write_startpoint_config_file(file_data string) (filename string, err error) {

	filename = fmt.Sprintf("/tmp/roxy-ssh-startpoint-%s.yaml", helper_functions_get_random_string(10))

	err = os.WriteFile(filename, []byte(file_data), 0644)

	return
}

func command_roxy_functions_ssh_routine(endpoint_uid string, sshuser string, identity_file string, localport uint16, remote_ssh_port uint16) (err error) {

	startpoint_config_data, err := command_roxy_functions_ssh_generate_startpoint_config_data(endpoint_uid, localport, remote_ssh_port)

	if err != nil {
		return
	}

	startpoint_config_filename, err := command_roxy_functions_ssh_write_startpoint_config_file(startpoint_config_data)
	if err != nil {
		return
	}

	/* start startpoint */

	/* get auth data */
	auth_data, err := command_auth_functions_get_auth_data()
	if err != nil {
		return
	}
	/* */

	auth_data_roxy := roxy.Proxy_auth_data_t{
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

	exposed_data := roxy.Exposed_data_t{
		ExposedEnable: false,
	}

	/* callback channel */

	callback_channel := make(roxy.Roxy_callback_channel)

	err = roxy.StartApp("startpoint", startpoint_config_filename, endpoint_uid, credentials, &auth_data_roxy, machine_data_bytes, &exposed_data, BuildType == "production", true, true, &callback_channel)

	/* wait for callback */
	select {
	case msg := <-callback_channel:
		switch msg {

		case roxy.ROXY_CALLBACK_STARTPOINT_STARTED:
			err = nil
			break

		case roxy.ROXY_CALLBACK_STARTPOINT_ERROR:
			err = fmt.Errorf("startpoint failed to start")
			break

		case roxy.ROXY_CALLBACK_TIMEOUT:
			err = fmt.Errorf("startpoint timed out")
			break
		}

	}
	if err != nil {
		return
	}

	/* deinit callback */
	roxy.Callback_deinit()
	/**/

	/* run ssh command */

	bin, err := helper_functions_find_bin("ssh")
	if err != nil {
		return
	}

	/*  */

	/* what for startpoint callback */

	var cmd *exec.Cmd

	if len(identity_file) > 0 {
		cmd = exec.Command(bin, fmt.Sprintf("%s@127.0.0.1", sshuser), "-p", fmt.Sprint(localport), "-i", identity_file)
	} else {
		cmd = exec.Command(bin, fmt.Sprintf("%s@127.0.0.1", sshuser), "-p", fmt.Sprint(localport))
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
	return
}
