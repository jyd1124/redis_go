package cluster

import (
	"redis_go/interface/resp"
	"redis_go/resp/reply"
)

func FlushDB(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replices := cluster.broadcast(c, cmdArgs)
	var errReply reply.ErrorReply
	for _, r := range replices {
		if reply.IsErrReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
	}
	if errReply == nil {
		return reply.MakeOkReply()
	}
	return reply.MakeErrReply("error : " + errReply.Error())
}
