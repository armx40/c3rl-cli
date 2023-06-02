package c3rliotproxy

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

var main_app_direction string
var main_app_startpoint_config_file string
var main_app_endpoint_uid string
var main_app_startpoint_uid string
var main_app_auth_data *Proxy_auth_data_t
var main_app_credentials *Host_device_credentials_t
var main_app_machine_data []byte

func StartApp(direction string, config_file string, endpoint_uid string, credentials *Host_device_credentials_t, auth_data *Proxy_auth_data_t, machine_data []byte) (err error) {

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

	main_app_endpoint_uid = endpoint_uid

	if direction == "startpoint" {
		main_app_startpoint_uid = credentials.UID
	} else {
		main_app_startpoint_uid = ""
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

// func main() {

// 	log.SetFlags(log.LstdFlags | log.Lshortfile)

// 	flag.StringVar(&main_app_direction, "direction", "", "endpoint or startpoint")
// 	flag.StringVar(&main_app_startpoint_config_file, "config", "", "startpoint config file")

// 	flag.Parse()

// 	// start_app("startpoint")
// 	start_app(main_app_direction)

//
// }
