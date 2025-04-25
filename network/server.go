// Package network 实现了区块链网络服务器和相关功能
package network

import (
	"fmt"
	"time"
)

// ServerOptions 定义了服务器的配置选项
type ServerOptions struct {
	Transports []Transport // 传输层实例列表
}

// Server 实现了区块链网络服务器
type Server struct {
	ServerOptions
	rpcCh  chan RPC      // RPC消息通道，用于接收网络消息
	quitCh chan struct{} // 退出信号通道，用于优雅关闭服务器
}

// NewServer 创建一个新的服务器实例
// opts: 服务器配置选项
func NewServer(opts ServerOptions) *Server {
	return &Server{
		ServerOptions: opts,
		rpcCh:         make(chan RPC),         // 创建RPC消息通道
		quitCh:        make(chan struct{}, 1), // 创建带缓冲的退出信号通道
	}
}

// Start 启动服务器并开始处理消息
// 这个方法会阻塞直到服务器收到退出信号
func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(5 * time.Second) // 创建定时器，每5秒触发一次

free:
	for {
		select {
		case rpc := <-s.rpcCh: // 处理接收到的RPC消息
			fmt.Printf("%v\n", rpc)
		case <-s.quitCh: // 收到退出信号
			break free
		case <-ticker.C: // 定时器触发
			fmt.Println("do stuff every x seconds")
		}
	}

	fmt.Println("Server shutdown")
}

// initTransports 初始化所有传输层
// 为每个传输层启动一个goroutine来处理消息
func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() { // 持续从传输层接收消息
				s.rpcCh <- rpc // 将消息转发到服务器的RPC通道
			}
		}(tr)
	}
}
