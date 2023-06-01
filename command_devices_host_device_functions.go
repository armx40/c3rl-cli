package main

import (
	"encoding/json"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func command_devices_host_device_functions_get_host_device_info() (data_out *Host_device_payloads_information_data_t, err error) {

	/* get all the info about the machine */
	data_out = &Host_device_payloads_information_data_t{}
	err = data_out.Get()
	/**/

	return
}

func command_devices_host_device_functions_add_device(name string, description string) (err error) {

	/* get auth data */
	auth_data, err := command_auth_functions_get_auth_data()
	if err != nil {
		return
	}
	/* */

	/* get all the info about the machine */
	host_device_data, err := command_devices_host_device_functions_get_host_device_info()
	if err != nil {
		return
	}
	/**/

	/* post this data  */
	var response generalPayloadV2

	request_data := make(map[string]interface{})
	request_data["name"] = name
	request_data["description"] = description
	request_data["device_data"] = *host_device_data

	headers := make(map[string]string)
	headers["Authorization"] = auth_data.Token

	resp, err := network_request(API_HOST+"devices?g=ahd", nil, headers, request_data)

	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &response)
	if err != nil {
		return err
	}
	/**/

	if response.Status != "success" {

		error_response, ok := response.Payload.(string)
		if !ok {
			err = fmt.Errorf("failed to get proper response")
			return
		}
		err = fmt.Errorf(error_response)
		return
	}
	var ok bool
	payload, ok := response.Payload.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("invalid response received")
		return
	}

	uid, ok := payload["uid"].(string)
	if !ok {
		err = fmt.Errorf("invalid response received")
		return
	}
	device_id_, ok := payload["did"].(string)
	if !ok {
		err = fmt.Errorf("invalid response received")
		return
	}
	user_id_, ok := payload["ui"].(string)
	if !ok {
		err = fmt.Errorf("invalid response received")
		return
	}

	user_id, err := primitive.ObjectIDFromHex(user_id_)
	if err != nil {
		return
	}
	device_id, err := primitive.ObjectIDFromHex(device_id_)
	if err != nil {
		return
	}
	reg_data := host_device_credentials_t{
		UID:      uid,
		DeviceID: device_id,
		UserID:   user_id,
	}

	err = command_devices_host_device_functions_generate_credentials(&reg_data, false)
	if err != nil {
		fmt.Printf("device registered but failed to write credentials file. You can still save the credentials using ")
	}
	return
}

func command_devices_host_device_functions_remove_device() (err error) {

	return
}

func command_devices_host_device_functions_get_credentials_from_server() (err error) {

	return
}

func command_devices_host_device_functions_generate_credentials(register_data *host_device_credentials_t, print_to_stdout bool) (err error) {

	if register_data == nil {
		/* not register data present */
	}

	home_dirname, err := os.UserHomeDir()
	if err != nil {
		return
	}

	/* check if .config/c3rl exist or not */
	_, err = os.Stat(home_dirname + "/.config/c3rl")
	if os.IsNotExist(err) {
		err = os.Mkdir(home_dirname+"/.config/c3rl", os.ModePerm)
		if err != nil {
			return
		}
	} else if err == nil {

	} else {
		return err
	}

	configuration_out_file := home_dirname + "/.config/c3rl/credentials.json"

	out_file := configuration_out_file

	if out_file == "" && !print_to_stdout {
		err = fmt.Errorf("no output file specified")
		return
	}

	/* get register data bytes */

	register_data_bytes, err := json.Marshal(register_data)
	if err != nil {
		return
	}

	/**/

	if print_to_stdout {
		fmt.Printf("%s\n", register_data_bytes)
		return
	}

	err = os.WriteFile(out_file, register_data_bytes, 0644)
	return
}
