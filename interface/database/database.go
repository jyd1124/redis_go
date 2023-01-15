package database

import (
	"redis_go/interface/resp"
)

// CmdLine is alias for [][]byte, represents a command line
type CmdLine = [][]byte

// Database is the interface for redis style storage engine
type Database interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply
	// AfterClientClose 关闭后的善后工作
	AfterClientClose(c resp.Connection)
	// Close 关闭
	Close()
}

// DataEntity 指代redis数据结构
type DataEntity struct {
	Data interface{}
}
