package c3rliotproxy

import (
	"encoding/binary"
	"net"

	pb "main/c3rl-iot-reverse-proxy/protofiles"

	"github.com/golang/protobuf/proto"
)

func packet_encapsulate(connection_id []byte, endpoint_host string, endpoint_port uint, startpoint_host string, startpoint_port uint, data []byte) (packet_bytes []byte, err error) {

	endpoint_host_uint32 := net.ParseIP(endpoint_host)
	startpoint_host_uint32 := net.ParseIP(startpoint_host)

	packet := pb.WebSocketPacketPayload{
		EHost:        binary.BigEndian.Uint32(endpoint_host_uint32.To4()),
		EPort:        uint32(endpoint_port),
		SHost:        binary.BigEndian.Uint32(startpoint_host_uint32.To4()),
		SPort:        uint32(startpoint_port),
		Data:         data,
		ConnectionId: connection_id,
		PacketType:   WEBSOCKET_PACKET_TYPE_PROXY,
	}

	packet_bytes, err = proto.Marshal(&packet)

	return
}

func packet_decapsulate(data *[]byte) (packet pb.WebSocketPacketPayload, err error) {

	err = proto.Unmarshal(*data, &packet)

	return
}
