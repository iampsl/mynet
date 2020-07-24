package tcp

import (
	"log"
	"net"
)

//ListenAndAccept listen and accept
func ListenAndAccept(address string, notify ServerNotify) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		} else {
			psocket := newsocket(tcpConn)
			notify.OnAccept(psocket)
			go readSocket(psocket, notify)
		}
	}
}

//socketnotify socket通知
type socketnotify interface {
	OnError(Socket, error)
	OnRead(Socket, []byte) int
}

func readSocket(psocket *socket, notify socketnotify) {
	defer psocket.Close()
	readbuffer := make([]byte, 1024)
	readsize := 0
	for {
		if readsize == len(readbuffer) {
			pnew := make([]byte, 2*len(readbuffer))
			copy(pnew, readbuffer)
			readbuffer = pnew
		}
		n, err := psocket.Read(readbuffer[readsize:])
		if err != nil {
			notify.OnError(psocket, err)
			break
		}
		readsize += n
		procTotal := 0
		for {
			if psocket.IsClose() {
				procTotal = readsize
				break
			}
			proc := notify.OnRead(psocket, readbuffer[procTotal:readsize])
			if proc == 0 {
				break
			}
			procTotal += proc
		}
		if procTotal > 0 {
			copy(readbuffer, readbuffer[procTotal:readsize])
			readsize -= procTotal
		}
	}
}
