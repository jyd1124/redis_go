package cluster

import "redis_go/interface/resp"

func MakeRouter() map[string]CmdFunc {
	routerMap := make(map[string]CmdFunc)
	// 对一个key的操作,可以单独转发
	routerMap["exists"] = defaultFunc // exists k1
	routerMap["type"] = defaultFunc
	routerMap["set"] = defaultFunc
	routerMap["setnx"] = defaultFunc
	routerMap["get"] = defaultFunc
	routerMap["getset"] = defaultFunc
	// 特殊
	routerMap["ping"] = Ping
	routerMap["rename"] = Rename
	routerMap["renamenx"] = Rename
	routerMap["flushdb"] = FlushDB
	routerMap["del"] = Del
	routerMap["select"] = Select

	return routerMap
}

// GET key; SET k1 v1
func defaultFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	key := string(cmdArgs[1])
	peer := cluster.peerPicker.PickNode(key)
	return cluster.relay(peer, c, cmdArgs)
}
