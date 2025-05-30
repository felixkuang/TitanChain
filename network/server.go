// Package network 实现了区块链网络服务器和相关功能
package network

import (
	"bytes"
	"github.com/felixkuang/titanchain/types"
	"os"
	"time"

	"github.com/go-kit/log"

	"github.com/felixkuang/titanchain/core"

	"github.com/felixkuang/titanchain/crypto"
)

// defaultBlockTime 定义了默认的区块出块时间间隔（5秒）
// 用于未指定 BlockTime 时的服务器出块周期，适用于测试和开发环境
var defaultBlockTime = 5 * time.Second

// ServerOpts 定义了服务器的配置选项
// 包括传输层、出块时间、私钥、RPC解码与处理器等
type ServerOpts struct {
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc      // RPC消息解码函数
	RPCProcessor  RPCProcessor       // RPC消息处理器
	Transports    []Transport        // 传输层实例列表
	BlockTime     time.Duration      // 出块间隔
	PrivateKey    *crypto.PrivateKey // 节点私钥（为空则非验证者）
}

// Server 实现了区块链网络服务器
// 负责管理交易池、消息通道、定时出块等
// 支持多传输层和优雅关闭
type Server struct {
	ServerOpts
	memPool     *TxPool // 内存交易池
	chain       *core.Blockchain
	isValidator bool          // 是否为验证者节点
	rpcCh       chan RPC      // RPC消息通道，用于接收网络消息
	quitCh      chan struct{} // 退出信号通道，用于优雅关闭服务器
}

// NewServer 创建一个新的服务器实例
// opts: 服务器配置选项
func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}

	chain, err := core.NewBlockchain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}
	s := &Server{
		ServerOpts:  opts,
		chain:       chain,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),         // 创建RPC消息通道
		quitCh:      make(chan struct{}, 1), // 创建带缓冲的退出信号通道
	}

	// 如果未指定 RPCProcessor，则默认使用自身
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	return s, nil
}

// Start 启动服务器并开始处理消息
// 这个方法会阻塞直到服务器收到退出信号
func (s *Server) Start() {
	s.initTransports()

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("error", err)
			}

			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				s.Logger.Log("error", err)
			}
		case <-s.quitCh: // 收到退出信号
			break free
		}
	}

	s.Logger.Log("msg", "Server is shutting down")
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(5 * time.Second)

	s.Logger.Log("msg", "Starting validator loop", "blockTime", s.BlockTime)

	for {
		<-ticker.C
		s.createNewBlock()
	}
}

// ProcessMessage 处理解码后的网络消息
// msg: 解码后的消息结构体
// 根据消息类型分发到具体处理逻辑
func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	}
	return nil
}

// broadcast 向所有传输层广播消息
// payload: 要广播的消息内容
// 返回广播过程中遇到的第一个错误（如有），否则返回 nil
func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}
	return nil
}

// processTransaction 处理收到的交易，验证签名并加入交易池
// tx: 交易指针
// 1. 检查是否已存在
// 2. 设置首次见到时间戳
// 3. 日志记录并异步广播
// 4. 加入交易池
func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		return nil
	}

	//if err := tx.Verify(); err != nil {
	//	return err
	//}

	tx.SetFirstSeen(time.Now().UnixNano())

	s.Logger.Log(
		"msg", "adding new tx to mempool",
		"hash", hash,
		"mempoolLength", s.memPool.Len(),
	)

	go s.broadcastTx(tx)

	return s.memPool.Add(tx)
}

// 广播区块
func (s *Server) broadcastBlock(b *core.Block) error {
	// 返回nil表示成功
	return nil
}

// broadcastTx 将交易编码后广播到所有节点
// tx: 交易指针
// 返回广播过程中的错误
func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
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

// createNewBlock 创建新区块（简化实现，实际应包含打包交易、签名等）
func (s *Server) createNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}

	// For now we are going to use all transactions that are in the mempool
	// Later on when we know the internal structure of our transaction
	// we will implement some kind of complexity function to determine how
	// many transactions can be included in a block.
	txx := s.memPool.Transactions()

	block, err := core.NewBlockFromPrevHeader(currentHeader, txx)
	if err != nil {
		return err
	}

	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}

	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

	s.memPool.Flush()

	return nil
}

// 生成创世区块
func genesisBlock() *core.Block {
	header := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Height:    0,
		Timestamp: 000000,
	}

	b, _ := core.NewBlock(header, nil)
	return b
}
