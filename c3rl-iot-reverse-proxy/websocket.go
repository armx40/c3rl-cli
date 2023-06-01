package c3rliotproxy

import (
	"fmt"
	"log"
	"os"
	"time"

	pb "main/c3rl-iot-reverse-proxy/protofiles"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

const WEBSOCKET_RECONNECT_WAIT_TIME = 5

var websocket_conn *websocket.Conn

const websocket_endpoint = "wss://es.c3rl.com/api/proxy/w"

const (
	WEBSOCKET_CONNECTION_TYPE_ENDPOINT   = 0
	WEBSOCKET_CONNECTION_TYPE_STARTPOINT = 1
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
			log.Println(err)
			log.Println("Will restart routine in 5 seconds")
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

	conn_type := WEBSOCKET_CONNECTION_TYPE_ENDPOINT

	if main_app_direction == "startpoint" {
		conn_type = WEBSOCKET_CONNECTION_TYPE_STARTPOINT
	}

	/* send endpoint uid in case of startpoint */
	uid_to_send := main_app_credentials.UID
	if main_app_direction == "startpoint" {
		uid_to_send = main_app_endpoint_uid
	}

	auth_data := pb.WebSocketAuthPayload{
		ConnectionType: uint32(conn_type),
		Uid:            []byte(uid_to_send),
		Token:          []byte(main_app_auth_data.Token),
		UserId:         []byte(main_app_credentials.UserID.Hex()),
		DeviceId:       []byte(main_app_credentials.DeviceID.Hex()),
		DeviceData:     main_app_machine_data,
	}

	auth_data_bytes, err := proto.Marshal(&auth_data)

	if err != nil {
		return
	}

	err = websocket_conn.WriteMessage(websocket.BinaryMessage, auth_data_bytes)

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