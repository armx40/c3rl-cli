package c3rliotroxy

import (
	"fmt"
	"log"
	pb "main/c3rl-iot-reverse-roxy/protofiles"
	"net"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

const STARTPOINT_REMOTE_CONNECTION_BUFFER_SIZE = 32 * 1024

type startpoint_connection_id_t string

func (s *startpoint_connection_id_t) bytes() []byte {
	return []byte(*s)
}

func (s *startpoint_connection_id_t) set_from_bytes(b []byte) {
	*s = startpoint_connection_id_t(b)
	return
}

type startpoint_config_host_port_t struct {
	StartPointHost string `yaml:"startPointHost"`
	StartPointPort uint   `yaml:"startPointPort"`
	EndPointHost   string `yaml:"endPointHost"`
	EndPointPort   uint   `yaml:"endPointPort"`
}

type startpoint_config_t struct {
	HostPorts []startpoint_config_host_port_t `yaml:"hostPorts"`
	UID       string                          `yaml:"uid"`
}

/************************************************************************/

type startpoint_remote_connection_t struct {
	Conn           net.Conn
	StartPointPort uint
	StartPointHost string
	EndPointPort   uint
	EndPointHost   string
	ConnOpened     bool
	Time           time.Time
	ConnectionID   startpoint_connection_id_t
	Buffer         []byte
}

func (e *startpoint_remote_connection_t) init() (err error) {

	log.Printf("now handling connection: %s\n", e.ConnectionID)

	/**/
	e.Buffer = make([]byte, STARTPOINT_REMOTE_CONNECTION_BUFFER_SIZE)
	/**/

	/* start read routine */

	err = e.read_routine()

	return
}

func (e *startpoint_remote_connection_t) read_routine() (err error) {

	var n int

	go func() {
		for {
			n, err = e.Conn.Read(e.Buffer)
			if err != nil {
				return
			}

			err = e.process_read(e.Buffer[:n])
			if err != nil {
				return
			}
		}

	}()
	return

}

func (e *startpoint_remote_connection_t) process_read(data []byte) (err error) {

	/* encapsulate packet and send to remote server */
	encapsulated_bytes, err := packet_encapsulate(main_app_startpoint_uid, e.ConnectionID.bytes(), e.EndPointHost, e.EndPointPort, e.StartPointHost, e.StartPointPort, data)
	if err != nil {
		return
	}
	/**/

	/* remote connection send */
	err = remote_connection_write(encapsulated_bytes)
	/**/

	return
}

func (e *startpoint_remote_connection_t) write(data []byte) (n int, err error) {
	return e.Conn.Write(data)
}

/************************************************************************/

type startpoint_connection_t struct {
	Listener       net.Listener
	StartPointPort uint
	StartPointHost string
	EndPointPort   uint
	EndPointHost   string
	ConnOpened     bool
	Time           time.Time
	ConnectionID   startpoint_connection_id_t
}

func (e *startpoint_connection_t) init() (err error) {

	err = e.open()
	if err != nil {
		return
	}

	err = e.receive_connections()
	if err != nil {
		return
	}
	return
}

func (e *startpoint_connection_t) destroy() (err error) {

	e.Listener.Close()

	e.ConnOpened = false

	delete(startpoint.StartpointConnections, e.ConnectionID)

	return
}

func (e *startpoint_connection_t) open() (err error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", e.StartPointHost, e.StartPointPort))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	curr_time := time.Now().UTC()

	/* prepare key */
	// key := fmt.Sprintf("%s:%d-%d", e.Host, e.Port, curr_time.Unix())
	key := e.ConnectionID
	/**/

	startpoint.StartpointConnections[key].ConnectionID = key
	startpoint.StartpointConnections[key].Listener = listener
	startpoint.StartpointConnections[key].ConnOpened = true
	startpoint.StartpointConnections[key].Time = curr_time

	return
}

func (e *startpoint_connection_t) receive_connections() (err error) {

	log.Printf("now listening on %s:%d for %s:%d\n", e.StartPointHost, e.StartPointPort, e.EndPointHost, e.EndPointPort)
	log.Printf("now listening on %s:%d for %s:%d\n", e.StartPointHost, e.StartPointPort, e.EndPointHost, e.EndPointPort)

	if !main_app_no_fmt_log {
		print_text := color.New(color.FgWhite)
		print_text.Printf("%s:", e.StartPointHost)
		print_text.Add(color.FgGreen)
		print_text.Printf("%d", e.StartPointPort)
		print_text.Add(color.FgWhite)
		print_text.Printf(" -> %s:", "endpoint")
		print_text.Add(color.FgGreen)
		print_text.Printf("%d\n", e.EndPointPort)
	}

	go func() {
		for {
			conn, err := e.Listener.Accept()
			if err != nil {
				log.Fatal(err)
				return
			}

			curr_time := time.Now().UTC()

			/* add startpoint remote connection */

			conn_id := startpoint.get_remote_connection_id(e.StartPointHost, e.StartPointPort, e.EndPointHost, e.EndPointPort)

			startpoint_remote_connection := startpoint_remote_connection_t{
				Conn:           conn,
				ConnectionID:   conn_id,
				StartPointPort: e.StartPointPort,
				StartPointHost: e.StartPointHost,
				EndPointPort:   e.EndPointPort,
				EndPointHost:   e.EndPointHost,
				Time:           curr_time,
			}

			startpoint.StartpointRemoteConnections[conn_id] = &startpoint_remote_connection

			startpoint.StartpointRemoteConnections[conn_id].init()
		}
	}()
	return
}

func (e *startpoint_connection_t) read_routine() (err error) {
	return
}

func (e *startpoint_connection_t) write_routine() (err error) {
	return
}

func (e *startpoint_connection_t) write(data []byte) (n int, err error) {
	return
}

/************************************************************************/

type startpoint_t struct {
	StartpointConnections        map[startpoint_connection_id_t]*startpoint_connection_t
	StartpointRemoteConnections  map[startpoint_connection_id_t]*startpoint_remote_connection_t
	StartpointConfig             startpoint_config_t
	StartpointConnectionsCounter uint64
}

func (e *startpoint_t) add_connection(settings *startpoint_config_host_port_t) (err error) {

	conn_id := e.get_connection_id(settings)

	curr_time := time.Now().UTC()

	conn := startpoint_connection_t{
		StartPointHost: settings.StartPointHost,
		StartPointPort: settings.StartPointPort,
		EndPointPort:   settings.EndPointPort,
		EndPointHost:   settings.EndPointHost,
		ConnectionID:   conn_id,
		Time:           curr_time,
	}

	e.StartpointConnections[conn_id] = &conn

	err = e.StartpointConnections[conn_id].init()

	return
}

func (e *startpoint_t) init_config() (err error) {
	// e.StartpointConfig = startpoint_config_t{}
	// /* read config file */

	// yamlFile, err := os.Open(main_app_startpoint_config_file)
	// if err != nil {
	// 	fmt.Println("FATAL: cannot read config file")
	// 	log.Fatalln("cannot read config file")
	// 	return
	// }
	// yamlData, _ := ioutil.ReadAll(yamlFile)

	// err = yaml.Unmarshal([]byte(yamlData), &e.StartpointConfig)

	// if err != nil {
	// 	fmt.Println("FATAL: unable to read contents of config file")
	// 	log.Fatalln("unable to read contents of config file")
	// 	return
	// }
	e.StartpointConfig = main_app_startpoint_config

	return
}

func (e *startpoint_t) init_connections_from_config() (err error) {
	for i := range e.StartpointConfig.HostPorts {
		err = e.add_connection(&e.StartpointConfig.HostPorts[i])
		if err != nil {
			return
		}
	}
	return
}

func (e *startpoint_t) get_remote_connection_id(local_host string, local_port uint, remote_host string, remote_port uint) (conn_id startpoint_connection_id_t) {

	/* temp uuid */
	uuid_str := uuid.NewString()
	/**/
	// conn_id = startpoint_connection_id_t(fmt.Sprintf("r-%s-%s-%d-%s-%d-%d", uuid_str, local_host, local_port, remote_host, remote_port, e.StartpointConnectionsCounter))

	conn_id = startpoint_connection_id_t(fmt.Sprintf("r-%s", uuid_str))
	// e.StartpointConnectionsCounter += 1
	return
}

func (e *startpoint_t) get_connection_id(settings *startpoint_config_host_port_t) (conn_id startpoint_connection_id_t) {
	conn_id = startpoint_connection_id_t(fmt.Sprintf("%s-%d-%s-%d-%d", settings.StartPointHost, settings.StartPointPort, settings.EndPointHost, settings.EndPointPort, e.StartpointConnectionsCounter))
	e.StartpointConnectionsCounter += 1
	return
}

func (e *startpoint_t) init() (err error) {

	log.Println("startpoint init")

	e.StartpointConnectionsCounter = 0

	e.StartpointConnections = make(map[startpoint_connection_id_t]*startpoint_connection_t)
	e.StartpointRemoteConnections = make(map[startpoint_connection_id_t]*startpoint_remote_connection_t)

	err = e.init_config()
	if err != nil {
		return
	}

	err = e.init_connections_from_config()
	if err != nil {
		return
	}

	if !main_app_no_fmt_log {
		print_text := color.New(color.FgWhite)
		print_text.Printf("Connected to: ")
		print_text.Add(color.FgGreen)
		print_text.Printf("%s\n", e.StartpointConfig.UID)
	}

	log.Println("startpoint inited")

	return
}

func (e *startpoint_t) handle_websocket(packet *pb.WebSocketPacketPayload, message *[]byte) (err error) {

	/* get connection id */
	var conn_id startpoint_connection_id_t
	conn_id.set_from_bytes(packet.ConnectionId)
	/**/

	/* find connection */
	conn := e.StartpointRemoteConnections[conn_id]
	/**/

	if conn == nil {
		err = fmt.Errorf("failed to find startpoint remote connection")
		log.Println(err)
		return
	}

	/* write data */
	_, err = conn.write(packet.Data)
	/**/

	return
}

var startpoint startpoint_t

func startpoint_init() (err error) {

	startpoint = startpoint_t{}

	err = startpoint.init()

	return
}
