package tcp

//Socket 套接字
type Socket interface {
	Write(...[]byte)      //写数据
	Close()               //关闭socket
	LocalAddr() string    //本地地址
	RemoteAddr() string   //远程地址
	SetData(interface{})  //设置自定义数据
	GetData() interface{} //得到自定义数据
}

//ServerNotify Server通知
type ServerNotify interface {
	OnAccept(Socket)           //接收到新的连接
	OnError(Socket, error)     //对应连接出错
	OnRead(Socket, []byte) int //对应连接读到数据
}

//ClientNotify 客户端通知
type ClientNotify interface {
	OnConnect(Socket, error)   //建立连接回调
	OnError(Socket, error)     //对应连接出错
	OnRead(Socket, []byte) int //对应连接读到数据
}
