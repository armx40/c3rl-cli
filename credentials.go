package main

import (
	"encoding/json"
	pb "main/c3rl-iot-reverse-proxy"
	"os"
)

func credentials_load_credentials(config_file string) (credentials *pb.Host_device_credentials_t, err error) {

	credentials_data, err := os.ReadFile(config_file)
	if err != nil {
		return
	}

	credentials = &pb.Host_device_credentials_t{}

	err = json.Unmarshal(credentials_data, &credentials)

	return
}
