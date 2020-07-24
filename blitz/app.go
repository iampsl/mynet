package blitz

import (
	"log"
	"mynet/blitz/pbproto"
)

//App 业务层
type App struct {
}

//OnAccept 接收连接
func (app *App) OnAccept(s Session) {
	log.Println(s, s.LocalAddr(), s.RemoteAddr())
}

//OnError 出错处理
func (app *App) OnError(s Session, err error) {
	log.Println(s, s.LocalAddr(), s.RemoteAddr(), err)
}

//OnMsg 消息处理
func (app *App) OnMsg(s Session, msg *pbproto.MsgBody) {

}
