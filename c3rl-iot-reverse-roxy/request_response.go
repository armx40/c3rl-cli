package c3rliotroxy

import (
	"encoding/json"
	"fmt"
	pb "main/c3rl-iot-reverse-roxy/protofiles"
	"os"

	"github.com/fatih/color"
)

func request_response_handle_request_packet(packet *pb.WebSocketPacketPayload) (err error) {

	/* handle special packet */

	var request_response_data request_response_request_packet_t

	err = json.Unmarshal(packet.Data, &request_response_data)
	if err != nil {
		return
	}

	switch request_response_data.GType {
	case "cls":
		err = request_response_handle_request_close(packet, &request_response_data)
		return
	default:
		err = fmt.Errorf("invalid gtype")
		return
	}

}

func request_response_handle_request_close(packet *pb.WebSocketPacketPayload, request_data *request_response_request_packet_t) (err error) {

	/* print close message and exist the program */

	closing_message, ok := request_data.Data.(string)

	if !ok {
		err = fmt.Errorf("invalid closing message")
	}

	print_text := color.New(color.FgRed)
	print_text.Printf("%s\n", closing_message)

	os.Exit(1)

	return
}
