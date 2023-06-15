package c3rliotroxy

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Proxy_auth_data_t struct {
	Token string
}

type Host_device_credentials_t struct {
	UID      string             `json:"uid"`
	DeviceID primitive.ObjectID `json:"did"`
	UserID   primitive.ObjectID `json:"ui"`
}

type Exposed_data_port_t struct {
	Type   string `json:"t"`
	Domain string `json:"d"`
}

type Exposed_data_t struct {
	ExposedPorts  map[uint16]Exposed_data_port_t `json:"ep"`
	ExposedEnable bool                           `json:"e"`
}

var main_app_direction string
var main_app_startpoint_config_file string
var main_app_endpoint_uid string
var main_app_startpoint_uid string
var main_app_auth_data *Proxy_auth_data_t
var main_app_credentials *Host_device_credentials_t
var main_app_machine_data []byte
var main_app_endpoint_exposed_data *Exposed_data_t

func StartApp(direction string, config_file string, endpoint_uid string, credentials *Host_device_credentials_t, auth_data *Proxy_auth_data_t, machine_data []byte, exposed_data *Exposed_data_t) (err error) {

	/* check if proxy program is already running */

	/* if yes then close this */

	/**/

	/* websocket client init */
	err = remote_connection_init()
	if err != nil {
		return err
	}
	/**/
	main_app_direction = direction
	main_app_startpoint_config_file = config_file
	main_app_auth_data = auth_data
	main_app_credentials = credentials
	main_app_machine_data = machine_data
	main_app_endpoint_exposed_data = exposed_data

	if direction == "startpoint" {
		main_app_startpoint_uid = credentials.UID
		main_app_endpoint_uid = endpoint_uid
	} else {
		main_app_startpoint_uid = ""
		main_app_endpoint_uid = credentials.UID
	}

	if main_app_direction == "startpoint" {
		err = startpoint_init()
	} else if main_app_direction == "endpoint" {
		err = endpoint_init()
	} else {
		return fmt.Errorf("invalid direction")
	}

	if err != nil {
		return err
	}
	/* dont not end */
	select {}
	/**/

}
