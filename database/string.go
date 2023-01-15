package database

import (
	"redis_go/interface/database"
	"redis_go/interface/resp"
	"redis_go/lib/utils"
	"redis_go/resp/reply"
)

// GET k1
func execGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeNullBulkReply()
	}
	bytes, _ := entity.Data.([]byte) // 这里假设只有[]byte类型,否则还需要判断类型是否正确
	return reply.MakeBulkReply(bytes)
}

// SET k1 v1
func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	entity := &database.DataEntity{
		Data: val,
	}
	db.PutEntity(key, entity)
	db.addAof(utils.ToCmdLine2("set", args...))
	return reply.MakeOkReply()
}

// SETNX k1 v1
func execSetNx(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	entity := &database.DataEntity{
		Data: val,
	}
	result := db.PutIfAbsent(key, entity)
	db.addAof(utils.ToCmdLine2("setnx", args...))
	return reply.MakeIntReply(int64(result))
}

// GETSET k1 v1
func execGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	entity, exist := db.GetEntity(key)
	db.PutEntity(key, &database.DataEntity{Data: val})
	db.addAof(utils.ToCmdLine2("getset", args...))
	if !exist {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeBulkReply(entity.Data.([]byte))
}

// STRLEN k
func execStrlen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.MakeNullBulkReply()
	}
	bytes := entity.Data.([]byte)
	return reply.MakeIntReply(int64(len(bytes)))
}

func init() {
	RegisterCommand("GET", execGet, 2) // GET K1
	RegisterCommand("SET", execSet, 3) // SET K1 V1
	RegisterCommand("SETNX", execSetNx, 3)
	RegisterCommand("GETSET", execGetSet, 3)
	RegisterCommand("STRLEN", execStrlen, 2)
}
