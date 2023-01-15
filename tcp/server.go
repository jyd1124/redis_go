package tcp

import (
	"context"
	"net"
	"os"
	"os/signal"
	"redis_go/interface/tcp"
	"redis_go/lib/logger"
	"sync"
	"syscall"
)

// Config
// tcp-server配置
type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	// closeChan
	closeChan := make(chan struct{})
	// 系统发来信号
	sigChan := make(chan os.Signal)
	// runtime会检测收到的这几个信号,然后转发
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info("start listen")
	ListenAndServe(listener, handler, closeChan)
	return nil
}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	// 当该方法没有被自然关闭(线程被强制关闭),通过接收系统信号量来关闭
	go func() {
		<-closeChan
		logger.Info("shutting down")
		_ = listener.Close()
		_ = handler.Close()
	}()
	// 方法正常关闭退出,关闭资源
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()
	// 等待已经启动的携程完成
	var waitDone sync.WaitGroup
	ctx := context.Background()
	for true {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accepted link")
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	// 即使某个携程出现异常,也需要等待运行的携程全部服务完,才可以退出
	waitDone.Wait()
}
