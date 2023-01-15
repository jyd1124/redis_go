package database

import (
	"redis_go/interface/resp"
	"redis_go/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}

// PING
func init() {
	RegisterCommand("ping", Ping, 1)
}
