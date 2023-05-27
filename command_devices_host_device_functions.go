package main

import (
	"fmt"
	"log"
	"strings"
)

func command_devices_host_device_functions_get_host_device_info() (data_out *host_device_payloads_information_data_t, err error) {

	/* get all the info about the machine */
	data_out = &host_device_payloads_information_data_t{}
	err = data_out.get()
	/**/

	return
}

func command_devices_host_device_functions_add_device(name string, description string) (err error) {

	/* get all the info about the machine */
	data, err := command_devices_host_device_functions_get_host_device_info()
	/**/

	log.Println(data)

	return
}

func command_devices_host_device_functions_remove_device() (err error) {

	return
}

func command_devices_host_device_functions_generate_credentials(out_file string, print_to_stdout bool) (err error) {

	out_file = strings.TrimSpace(out_file)

	if out_file == "" && !print_to_stdout {
		err = fmt.Errorf("no output file specified")
		return
	}

	return
}
