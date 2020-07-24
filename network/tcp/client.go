package tcp

import (
	"net"
	"time"
)

//Connect 建立连接
func Connect(address string, timeout time.Duration, notify ClientNotify) {
	go doConnect(address, timeout, notify)
}

func doConnect(address string, timeout time.Duration, notify ClientNotify) {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		notify.OnConnect(nil, err)
		return
	}
	psocket := newsocket(conn)
	notify.OnConnect(psocket, nil)
	go readSocket(psocket, notify)
}
