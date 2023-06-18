package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type cli_version_t struct {
	Major       int
	Minor       int
	Patch       int
	BuildNumber int
}

func command_update_functions_get_current_version() (version cli_version_t, err error) {

	major, err := strconv.Atoi(command_version_major)
	if err != nil {
		return
	}

	minor, err := strconv.Atoi(command_version_minor)
	if err != nil {
		return
	}

	patch, err := strconv.Atoi(command_version_patch)
	if err != nil {
		return
	}

	build_number, err := strconv.Atoi(command_version_build_number)
	if err != nil {
		return
	}

	version = cli_version_t{
		Major:       major,
		Minor:       minor,
		Patch:       patch,
		BuildNumber: build_number,
	}

	return
}

func command_update_functions_get_lastest_version() (version cli_version_t, err error) {

	/* call github api to get latest version */

	resp, err := network_request("https://api.github.com/repos/c3rl/c3rl-cli-releases/releases/latest", nil, nil, nil)
	if err != nil {
		return
	}

	var resp_data map[string]interface{}

	err = json.Unmarshal(resp, &resp_data)
	if err != nil {
		return
	}

	/* get latest version string */
	version_string, ok := resp_data["name"].(string)
	if !ok {
		err = fmt.Errorf("failed to get version string")
		return
	}
	/**/

	version_string_split := strings.Split(version_string, ".")

	if len(version_string_split) != 3 {
		err = fmt.Errorf("invalid version string")
		return
	}

	major, err := strconv.Atoi(version_string_split[0])
	if err != nil {
		return
	}

	minor, err := strconv.Atoi(version_string_split[1])
	if err != nil {
		return
	}

	patch, err := strconv.Atoi(version_string_split[2])
	if err != nil {
		return
	}

	version = cli_version_t{
		Major: major,
		Minor: minor,
		Patch: patch,
	}

	return
}

func command_update_functions_check_if_update_is_required() (update_required bool, latest_version cli_version_t, err error) {

	current_version, err := command_update_functions_get_current_version()
	if err != nil {
		return
	}
	latest_version, err = command_update_functions_get_lastest_version()
	if err != nil {
		return
	}

	if latest_version.Major > current_version.Major {
		update_required = true
		return
	} else {
		if latest_version.Minor > current_version.Minor {
			update_required = true
			return
		} else {
			if latest_version.Patch > current_version.Patch {
				update_required = true
				return
			}
		}
	}

	return

}
func command_update_functions_generate_update_script(version cli_version_t, update_location string) (script_file_data string, err error) {

	current_version, err := command_update_functions_get_current_version()
	if err != nil {
		return
	}

	current_version_string := fmt.Sprintf("%d.%d.%d", current_version.Major, current_version.Minor, current_version.Patch)

	version_string := fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)

	current_os := "linux"
	current_arch := "amd64"
	if runtime.GOOS == "darwin" {
		current_os = "darwin"
	}

	if runtime.GOARCH == "arm" {
		current_arch = "arm"
	}

	if runtime.GOARCH == "arm64" {
		current_arch = "arm64"
	}

	cli_filename := fmt.Sprintf(`c3rl-cli_%s_%s_%s`, version_string, current_os, current_arch)

	script_file_data = fmt.Sprintf(`sudo mv %s %s.%s
curl -L -o /tmp/%s.tar https://github.com/c3rl/c3rl-cli-releases/releases/download/%s/%s.tar
tar -xzf /tmp/%s.tar -C /tmp/
chmod +x /tmp/%s
sudo mv /tmp/%s %s
`, update_location, update_location, current_version_string, cli_filename, version_string, cli_filename, cli_filename, cli_filename, cli_filename, update_location)

	return
}

func command_update_functions_update_routine(version cli_version_t, update_location string) (err error) {

	/* generate service file data */
	update_file_data, err := command_update_functions_generate_update_script(version, update_location)
	if err != nil {
		return
	}
	/**/

	/* write update script */
	tmp_script_filename := "/tmp/c3rl-cli-update-script.sh"

	err = os.WriteFile(tmp_script_filename, []byte(update_file_data), 0644)
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

	return err
}
