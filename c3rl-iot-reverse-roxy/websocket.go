package c3rliotroxy

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "main/c3rl-iot-reverse-roxy/protofiles"

	"github.com/fatih/color"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

const WEBSOCKET_RECONNECT_WAIT_TIME = 5

var websocket_conn *websocket.Conn

// const websocket_endpoint = "ws://localhost:1797/api/proxy/w"

const websocket_endpoint = "wss://roxy.c3rl.com/api/roxy/w"

const (
	WEBSOCKET_CONNECTION_TYPE_ENDPOINT   = 0
	WEBSOCKET_CONNECTION_TYPE_STARTPOINT = 1 << 1
	WEBSOCKET_CONNECTION_TYPE_EXPOSE     = 1 << 2
)

const (
	WEBSOCKET_PACKET_TYPE_REQUEST_RESPONSE = 1
	WEBSOCKET_PACKET_TYPE_PROXY            = 2
)

func websocket_init() (err error) {

	go func() {
		for {
			log.Println("websocket initing")
			err := websocket_open()
			if err != nil {
				goto ROUTINE_END
			}
			log.Println("websocket inited")
			err = websocket_init_auth()
			if err != nil {
				goto ROUTINE_END
			}

			err = websocket_verify_auth()
			if err != nil {
				goto ROUTINE_END
			}

			err = websocket_receive_routine()
			if err != nil {
				goto ROUTINE_END
			}

		ROUTINE_END:

			if main_app_direction == "startpoint" {
				fmt.Printf("cannot establish start point connection\n")
				err = fmt.Errorf("cannot establish start point connection")
				os.Exit(1)
				return
			}

			fmt.Printf("disconnected will retry in 5 seconds\n")

			log.Println(err)
			log.Println("will restart routine in 5 seconds")
			time.Sleep(WEBSOCKET_RECONNECT_WAIT_TIME * time.Second)

		}
	}()
	return
}

func websocket_open() (err error) {
	websocket_conn, _, err = websocket.DefaultDialer.Dial(websocket_endpoint, nil)
	if err != nil {
		return
	}
	return
}

func websocket_init_auth() (err error) {

	var expose_data []byte

	conn_type := WEBSOCKET_CONNECTION_TYPE_ENDPOINT

	if main_app_direction == "startpoint" {
		conn_type = WEBSOCKET_CONNECTION_TYPE_STARTPOINT
	}

	if main_app_direction == "endpoint" {
		conn_type = WEBSOCKET_CONNECTION_TYPE_ENDPOINT

		if main_app_endpoint_exposed_data != nil {
			if main_app_endpoint_exposed_data.ExposedEnable {
				conn_type = WEBSOCKET_CONNECTION_TYPE_ENDPOINT | WEBSOCKET_CONNECTION_TYPE_EXPOSE
			}
			expose_data, err = json.Marshal(main_app_endpoint_exposed_data)
			if err != nil {
				return
			}
		}

	}

	/* send endpoint uid in case of startpoint */

	endpoint_uid := main_app_endpoint_uid

	var device_id_to_send string

	device_id_to_send = main_app_credentials.DeviceID.Hex()

	var user_id_to_send string

	user_id_to_send = main_app_credentials.UserID.Hex()

	auth_data := pb.WebSocketAuthPayload{
		ConnectionType: uint32(conn_type),
		Token:          []byte(main_app_auth_data.Token),
		UserId:         []byte(user_id_to_send),
		DeviceId:       []byte(device_id_to_send),
		DeviceData:     main_app_machine_data,
		EndpointUid:    []byte(endpoint_uid),
		StartpointUid:  []byte(main_app_startpoint_uid),
		ExposedData:    expose_data,
	}

	auth_data_bytes, err := proto.Marshal(&auth_data)

	if err != nil {
		return
	}

	err = websocket_conn.WriteMessage(websocket.BinaryMessage, auth_data_bytes)

	return
}

func websocket_verify_auth() (err error) {

	_, message, err := websocket_conn.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return err
	}

	/* process the auth response data */
	/**/

	/* print proper server response */

	if main_app_direction == "endpoint" {
		if main_app_endpoint_exposed_data != nil {
			if main_app_endpoint_exposed_data.ExposedEnable {

				var response expose_response_payload_t

				err = json.Unmarshal(message, &response)
				if err != nil {
					return
				}

				for i := range response.ExposedPorts {
					print_text := color.New(color.FgWhite)
					print_text.Printf("%s:", "tcp.exp.c3rl.com")
					print_text.Add(color.FgGreen)
					print_text.Printf("%d", response.ExposedPorts[i])
					print_text.Add(color.FgWhite)
					print_text.Printf(" -> %s:", "127.0.0.1")
					print_text.Add(color.FgGreen)
					print_text.Printf("%d\n", i)

					log.Printf("exposed port: %d on remote port: %d\n", i, response.ExposedPorts[i])
				}

				for i := range response.ExposedDomainPorts {
					print_text := color.New(color.FgGreen)
					print_text.Printf("https://%s", response.ExposedDomainPorts[i])
					print_text.Add(color.FgWhite)
					print_text.Printf(" -> %s:", "127.0.0.1")
					print_text.Add(color.FgGreen)
					print_text.Printf("%d\n", i)

					log.Printf("exposed port: %d on remote port: %d\n", i, response.ExposedPorts[i])
				}
			} else {
				var response general_payload_v2_t

				err = json.Unmarshal(message, &response)
				if err != nil {
					return
				}

				if response.Code == APIv2CODEOK {
					print_text := color.New(color.FgWhite)
					print_text.Printf("%s: ", "endpoint")
					print_text.Add(color.FgGreen)
					print_text.Printf("%s\n", response.Status)

				} else {
					print_text := color.New(color.FgWhite)
					print_text.Printf("%s: ", "endpoint")
					print_text.Add(color.FgRed)
					print_text.Printf("%s\n", response.Status)
				}
			}
		} else {
			var response general_payload_v2_t

			err = json.Unmarshal(message, &response)
			if err != nil {
				return
			}

			if response.Code == APIv2CODEOK {
				print_text := color.New(color.FgWhite)
				print_text.Printf("%s:", "endpoint: ")
				print_text.Add(color.FgGreen)
				print_text.Printf("%s", response.Status)

			} else {
				print_text := color.New(color.FgWhite)
				print_text.Printf("%s: ", "endpoint")
				print_text.Add(color.FgRed)
				print_text.Printf("%s\n", response.Status)
			}
		}
	}

	if main_app_direction == "startpoint" {
		var response general_payload_v2_t

		err = json.Unmarshal(message, &response)
		if err != nil {
			return
		}

		if response.Code == APIv2CODEOK {
			print_text := color.New(color.FgWhite)
			print_text.Printf("%s: ", "startpoint")
			print_text.Add(color.FgGreen)
			print_text.Printf("%s\n", response.Status)

		} else {
			print_text := color.New(color.FgWhite)
			print_text.Printf("%s: ", "startpoint")
			print_text.Add(color.FgRed)
			print_text.Printf("%s\n", response.Status)
		}
	}

	return
}

func websocket_receive_routine() (err error) {
	for {

		_, message, err := websocket_conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return err
		}

		if main_app_direction == "endpoint" {
			err = endpoint.handle_websocket(&message)
		} else if main_app_direction == "startpoint" {
			err = startpoint.handle_websocket(&message)
		} else {
			err = fmt.Errorf("invalid app direction")
			return err
		}

	}

}

func websocket_write(data []byte) (err error) {
	return websocket_conn.WriteMessage(websocket.BinaryMessage, data)
}
