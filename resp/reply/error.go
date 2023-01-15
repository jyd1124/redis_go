package reply

// UnknownErrReply UNKNOWN_ERR
type UnknownErrReply struct{}

var unknownErrBytes = []byte("-Err unknown\r\n")

func (u *UnknownErrReply) Error() string {
	return "Err unknown"
}

func (u *UnknownErrReply) ToBytes() []byte {
	return unknownErrBytes
}

func MakeUnknownErrReply() *UnknownErrReply {
	return &UnknownErrReply{}
}

// ArgNumErrReply ARG_NUM_ERR
type ArgNumErrReply struct {
	// 用户命令
	Cmd string
}

func (a *ArgNumErrReply) Error() string {
	return "ERR wrong number of arguments for " + a.Cmd + " command"
}

func (a *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for " + a.Cmd + " command\r\n")
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{
		Cmd: cmd,
	}
}

// SyntaxErrReply SYNTAX_ERR
type SyntaxErrReply struct{}

var syntaxErrBytes = []byte("-Err syntax error\r\n")

func (s *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

func (s *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

var theSyntaxErrReply = new(SyntaxErrReply)

func MakeSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}

// WrongTypeErrReply WRONG_TYPE_ERR
type WrongTypeErrReply struct{}

var wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")

func (w *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

func (w *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

var theWrongTypeReply = new(WrongTypeErrReply)

func MakeWrongTypeErrReply() *WrongTypeErrReply {
	return theWrongTypeReply
}

// ProtocloErrReply PROTOCOL_ERR
type ProtocloErrReply struct {
	Msg string
}

func (p *ProtocloErrReply) Error() string {
	return "ERR Protocol error:" + p.Msg
}

func (p *ProtocloErrReply) ToBytes() []byte {
	return []byte("-ERR Protocol error: " + p.Msg + "\r\n")
}

var theProtocolErrReply = new(ProtocloErrReply)

func MakeProtocolErrReply() *ProtocloErrReply {
	return theProtocolErrReply
}
