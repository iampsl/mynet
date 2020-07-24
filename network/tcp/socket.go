package tcp

import (
	"net"
	"sync"
	"sync/atomic"
)

//socket 对net.Conn 的包装
type socket struct {
	conn      net.Conn   //TCP底层连接
	buffers   [2]*buffer //双发送缓存
	sendIndex uint       //发送缓存索引
	notify    chan int   //通知通道
	isclose   uint32     //指示socket是否关闭

	data interface{} //自定义数据

	m          sync.Mutex //锁
	bclose     bool       //是否关闭
	writeIndex uint       //插入缓存索引
}

//newsocket 创建一个socket
func newsocket(c net.Conn) *socket {
	if c == nil {
		//c为nil,抛出异常
		panic("c is nil")
	}
	//初始化结构体
	var psocket = new(socket)
	psocket.conn = c
	psocket.buffers[0] = new(buffer)
	psocket.buffers[1] = new(buffer)
	psocket.sendIndex = 0
	psocket.notify = make(chan int, 1)
	psocket.isclose = 0
	psocket.bclose = false
	psocket.writeIndex = 1
	//启动发送协程
	go psocket.dosend()
	return psocket
}

func (my *socket) dosend() {
	writeErr := false
	for {
		_, ok := <-my.notify
		if !ok {
			return
		}
		my.m.Lock()
		my.writeIndex = my.sendIndex
		my.m.Unlock()
		my.sendIndex = (my.sendIndex + 1) % 2
		if !writeErr {
			var sendSplice = my.buffers[my.sendIndex].Data()
			for len(sendSplice) > 0 {
				n, err := my.conn.Write(sendSplice)
				if err != nil {
					writeErr = true
					break
				}
				sendSplice = sendSplice[n:]
			}
		}
		my.buffers[my.sendIndex].Clear()
	}
}

//Read 读数据
func (my *socket) Read(b []byte) (n int, err error) {
	return my.conn.Read(b)
}

//WriteBytes 写数据
func (my *socket) Write(b ...[]byte) {
	my.m.Lock()
	if my.bclose {
		my.m.Unlock()
		return
	}
	dataLen := my.buffers[my.writeIndex].Len()
	writeLen := 0
	for i := 0; i < len(b); i++ {
		writeLen += len(b[i])
		my.buffers[my.writeIndex].Append(b[i])
	}
	if dataLen == 0 && writeLen != 0 {
		my.notify <- 0
	}
	my.m.Unlock()
}

//Close 关闭一个tcpsocket, 释放系统资源
func (my *socket) Close() {
	my.m.Lock()
	if my.bclose {
		my.m.Unlock()
		return
	}
	my.bclose = true
	my.conn.Close()
	close(my.notify)
	my.m.Unlock()
	atomic.StoreUint32(&(my.isclose), 1)
}

//IsClose 判断tcpsocket是否关闭
func (my *socket) IsClose() bool {
	val := atomic.LoadUint32(&(my.isclose))
	if val > 0 {
		return true
	}
	return false
}

//LocalAddr local address
func (my *socket) LocalAddr() string {
	return my.conn.LocalAddr().String()
}

//RemoteAddr remote address
func (my *socket) RemoteAddr() string {
	return my.conn.RemoteAddr().String()
}

//GetData 得到自定义数据
func (my *socket) GetData() interface{} {
	return my.data
}

//SetData 设置自定义数据
func (my *socket) SetData(data interface{}) {
	my.data = data
}
