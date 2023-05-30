package c3rliotproxy

import "sync"

const REMOTE_CONNECTION_TYPE = "websocket"

var remote_connection_mutex sync.Mutex

func remote_connection_init() (err error) {

	if REMOTE_CONNECTION_TYPE == "websocket" {
		err = websocket_init()
		if err != nil {
			return
		}
	}

	return

}

func remote_connection_write(data []byte) (err error) {

	remote_connection_mutex.Lock()
	defer remote_connection_mutex.Unlock()

	if REMOTE_CONNECTION_TYPE == "websocket" {
		err = websocket_write(data)
		if err != nil {
			return
		}
	}
	return
}
