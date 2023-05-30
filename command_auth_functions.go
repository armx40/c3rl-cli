package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/zalando/go-keyring"
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

	/**/

	/* store auth payload */
	err = keyring.Set("c3rl-cli", "auth", string(marshaled_bytes))
	if err != nil {
		log.Fatal(err)
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

	resp, err := network_request("https://c3rl.com/api/esc3rl/user?g=ech", params, headers, nil)
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

func command_auth_functions_get_auth_data() (data userTokenPayload, err error) {
	data_str, err := keyring.Get("c3rl-cli", "auth")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal([]byte(data_str), &data)
	if err != nil {
		return
	}
	return
}

func command_auth_functions_generate_host_device_certificate() (err error) {

	return
}
