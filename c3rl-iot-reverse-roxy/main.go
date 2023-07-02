package c3rliotroxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/yaml.v2"
)

type Proxy_auth_data_t struct {
	Token string
}

type Roxy_callback_channel chan int

const (
	ROXY_CALLBACK_STARTPOINT_STARTED = 0
	ROXY_CALLBACK_STARTPOINT_ERROR   = 1
	ROXY_CALLBACK_ENDPOINT_STARTED   = 2
	ROXY_CALLBACK_ENDPOINT_ERROR     = 3
	ROXY_CALLBACK_EXPOSE_STARTED     = 4
	ROXY_CALLBACK_EXPOSE_ERROR       = 5
	ROXY_CALLBACK_TIMEOUT            = 6
)

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
var main_app_endpoint_uid string
var main_app_startpoint_uid string
var main_app_auth_data *Proxy_auth_data_t
var main_app_credentials *Host_device_credentials_t
var main_app_machine_data []byte
var main_app_endpoint_exposed_data *Exposed_data_t
var main_app_is_production bool
var main_app_no_fmt_log bool
var main_app_dont_wait bool
var main_app_callback_channel *Roxy_callback_channel

// var main_app_startpoint_config_file string
var main_app_startpoint_config startpoint_config_t

func StartApp(direction string, config_file string, endpoint_uid string, credentials *Host_device_credentials_t, auth_data *Proxy_auth_data_t, machine_data []byte, exposed_data *Exposed_data_t, is_production bool, no_fmt_log bool, dont_wait bool, callback_channel *Roxy_callback_channel) (err error) {

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
	// main_app_startpoint_config_file = config_file
	main_app_auth_data = auth_data
	main_app_credentials = credentials
	main_app_machine_data = machine_data
	main_app_endpoint_exposed_data = exposed_data
	main_app_is_production = is_production
	main_app_no_fmt_log = no_fmt_log
	main_app_dont_wait = dont_wait
	main_app_callback_channel = callback_channel

	if main_app_callback_channel != nil {
		callback_init()
	}
	if main_app_is_production {
		websocket_endpoint = "wss://roxy.c3rl.com/api/roxy/w"
	} else {
		websocket_endpoint = "ws://localhost:1797/api/proxy/w"
	}

	if direction == "startpoint" {

		yamlFile, errd := os.Open(config_file)
		if errd != nil {
			err = errd
			fmt.Println("FATAL: cannot read config file")
			log.Fatalln("cannot read config file")
			return
		}
		yamlData, errd := ioutil.ReadAll(yamlFile)
		if errd != nil {
			err = errd
			fmt.Println("FATAL: cannot read config file")
			log.Fatalln("cannot read config file")
			return
		}

		err = yaml.Unmarshal([]byte(yamlData), &main_app_startpoint_config)
		if err != nil {
			fmt.Println("FATAL: unable to read contents of config file")
			log.Fatalln("unable to read contents of config file")
			return
		}

		if len(main_app_startpoint_config.UID) > 0 {
			main_app_endpoint_uid = main_app_startpoint_config.UID
		}
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
	if !main_app_dont_wait {
		select {}
	}
	/**/
	return
}
