package parser

import (
	"bufio"
	"errors"
	"io"
	"redis_go/interface/resp"
	"redis_go/lib/logger"
	"redis_go/resp/reply"
	"runtime/debug"
	"strconv"
	"strings"
)

// Payload c -> s 解析完的数据
type Payload struct {
	Data resp.Reply
	Err  error
}

// readState 解析器状态
type readState struct {
	// 单行/多行
	readingMultiLine bool
	// 需要解析的参数数量
	expectedArgsCount int
	msgType           byte
	// 已经解析的参数
	args    [][]byte
	bulkLen int64
}

// 解析器是否解析完成
func (s *readState) finished() bool {
	return s.expectedArgsCount > 0 && len(s.args) == s.expectedArgsCount
}

// ParseStream 返回只可读chan,将解析io结果通过管道输出;代码不会卡在这边(异步)
func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	// 为每个用户生成一个解析器
	go parse0(reader, ch)
	return ch
}

func parse0(reader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	bufReader := bufio.NewReader(reader)
	var state readState
	var err error
	var msg []byte
	for true {
		var ioErr bool
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil { // 出现IO错误
			if ioErr {
				ch <- &Payload{
					Err: err,
				}
				close(ch)
				return
			} // 出现协议错误
			ch <- &Payload{
				Err: err,
			}
			state = readState{}
			continue
		}
		// 判断是不是多行解析模式
		// *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
		if !state.readingMultiLine {
			if msg[0] == '*' { // *3\r\n
				err = parseMultiBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: errors.New("protocol error: " + string(msg)),
					}
					state = readState{}
					continue
				}
				if state.expectedArgsCount == 0 {
					ch <- &Payload{
						Data: &reply.EmptyMultiBulkReply{}, // 空字符串
					}
					state = readState{}
					continue
				}
			} else if msg[0] == '$' { // $4\r\nPING\r\n
				err = parseBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: errors.New("protocol error: " + string(msg)),
					}
					state = readState{}
					continue
				}
				if state.bulkLen == -1 { // $-1\r\n
					ch <- &Payload{
						Data: &reply.NullBulkReply{},
					}
					state = readState{}
					continue
				}
			} else { // + - :
				result, err := parseSingleLineReply(msg)
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = readState{}
				continue
			}
		} else {
			err := readBody(msg, &state)
			if err != nil {
				ch <- &Payload{
					Err: errors.New("protocol error: " + string(msg)),
				}
				state = readState{}
				continue
			}
			if state.finished() {
				var result resp.Reply
				if state.msgType == '*' {
					result = reply.MakeMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					result = reply.MakeBulkReply(state.args[0])
				}
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = readState{}
			}
		}
	}
}

// 根据bulkLen状态,从字节流读取一行(只是按照标识切片)
// *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error
	if state.bulkLen == 0 { // 1.\r\n 切片
		msg, err = bufReader.ReadBytes('\n')
		if err != nil { // 出现IO错误
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' { // 协议错误
			return nil, false, errors.New("protocol error" + string(msg))
		}
	} else { // 2.之前读到了$数字,严格读取字符个数
		msg = make([]byte, state.bulkLen+2)
		_, err := io.ReadFull(bufReader, msg)
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error" + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}

// *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
// *开头,修改解析器的状态
func parseMultiBulkHeader(msg []byte, state *readState) error {
	var err error
	var expectedLine uint64
	expectedLine, err = strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 32)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}
	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		state.msgType = msg[0]                       // *
		state.readingMultiLine = true                // 多行状态
		state.expectedArgsCount = int(expectedLine)  // 3
		state.args = make([][]byte, 0, expectedLine) // 创建参数数组
		return nil
	} else {
		return errors.New("protocol error: " + string(msg))
	}
}

// $4\r\nPING\r\n
// $开头的单行字符串,修改解析器状态
func parseBulkHeader(msg []byte, state *readState) error {
	var err error
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}
	if state.bulkLen == -1 { // null bulk
		return nil
	} else if state.bulkLen > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectedArgsCount = 1
		state.args = make([][]byte, 0, 1)
		return nil
	} else {
		return errors.New("protocol error: " + string(msg))
	}
}

// +OK\r\n -err\r\n :5\r\n
func parseSingleLineReply(msg []byte) (resp.Reply, error) {
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result resp.Reply
	switch msg[0] {
	case '+':
		result = reply.MakeStatusReply(str[1:])
	case '-':
		result = reply.MakeErrReply(str[1:])
	case ':':
		val, err := strconv.ParseInt(str[1:], 10, 64)
		if err != nil {
			return nil, errors.New("protocol error: " + string(msg))
		}
		result = reply.MakeIntReply(val)
	}
	return result, nil
}

// $3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
// PING\r\n
// 读完头部,按照头部信息,解析body数据
func readBody(msg []byte, state *readState) error {
	line := msg[0 : len(msg)-2]
	var err error
	// $3
	if line[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("protocol error: " + string(msg))
		}
		// $0\r\n
		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}
	return nil
}
