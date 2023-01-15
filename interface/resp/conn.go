package resp

// Connection
// redis客户端连接(协议层)
type Connection interface {
	// Write 给客户端写
	Write([]byte) error
	// GetDBIndex 查询客户端用的DB号
	GetDBIndex() int
	// SelectDB 切换DB
	SelectDB(int)
}
