package c3rliotproxy

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/pion/dtls/v2"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"
)

const (
	UDP_CONNECTION_STATE_CLOSED  = 0
	UDP_CONNECTION_STATE_OPENED  = 1
	UDP_CONNECTION_STATE_OPENING = 2
	UDP_CONNECTION_STATE_CLOSING = 3
)

var udp_remote_connection_connection *dtls.Conn
var udp_remote_connection_context context.Context
var udp_remote_connection_cancel context.CancelFunc
var udp_remote_connection_timeout_ticker *time.Ticker
var udp_remote_connection_state = UDP_CONNECTION_STATE_CLOSED

const UDP_REMOTE_CONNECTION_TIMEOUT_SECONDS = 150
const UDP_REMOTE_CONNECTION_ENDPOINT_SERVER_HOST = "127.0.0.1"
const UDP_REMOTE_CONNECTION_ENDPOINT_SERVER_PORT = 8080
const UDP_REMOTE_CONNECTION_STARTPOINT_SERVER_HOST = "127.0.0.1"
const UDP_REMOTE_CONNECTION_STARTPOINT_SERVER_PORT = 8081

func udp_remote_connection_init(direction string) (err error) {

	if udp_remote_connection_state == UDP_CONNECTION_STATE_OPENING {
		log.Println("remote connection is in the state of opening")
		return
	}

	if udp_remote_connection_state != UDP_CONNECTION_STATE_CLOSED {
		log.Printf("cannot start remote connection with current state: %d\n", udp_remote_connection_state)
		return
	}

	err = udp_remote_connection_open(direction)
	if err != nil {
		return
	}

	err = udp_remote_connection_init_auth()
	if err != nil {
		return
	}

	err = udp_remote_connection_handle_messages()
	if err != nil {
		return
	}

	err = udp_remote_connection_receive_routine()
	if err != nil {
		log.Println("remote connection receive routine failed")
		return
	}
	return
}

func udp_remote_connection_open(direction string) (err error) {

	remote_host := ""
	remote_port := uint(0)
	if direction == "endpoint" {
		remote_host = UDP_REMOTE_CONNECTION_ENDPOINT_SERVER_HOST
		remote_port = UDP_REMOTE_CONNECTION_ENDPOINT_SERVER_PORT
	} else if direction == "startpoint" {
		remote_host = UDP_REMOTE_CONNECTION_STARTPOINT_SERVER_HOST
		remote_port = UDP_REMOTE_CONNECTION_STARTPOINT_SERVER_PORT
	} else {
		return fmt.Errorf("invalid direction: %s", direction)
	}

	udp_remote_connection_state = UDP_CONNECTION_STATE_OPENING

	log.Printf("opening remote connection for %s:%d\n", remote_host, remote_port)

	addr := &net.UDPAddr{IP: net.ParseIP(remote_host), Port: int(remote_port)}

	// Generate a certificate and private key to secure the connection
	certificate, err := selfsign.GenerateSelfSigned()
	if err != nil {
		return
	}

	//
	// Everything below is the pion-DTLS API! Thanks for using it ❤️.
	//

	// Prepare the configuration of the DTLS connection
	config := &dtls.Config{
		Certificates:         []tls.Certificate{certificate},
		InsecureSkipVerify:   true,
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
	}

	// Connect to a DTLS server
	udp_remote_connection_context, udp_remote_connection_cancel = context.WithTimeout(context.Background(), 30*time.Second)

	udp_remote_connection_connection, err = dtls.DialWithContext(udp_remote_connection_context, "udp", addr, config)
	if err != nil {
		return
	}

	log.Printf("opened remote connection for %s:%d\n", remote_host, remote_port)
	udp_remote_connection_state = UDP_CONNECTION_STATE_OPENED
	udp_remote_connection_timeout_ticker = time.NewTicker(UDP_REMOTE_CONNECTION_TIMEOUT_SECONDS * time.Second)

	return
}

func udp_remote_connection_init_auth() (err error) {
	udp_remote_connection_connection.Write([]byte{})
	return
}

func udp_remote_connection_destroy() (err error) {
	log.Println("destroying remote connection")
	err = udp_remote_connection_connection.Close()
	udp_remote_connection_timeout_ticker.Stop()
	udp_remote_connection_state = UDP_CONNECTION_STATE_CLOSED
	return
}

func udp_remote_connection_handle_messages() (err error) {

	go func() {
		log.Println("now handling remote connection messages")
		break_ := false
		for {
			select {

			case _ = <-udp_remote_connection_timeout_ticker.C:
				log.Printf("timeout ticker for remote connection executed")
				udp_remote_connection_destroy()
				break_ = true
				break
			}
			if break_ {
				break
			}
		}
		log.Println("now NOT handling remote connection messages")
	}()
	return
}

func udp_remote_connection_receive_routine() (err error) {

	size := 32 * 1024

	buf := make([]byte, size)

	written := int64(0)

	go func() {
		for {
			nr, er := udp_remote_connection_connection.Read(buf)
			udp_remote_connection_timeout_ticker.Reset(UDP_REMOTE_CONNECTION_TIMEOUT_SECONDS * time.Second)
			if nr > 0 {
				nw, ew := udp_remote_connection_receive_read(buf[:nr])
				if nw < 0 || nr < int(nw) {
					nw = 0
					if ew == nil {
						ew = errors.New("invalid write result")
					}
				}
				written += int64(nw)
				if ew != nil {
					err = ew
					break
				}
				if nr != int(nw) {
					err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}

		}
	}()

	return
}

func udp_remote_connection_receive_read(data []byte) (nw int64, err error) {

	/* decapsulate packet  */
	packet, err := packet_decapsulate(&data)
	if err != nil {
		return
	}
	/**/

	/* get host string string */
	host_bytes := make(net.IP, 4)
	binary.BigEndian.PutUint32(host_bytes, packet.EHost)
	/**/

	return
}

func udp_remote_connection_write(data []byte, host string, port uint) (written int, err error) {

	/* check if remote connection is opened */
	if udp_remote_connection_state == UDP_CONNECTION_STATE_CLOSED {
		log.Printf("connection is closed will restart %d\n", udp_remote_connection_state)
		err = udp_remote_connection_init("endpoint")
		return
	}

	if udp_remote_connection_state == UDP_CONNECTION_STATE_OPENING {
		return
	}
	/**/

	udp_remote_connection_timeout_ticker.Reset(UDP_REMOTE_CONNECTION_TIMEOUT_SECONDS * time.Second)
	return
}
