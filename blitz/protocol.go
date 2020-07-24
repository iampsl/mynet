package blitz

import (
	"encoding/binary"
	"log"
	"mynet/blitz/pbproto"
	"mynet/network/tcp"

	"github.com/golang/protobuf/proto"
)

//Session 会话
type Session interface {
	Write(int, proto.Message)
	Close()
	LocalAddr() string
	RemoteAddr() string
}

//ServerNotify Server通知
type ServerNotify interface {
	OnAccept(Session)
	OnError(Session, error)
	OnMsg(Session, *pbproto.MsgBody)
}

//ClientNotify Client通知
type ClientNotify interface {
	OnConnect(Session, error)
	OnError(Session, error)
	OnMsg(Session, *pbproto.MsgBody)
}

type pbsession struct {
	socket tcp.Socket
}

func (s *pbsession) setSocket(socket tcp.Socket) {
	s.socket = socket
}

func (s *pbsession) getSocket() tcp.Socket {
	return s.socket
}

func (s *pbsession) Write(id int, m proto.Message) {
	data, err := proto.Marshal(m)
	if err != nil {
		log.Println(err)
		return
	}
	var body pbproto.MsgBody
	body.ID = uint32(id)
	body.Data = data
	bodyData, err := proto.Marshal(&body)
	if err != nil {
		log.Println(err)
		return
	}
	var head [4]byte
	binary.BigEndian.PutUint32(head[:], uint32(len(bodyData)))
	s.socket.Write(head[:], bodyData)
}

func (s *pbsession) Close() {
	s.socket.Close()
}

func (s *pbsession) LocalAddr() string {
	return s.socket.LocalAddr()
}

func (s *pbsession) RemoteAddr() string {
	return s.socket.RemoteAddr()
}

//Protocol 协议层
type Protocol struct {
	notify ServerNotify
}

//SetNotify 设置被通知者
func (p *Protocol) SetNotify(n ServerNotify) {
	p.notify = n
}

//OnAccept 接收socket
func (p *Protocol) OnAccept(tcpso tcp.Socket) {
	psession := new(pbsession)
	tcpso.SetData(psession)
	psession.setSocket(tcpso)
	p.notify.OnAccept(psession)
}

//OnError socket出错
func (p *Protocol) OnError(tcpso tcp.Socket, err error) {
	psession := tcpso.GetData().(*pbsession)
	p.notify.OnError(psession, err)
}

//OnRead  socket数据
func (p *Protocol) OnRead(tcpso tcp.Socket, data []byte) int {
	if len(data) < 4 {
		return 0
	}
	length := int(binary.BigEndian.Uint32(data))
	if length > len(data) {
		return 0
	}
	if length < 4 {
		tcpso.Close()
		return length
	}
	var msg pbproto.MsgBody
	err := proto.Unmarshal(data[4:length], &msg)
	if err != nil {
		tcpso.Close()
		return length
	}
	p.notify.OnMsg(tcpso.GetData().(*pbsession), &msg)
	return length
}
