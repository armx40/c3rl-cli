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

const WEBSOCKET_RECONNECT_WAIT_TIME = 10
const WEBSOCKET_RECONNECT_RETRIES = 10

var websocket_conn *websocket.Conn

var websocket_endpoint = ""

// const websocket_endpoint = "ws://localhost:1797/api/proxy/w"

// const websocket_endpoint = "wss://roxy.c3rl.com/api/roxy/w"

const (
	WEBSOCKET_CONNECTION_TYPE_ENDPOINT   = 0
	WEBSOCKET_CONNECTION_TYPE_STARTPOINT = 1 << 1
	WEBSOCKET_CONNECTION_TYPE_EXPOSE     = 1 << 2
)

const (
	WEBSOCKET_PACKET_TYPE_REQUEST_RESPONSE          = 1
	WEBSOCKET_PACKET_TYPE_PROXY                     = 2
	WEBSOCKET_PACKET_TYPE_REQUEST_RESPONSE_REQUEST  = 3
	WEBSOCKET_PACKET_TYPE_REQUEST_RESPONSE_RESPONSE = 4
)

func websocket_init() (err error) {

	websocket_retries := WEBSOCKET_RECONNECT_RETRIES

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
				log.Println(err)
				goto ROUTINE_END
			}

			go func() {
				if main_app_direction == "startpoint" {
					callback_send_message(ROXY_CALLBACK_STARTPOINT_STARTED)
				}
				if main_app_direction == "endpoint" {
					callback_send_message(ROXY_CALLBACK_ENDPOINT_STARTED)
				}
			}()

			websocket_retries = 5

			err = websocket_receive_routine()
			if err != nil {
				goto ROUTINE_END
			}

		ROUTINE_END:

			log.Println(websocket_retries)

			if main_app_direction == "startpoint" {
				// fmt.Printf("cannot establish start point connection\n")
				err = fmt.Errorf("cannot establish start point connection")
				callback_send_message(ROXY_CALLBACK_STARTPOINT_ERROR)
				os.Exit(1)
				return
			}

			if main_app_direction == "endpoint" && main_app_endpoint_exposed_data.ExposedEnable {
				// fmt.Printf("cannot establish expose connection\n")
				err = fmt.Errorf("cannot establish expose connection")
				callback_send_message(ROXY_CALLBACK_EXPOSE_ERROR)
				os.Exit(1)
				return
			}

			websocket_retries = websocket_retries - 1
			if websocket_retries == 0 {
				print_text := color.New(color.FgRed)
				print_text.Println("retries exhausted. Exiting...")

				os.Exit(1)
			}

			print_text := color.New(color.FgYellow)
			print_text.Printf("retrying in %d seconds\n", WEBSOCKET_RECONNECT_WAIT_TIME)

			log.Println(err)
			log.Printf("will restart routine in %d seconds\n", WEBSOCKET_RECONNECT_WAIT_TIME)
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

				var response_g general_payload_v2_t

				err = json.Unmarshal(message, &response_g)
				if err != nil {
					return
				}

				if response_g.Code == APIv2CODEOK {

					interm_bytes, errd := json.Marshal(response_g.Payload)
					if err != nil {
						err = errd
						return
					}

					var response expose_response_payload_t

					err = json.Unmarshal(interm_bytes, &response)
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
					print_text := color.New(color.FgWhite)
					print_text.Printf("%s: ", "endpoint")
					print_text.Add(color.FgRed)
					print_text.Printf("%s\n", response_g.Status)
					print_text.Printf("%s\n", response_g.Payload)
					return fmt.Errorf(response_g.Payload.(string))
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
					print_text.Printf("%s\n", response.Payload)
					return fmt.Errorf(response.Payload.(string))
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
				print_text.Printf("%s\n", response.Payload)
				return fmt.Errorf(response.Payload.(string))
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
			if !main_app_no_fmt_log {
				print_text := color.New(color.FgWhite)
				print_text.Printf("%s: ", "startpoint")
				print_text.Add(color.FgGreen)
				print_text.Printf("%s\n", response.Status)
			}

		} else {
			if !main_app_no_fmt_log {
				print_text := color.New(color.FgWhite)
				print_text.Printf("%s: ", "startpoint")
				print_text.Add(color.FgRed)
				print_text.Printf("%s\n", response.Status)
				print_text.Printf("%s\n", response.Payload)
			}
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

		/* decapsulate packet and process */
		packet, err := packet_decapsulate(&message)
		if err != nil {
			log.Println(err)
			return err
		}
		/**/

		/* check if packet is for some api functions */

		if packet.PacketType == WEBSOCKET_PACKET_TYPE_REQUEST_RESPONSE_REQUEST {
			err = request_response_handle_request_packet(&packet)
			if err != nil {
				log.Println(err)
			}
			continue
		}

		if packet.PacketType == WEBSOCKET_PACKET_TYPE_REQUEST_RESPONSE_RESPONSE {
			// err = request_response_handle_request_packet(&packet)
			// if err != nil {
			// 	log.Println(err)
			// }
			continue
		}

		/* proceed to proxy functions */

		if main_app_direction == "endpoint" {
			err = endpoint.handle_websocket(&packet, &message)
		} else if main_app_direction == "startpoint" {
			err = startpoint.handle_websocket(&packet, &message)
		} else {
			err = fmt.Errorf("invalid app direction")
			return err
		}

	}

}

func websocket_write(data []byte) (err error) {
	return websocket_conn.WriteMessage(websocket.BinaryMessage, data)
}
