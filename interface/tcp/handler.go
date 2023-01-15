package tcp

import (
	"context"
	"net"
)

// Handler
// 处理业务逻辑接口
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}
