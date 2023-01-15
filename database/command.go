package database

import "strings"

// 记录系统里所有的指令
var cmdTable = make(map[string]*command)

// 命令结构体
type command struct {
	// 命令执行函数
	exector ExecFunc
	// 需要的参数个数
	arity int
}

// RegisterCommand 后续添加命令,扩展cmd集
func RegisterCommand(name string, exector ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		exector: exector,
		arity:   arity,
	}
}
