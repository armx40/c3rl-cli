package c3rliotroxy

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"

	pb "main/c3rl-iot-reverse-roxy/protofiles"

	"github.com/fatih/color"
)

const ENDPOINT_CONNECTION_TIMEOUT_SECONDS = 5 * 60 // 5 minutes

const ENDPOINT_CONNECTION_BUFFER_SIZE = 32 * 1024

const ENDPOINT_DEADLINE_ENABLED = true

type endpoint_connection_id_t string

func (s endpoint_connection_id_t) bytes() []byte {
	return []byte(s)
}

func (s endpoint_connection_id_t) set_from_bytes(b []byte) {
	s = endpoint_connection_id_t(b)
	return
}

type endpoint_connection_t struct {
	Conn          net.Conn
	EPort         uint
	EHost         string
	ConnOpened    bool
	Time          time.Time
	ConnectionID  endpoint_connection_id_t
	Endpoint      *endpoint_t
	Buffer        []byte
	StartPointUID string
}

func (e *endpoint_connection_t) init() (err error) {

	/**/
	e.Buffer = make([]byte, ENDPOINT_CONNECTION_BUFFER_SIZE)
	/**/

	err = e.open()
	if err != nil {
		return
	}

	err = e.read_routine()
	if err != nil {
		return
	}

	return
}

func (e *endpoint_connection_t) destroy() (err error) {

	e.Conn.Close()

	e.ConnOpened = false

	delete(endpoint.EndpointConnections, e.ConnectionID)

	log.Printf("destroyed endpoint connection for host: %s and port %d with conn id: %s, current count: %d\n", e.EHost, e.EPort, e.ConnectionID, len(endpoint.EndpointConnections))

	return
}

func (e *endpoint_connection_t) open() (err error) {

	log.Printf("opening endpoint for host: %s and port %d with id: %s\n", e.EHost, e.EPort, e.ConnectionID)

	print_text := color.New(color.FgWhite)
	print_text.Printf("opening endpoint for host: %s:", e.EHost)
	print_text.Add(color.FgGreen)
	print_text.Printf("%d", e.EPort)
	print_text.Add(color.FgWhite)
	print_text.Printf(" with id: ")
	print_text.Add(color.FgBlue)
	print_text.Printf("%s\n", e.ConnectionID)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", e.EHost, e.EPort))
	if err != nil {
		return err
	}

	curr_time := time.Now().UTC()

	/* prepare key */
	key := e.ConnectionID
	/**/

	/* set timeout. why?? */
	if ENDPOINT_DEADLINE_ENABLED {
		conn.SetDeadline(time.Now().Add(ENDPOINT_CONNECTION_TIMEOUT_SECONDS * time.Second))
	}
	/**/

	endpoint.EndpointConnections[key].ConnectionID = key //??
	endpoint.EndpointConnections[key].Conn = conn
	endpoint.EndpointConnections[key].ConnOpened = true
	endpoint.EndpointConnections[key].Time = curr_time

	log.Printf("opened endpoint for host: %s and port %d\n", e.EHost, e.EPort)

	return
}

func (e *endpoint_connection_t) read_routine() (err error) {

	log.Printf("starting endpoint read routine for host: %s and port %d current count: %d\n", e.EHost, e.EPort, len(endpoint.EndpointConnections))

	var n int

	go func() {
		for {

			n, err = e.Conn.Read(e.Buffer)
			if err != nil {
				e.destroy()
				return
			}

			/* reset deadline */
			if ENDPOINT_DEADLINE_ENABLED {
				e.Conn.SetDeadline(time.Now().Add(ENDPOINT_CONNECTION_TIMEOUT_SECONDS * time.Second))
			}
			/**/

			err = e.process_read(e.Buffer[:n])
		}
	}()
	return
}

func (e *endpoint_connection_t) process_read(data []byte) (err error) {

	/* encapsulate packet and send to remote server */
	encapsulated_bytes, err := packet_encapsulate(e.StartPointUID, e.ConnectionID.bytes(), e.EHost, e.EPort, "127.0.0.1", 0, data)
	if err != nil {
		return
	}
	/**/

	/* remote connection send */
	err = remote_connection_write(encapsulated_bytes)
	/**/

	return
}

func (e *endpoint_connection_t) process_packet(packet *pb.WebSocketPacketPayload) (n int, err error) {
	n, err = e.Conn.Write(packet.Data)
	if err != nil {
		e.destroy()
		return
	}
	/* reset deadline */
	if ENDPOINT_DEADLINE_ENABLED {
		e.Conn.SetDeadline(time.Now().Add(ENDPOINT_CONNECTION_TIMEOUT_SECONDS * time.Second))
	}
	/**/
	return
}

/************************************************************************/

type endpoint_t struct {
	EndpointConnections        map[endpoint_connection_id_t]*endpoint_connection_t
	EndpointConnectionsCounter uint64
}

func (e *endpoint_t) init() (err error) {

	log.Println("endpoint init")

	e.EndpointConnectionsCounter = 0

	e.EndpointConnections = make(map[endpoint_connection_id_t]*endpoint_connection_t)

	log.Println("endpoint inited")

	return
}

func (e *endpoint_t) handle_websocket(packet *pb.WebSocketPacketPayload, message *[]byte) (err error) {

	/* check connection id and perform */
	err = e.process_packet(packet)
	/**/

	return
}

func (e *endpoint_t) process_packet(packet *pb.WebSocketPacketPayload) (err error) {

	/* first get connection  */
	e_conn := e.get_connection_from_packet(packet)
	/**/

	/* if connection exist then forward the packet */
	if e_conn != nil {
		_, err = e_conn.process_packet(packet)
		return
	}
	/**/

	/* connection doesnt exist create a new one */
	err = e.add_connection_from_packet(packet)
	if err != nil {
		return
	}
	/**/

	/* again get connection */
	e_conn = e.get_connection_from_packet(packet)
	/**/

	/* for some reason if the connection doesnt exist then return error */
	if e_conn == nil {
		err = fmt.Errorf("connection doesnt exist even after creation")
		return
	}
	/**/

	/* process the first write */
	_, err = e_conn.process_packet(packet)
	return
	/**/
}

func (e *endpoint_t) get_connection_from_packet(packet *pb.WebSocketPacketPayload) (e_conn *endpoint_connection_t) {

	conn_id := endpoint_connection_id_t(packet.ConnectionId)

	e_conn, ok := e.EndpointConnections[conn_id]
	if !ok {
		e_conn = nil
		return
	}
	return
}

func (e *endpoint_t) add_connection_from_packet(packet *pb.WebSocketPacketPayload) (err error) {

	curr_time := time.Now().UTC()

	host := make(net.IP, 4)
	binary.BigEndian.PutUint32(host, packet.EHost)

	conn_id := endpoint_connection_id_t(packet.ConnectionId)

	conn := endpoint_connection_t{
		ConnectionID:  conn_id,
		Time:          curr_time,
		EHost:         host.To4().String(),
		EPort:         uint(packet.EPort),
		StartPointUID: string(packet.StartpointUid),
	}

	e.EndpointConnections[conn_id] = &conn

	err = e.EndpointConnections[conn_id].init()

	return
}

var endpoint endpoint_t

func endpoint_init() (err error) {

	endpoint = endpoint_t{}

	err = endpoint.init()

	return
}
