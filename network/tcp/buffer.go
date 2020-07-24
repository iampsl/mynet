package tcp

// buffer 缓冲 封装切片
type buffer struct {
	splice []byte
}

// Data 取得底层数据切片
func (pbuffer *buffer) Data() []byte {
	return pbuffer.splice
}

//Clear 清空数据
func (pbuffer *buffer) Clear() {
	pbuffer.splice = pbuffer.splice[:0]
}

//Len 数据长度
func (pbuffer *buffer) Len() int {
	return len(pbuffer.splice)
}

//Append 添加数据
func (pbuffer *buffer) Append(data []byte) {
	pbuffer.splice = append(pbuffer.splice, data...)
}
