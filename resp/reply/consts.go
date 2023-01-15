package reply

// PongReply PING
type PongReply struct{}

var pongBytes = []byte("+PONG\r\n")

func (r PongReply) ToBytes() []byte {
	return pongBytes
}

var thePongReply = new(PongReply)

func MakePongReply() *PongReply {
	return thePongReply
}

// OkReply OK
type OkReply struct{}

var okBytes = []byte("+OK\r\n")

func (r OkReply) ToBytes() []byte {
	return okBytes
}

var theOkReply = new(OkReply)

func MakeOkReply() *OkReply {
	return theOkReply
}

// NullBulkReply NULL (空)
type NullBulkReply struct{}

var nullBulkBytes = []byte("$-1\r\n")

func (r NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

var theNullBulkReply = new(NullBulkReply)

func MakeNullBulkReply() *NullBulkReply {
	return theNullBulkReply
}

// EmptyMultiBulkReply EMPTY_ARRAY (空数组)
type EmptyMultiBulkReply struct{}

var emptyMultiBulkBytes = []byte("*0\r\n")

func (r EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

var theEmptyMultiBulkReply = new(EmptyMultiBulkReply)

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return theEmptyMultiBulkReply
}

// NoReply NO (真空,什么都没有)
type NoReply struct{}

var noBytes = []byte("")

func (r NoReply) ToBytes() []byte {
	return noBytes
}

var theNoReply = new(NoReply)

func MakeNoReply() *NoReply {
	return theNoReply
}
