package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func command_auth_functions_login(username string, password string) (err error) {

	var response generalPayloadV2

	request_data := make(map[string]string)
	request_data["user"] = username
	request_data["pass"] = password

	resp, err := network_request(API_HOST+"auth?g=lgn", nil, nil, request_data)

	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &response)
	if err != nil {
		return err
	}

	if response.Status != "success" {
		return fmt.Errorf("login failed")
	}

	/* decode data */

	marshaled_bytes, err := json.Marshal(response.Payload)
	if err != nil {
		return err
	}

	var token_data userTokenPayload

	err = json.Unmarshal(marshaled_bytes, &token_data)
	if err != nil {
		return err
	}

	err = command_auth_functions_set_auth_data(&token_data)
	if err != nil {
		return err
	}
	/**/

	return
}

func command_auth_functions_echo() (err error) {

	var response generalPayloadV2

	params := make(map[string]string)
	params["g"] = "ech"

	/* get auth data */
	auth_data, err := command_auth_functions_get_auth_data()
	if err != nil {
		return
	}
	/* */

	headers := make(map[string]string)
	headers["Authorization"] = auth_data.Token

	resp, err := network_request(API_HOST+"auth?g=ech", params, headers, nil)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &response)
	if err != nil {
		return err
	}

	if response.Status != "success" {
		return fmt.Errorf("echo failed")
	}

	if response.Payload != "ohce" {
		return fmt.Errorf("unknown error")
	}
	return nil
}

func command_auth_functions_set_auth_data(token_data *userTokenPayload) (err error) {

	/* store auth payload */
	// json_bytes, err := json.Marshal(token_data)
	// if err != nil {
	// 	err = fmt.Errorf("failed to get json data")
	// 	return
	// }
	// err = keyring.Set("c3rl-cli", "auth", string(json_bytes))
	// if err != nil {
	// 	err = fmt.Errorf("Failed to store auth token. Make sure you have access to system's keyring.")
	// 	return
	// }
	/**/

	/* store in auth file in config */
	yaml_data, err := yaml.Marshal(&token_data)
	if err != nil {
		err = fmt.Errorf("failed to decode auth data")
		return
	}

	user_home_dir, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("failed to get user home directory")
		return
	}

	/* check if .config/c3rl exist or not */
	_, err = os.Stat(user_home_dir + "/.config/c3rl")
	if os.IsNotExist(err) {
		err = os.Mkdir(user_home_dir+"/.config/c3rl", os.ModePerm)
		if err != nil {
			return
		}
	} else if err == nil {

	} else {
		return err
	}

	err = os.WriteFile(user_home_dir+"/.config/c3rl/auth.yaml", yaml_data, 0644)
	if err != nil {
		err = fmt.Errorf("failed to store auth token")
		return
	}
	/**/

	return

}
func command_auth_functions_get_auth_data() (data userTokenPayload, err error) {

	/* get from keyring */
	// data_str, err := keyring.Get("c3rl-cli", "auth")
	// if err != nil {
	// 	return
	// }

	// err = json.Unmarshal([]byte(data_str), &data)
	// if err != nil {
	// 	return
	// }
	/* */

	/* get from auth yaml */
	user_dir, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("failed to get user home directory")
		return
	}
	auth_yaml_data, err := os.ReadFile(user_dir + "/.config/c3rl/auth.yaml")
	if err != nil {
		err = fmt.Errorf("failed to get auth token")
		return
	}
	err = yaml.Unmarshal(auth_yaml_data, &data)
	if err != nil {
		err = fmt.Errorf("failed to decode auth yaml")
		return
	}
	/**/
	return
}

func command_auth_functions_generate_host_device_certificate() (err error) {

	return
}
