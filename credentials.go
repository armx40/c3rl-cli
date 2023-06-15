package main

import (
	pb "main/c3rl-iot-reverse-roxy"
	"os"

	"gopkg.in/yaml.v3"
)

func credentials_load_credentials(config_file string) (credentials *pb.Host_device_credentials_t, err error) {

	if config_file == "" {
		home_dirname, err_ := os.UserHomeDir()
		if err != nil {
			err = err_
			return
		}
		config_file = home_dirname + "/.config/c3rl/credentials.yaml"
	}

	credentials_data, err := os.ReadFile(config_file)
	if err != nil {
		return
	}

	credentials = &pb.Host_device_credentials_t{}

	err = yaml.Unmarshal(credentials_data, &credentials)

	return
}
