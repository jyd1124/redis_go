package resp

// Reply S -> C的回复
type Reply interface {
	// ToBytes 回复内容转换成字节
	ToBytes() []byte
}
