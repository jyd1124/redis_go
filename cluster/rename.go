package cluster

import (
	"redis_go/interface/resp"
	"redis_go/resp/reply"
)

// Rename rename k1 k2
func Rename(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeErrReply("ERR Wrong number args")
	}
	src := string(cmdArgs[1]) // k1
	dest := string(cmdArgs[2])
	srcPeer := cluster.peerPicker.PickNode(src) // 192.0.0.1:6379
	destPeer := cluster.peerPicker.PickNode(dest)
	if srcPeer != destPeer {
		return reply.MakeErrReply("ERR rename must within on peer")
	}
	return cluster.relay(srcPeer, c, cmdArgs)
}
