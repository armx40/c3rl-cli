package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
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
