package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"time"
)

func command_proxy_functions_linux_endpoint_generate_service_file_data(cli_location string, credential_file string) (service_file_data string, err error) {

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
		credential_file = current_user_home_dir + "/.config/c3rl/credentials.json"
	}

	service_file_data = fmt.Sprintf(`[Unit]
Description=c3rl-cli endpoint service
After=network.target

[Service]
User=%s
Group=%s
ExecStart=%s p e --cr=%s

[Install]
WantedBy=multi-user.target
`, currentUser.Username, currentUser.Username, cli_location, credential_file,
	)
	return

}

func command_proxy_functions_linux_endpoint_generate_install_script(service_file_location string) (script_file_data string, err error) {

	script_file_data = fmt.Sprintf(`sudo cp %s /etc/systemd/system/c3rl-cli-endpoint.service
sudo systemctl enable c3rl-cli-endpoint.service
sudo systemctl start c3rl-cli-endpoint.service
`, service_file_location)

	return
}

func command_proxy_functions_linux_endpoint_generate_uninstall_script() (script_file_data string, err error) {

	script_file_data = `sudo systemctl stop c3rl-cli-endpoint.service
sudo systemctl disable c3rl-cli-endpoint.service
sudo rm /etc/systemd/system/c3rl-cli-endpoint.service
sudo systemctl start c3rl-cli-endpoint.service
sudo systemctl daemon-reload
sudo systemctl reset-failed
`

	return
}

func command_proxy_functions_linux_endpoint_install_routine(cli_location string, credential_file string) (err error) {

	/* generate service file data */
	service_file_data, err := command_proxy_functions_linux_endpoint_generate_service_file_data(cli_location, credential_file)
	if err != nil {
		return
	}
	/**/

	/* save service file in tmp */
	tmp_service_filename := fmt.Sprintf("/tmp/c3rl-cli-endpoint-install-%d.service", time.Now().Unix())

	err = os.WriteFile(tmp_service_filename, []byte(service_file_data), 0644)
	if err != nil {
		return
	}
	/**/

	/* generate install script data */
	install_script_data, err := command_proxy_functions_linux_endpoint_generate_install_script(tmp_service_filename)
	if err != nil {
		return
	}
	/**/

	/* save script file in tmp */
	tmp_script_filename := fmt.Sprintf("/tmp/c3rl-cli-endpoint-install-script-%d.sh", time.Now().Unix())

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

func command_proxy_functions_linux_endpoint_uninstall_routine() (err error) {

	/* generate install script data */
	uninstall_script_data, err := command_proxy_functions_linux_endpoint_generate_uninstall_script()
	if err != nil {
		return
	}
	/**/

	/* save script file in tmp */
	tmp_script_filename := fmt.Sprintf("/tmp/c3rl-cli-endpoint-uninstall-script-%d.sh", time.Now().Unix())

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
