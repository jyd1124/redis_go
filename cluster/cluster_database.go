package cluster

import (
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
	"redis_go/config"
	database2 "redis_go/database"
	"redis_go/interface/database"
	"redis_go/interface/resp"
	"redis_go/lib/consistenthash"
	"redis_go/lib/logger"
	"redis_go/resp/reply"
	"strings"
)

type ClusterDatabase struct {
	// 记录自己地址
	self string
	// 整个集群的节点
	nodes []string
	// 节点管理器
	peerPicker *consistenthash.NodeMap
	// 客户端连接池
	peerConnection map[string]*pool.ObjectPool
	// 数据库
	db database.Database
}

func MakeClusterDatabase() *ClusterDatabase {
	cluster := &ClusterDatabase{
		self:           config.Properties.Self,
		db:             database2.NewStandaloneDatabase(),
		peerPicker:     consistenthash.NewNodeMap(nil),
		peerConnection: make(map[string]*pool.ObjectPool),
	}
	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, peer := range config.Properties.Peers {
		nodes = append(nodes, peer)
	}
	nodes = append(nodes, config.Properties.Self)
	cluster.peerPicker.AddNode(nodes...)
	ctx := context.Background()
	for _, peer := range config.Properties.Peers {
		cluster.peerConnection[peer] = pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{Peer: peer})
	}
	cluster.nodes = nodes
	return cluster
}

type CmdFunc func(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply

var router = MakeRouter()

func (cluster *ClusterDatabase) Exec(client resp.Connection, args [][]byte) (result resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			result = reply.MakeUnknownErrReply()
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		return reply.MakeErrReply("not supported cmd")
	}
	result = cmdFunc(cluster, client, args)
	return
}

func (cluster *ClusterDatabase) AfterClientClose(client resp.Connection) {
	cluster.AfterClientClose(client)
}

func (cluster *ClusterDatabase) Close() {
	cluster.db.Close()
}
