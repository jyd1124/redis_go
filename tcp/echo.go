package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
	"redis_go/lib/logger"
	"redis_go/lib/sync/atomic"
	"redis_go/lib/sync/wait"
	"sync"
	"time"
)

// EchoClient
// 客户端
type EchoClient struct {
	Conn net.Conn
	// 包装waitGroup(添加超时),会自动初始化
	Waiting wait.Wait
}

func (client *EchoClient) Close() error {
	client.Waiting.WaitWithTimeout(10 * time.Second)
	_ = client.Conn.Close()
	return nil
}

// EchoHandler
// 业务引擎
type EchoHandler struct {
	// 业务引擎的client(set)
	activeConn sync.Map
	// 判断是否正在关闭(存在并发问题,利用原子bool)
	closing atomic.Boolean
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	// 业务引擎在关闭
	if handler.closing.Get() {
		_ = conn.Close()
	}
	client := &EchoClient{
		Conn: conn,
	}
	handler.activeConn.Store(client, struct{}{})
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		// 接收出现err
		if err != nil {
			// 连接关闭
			if err == io.EOF {
				logger.Info("Connecting close")
				handler.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		// 业务完成才会关闭
		client.Waiting.Add(1)
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}
}

func (handler *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	handler.closing.Set(true)
	handler.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})
	return nil
}
