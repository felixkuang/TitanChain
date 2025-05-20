// Package network 实现了区块链网络服务器和相关功能
package network

import (
	"fmt"
	"time"

	"github.com/felixkuang/titanchain/core"
	"github.com/sirupsen/logrus"

	"github.com/felixkuang/titanchain/crypto"
)

// defaultBlockTime 定义了默认的区块出块时间间隔（5秒）
// 用于未指定 BlockTime 时的服务器出块周期，适用于测试和开发环境
var defaultBlockTime = 5 * time.Second

// ServerOpts 定义了服务器的配置选项
// 包括传输层、出块时间、私钥等
type ServerOpts struct {
	RPCHandler RPCHandler
	Transports []Transport        // 传输层实例列表
	BlockTime  time.Duration      // 出块间隔
	PrivateKey *crypto.PrivateKey // 节点私钥（为空则非验证者）
}

// Server 实现了区块链网络服务器
// 负责管理交易池、消息通道、定时出块等
// 支持多传输层和优雅关闭
type Server struct {
	ServerOpts
	blockTime   time.Duration // 出块间隔
	memPool     *TxPool       // 内存交易池
	isValidator bool          // 是否为验证者节点
	rpcCh       chan RPC      // RPC消息通道，用于接收网络消息
	quitCh      chan struct{} // 退出信号通道，用于优雅关闭服务器
}

// NewServer 创建一个新的服务器实例
// opts: 服务器配置选项
func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}

	s := &Server{
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),         // 创建RPC消息通道
		quitCh:      make(chan struct{}, 1), // 创建带缓冲的退出信号通道
	}

	if opts.RPCHandler == nil {
		opts.RPCHandler = NewDefaultRPCHandler(s)
	}

	s.ServerOpts = opts

	return s
}

// Start 启动服务器并开始处理消息
// 这个方法会阻塞直到服务器收到退出信号
func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.blockTime) // 创建定时器，每5秒触发一次

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			if err := s.RPCHandler.HandleRPC(rpc); err != nil {
				logrus.Error(err)
			}
		case <-s.quitCh: // 收到退出信号
			break free
		case <-ticker.C: // 定时器触发
			if s.isValidator {
				s.createNewBlock()
			}
		}
	}
	fmt.Println("Server shutdown")
}

// handleTransaction 处理收到的交易，验证签名并加入交易池
func (s *Server) ProcessTransaction(from NetAddr, tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash": hash,
		}).Info("transaction already in mempool")

		return nil
	}

	//if err := tx.Verify(); err != nil {
	//	return err
	//}

	tx.SetFirstSeen(time.Now().UnixNano())

	logrus.WithFields(logrus.Fields{
		"hash":           hash,
		"mempool length": s.memPool.Len(),
	}).Info("adding new tx to the mempool")

	return s.memPool.Add(tx)
}

// createNewBlock 创建新区块（简化实现，实际应包含打包交易、签名等）
func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
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
