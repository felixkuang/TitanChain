// Package network 实现了区块链网络服务器和相关功能
package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/felixkuang/titanchain/types"

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
	Transport     Transport
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
	mempool     *TxPool // 内存交易池
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
		opts.Logger = log.With(opts.Logger, "addr", opts.Transport.Addr())
	}

	chain, err := core.NewBlockchain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}
	s := &Server{
		ServerOpts:  opts,
		chain:       chain,
		mempool:     NewTxPool(1000),
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

	s.boostrapNodes()

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
				if err != core.ErrBlockKnown {
					s.Logger.Log("error", err)
				}
			}
		case <-s.quitCh: // 收到退出信号
			break free
		}
	}

	s.Logger.Log("msg", "Server is shutting down")
}

func (s *Server) boostrapNodes() {
	for _, tr := range s.Transports {
		if s.Transport.Addr() != tr.Addr() {
			if err := s.Transport.Connect(tr); err != nil {
				s.Logger.Log("error", "could not connect to remote", "err", err)
			}
			s.Logger.Log("msg", "connected to remote", "remote", tr.Addr())

			// Send the getStatusMessage so we can sync (if needed)Add commentMore actions
			if err := s.sendGetStatusMessage(tr); err != nil {
				s.Logger.Log("error", "sendGetStatusMessage", "err", err)
			}
		}
	}
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.BlockTime)

	s.Logger.Log("msg", "Starting validator loop", "blockTime", s.BlockTime)

	for {
		<-ticker.C
		err := s.createNewBlock()
		if err != nil {
			s.Logger.Log("error", err).Error()
		}
	}
}

// ProcessMessage 处理解码后的网络消息
// msg: 解码后的消息结构体
// 根据消息类型分发到具体处理逻辑
func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	case *core.Block:
		return s.processBlock(t)
	case *GetStatusMessage:
		return s.processGetStatusMessage(msg.From, t)
	case *StatusMessage:
		return s.processStatusMessage(msg.From, t)
	case *GetBlocksMessage:
		return s.processGetBlocksMessage(msg.From, t)
	}

	return nil
}

func (s *Server) processGetBlocksMessage(from NetAddr, data *GetBlocksMessage) error {
	panic("here")
	fmt.Printf("got get blocks message => %+v\n", data)

	return nil
}

// TODO: Remove the logic from the main function to here
// Normally Transport which is our own transport should do the trick.
func (s *Server) sendGetStatusMessage(tr Transport) error {
	var (
		getStatusMsg = new(GetStatusMessage)
		buf          = new(bytes.Buffer)
	)
	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeGetStatus, buf.Bytes())
	if err := s.Transport.SendMessage(tr.Addr(), msg.Bytes()); err != nil {
		return err
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

// processStatusMessage 处理收到的 StatusMessage 消息
// from: 发送方网络地址
// data: 状态消息内容
// 主要用于节点间同步当前区块高度和节点ID
func (s *Server) processStatusMessage(from NetAddr, data *StatusMessage) error {
	if data.CurrentHeight <= s.chain.Height() {
		s.Logger.Log("msg", "cannot sync blockHeight to low", "ourHeight", s.chain.Height(), "theirHeight", data.CurrentHeight, "addr", from)
		return nil
	}

	// In this case we are 100% sure that the node has blocks heigher than us.
	getBlocksMessage := &GetBlocksMessage{
		From: s.chain.Height(),
		To:   0,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(getBlocksMessage); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeGetBlocks, buf.Bytes())

	return s.Transport.SendMessage(from, msg.Bytes())
}

// processGetStatusMessage 处理收到的 GetStatusMessage 消息
// from: 发送方网络地址
// data: GetStatus 消息内容
// 主要用于响应节点状态请求，返回当前节点的区块高度和ID
func (s *Server) processGetStatusMessage(from NetAddr, data *GetStatusMessage) error {
	fmt.Printf("=> received Getstatus msg from %s => %+v\n", from, data)

	statusMessage := &StatusMessage{
		CurrentHeight: s.chain.Height(),
		ID:            s.ID,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(statusMessage); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeStatus, buf.Bytes())

	// 通过传输层将状态消息发送给请求方
	return s.Transport.SendMessage(from, msg.Bytes())
}

// processBlock 处理收到的区块消息
// b: 区块指针
// 1. 将区块添加到本地区块链
// 2. 异步广播区块到其他节点
func (s *Server) processBlock(b *core.Block) error {
	if err := s.chain.AddBlock(b); err != nil {
		return err
	}

	go func() {
		err := s.broadcastBlock(b)
		if err != nil {
			s.Logger.Log("error", err)
		}
	}()

	return nil
}

// broadcastBlock 将区块编码后广播到所有节点
// b: 区块指针
// 返回广播过程中的错误
func (s *Server) broadcastBlock(b *core.Block) error {
	buf := &bytes.Buffer{}
	if err := b.Encode(core.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeBlock, buf.Bytes())
	return s.broadcast(msg.Bytes())
}

// processTransaction 处理收到的交易，验证签名并加入交易池
// tx: 交易指针
// 1. 检查是否已存在
// 2. 设置首次见到时间戳
// 3. 日志记录并异步广播
// 4. 加入交易池
func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.mempool.Contains(hash) {
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	//s.Logger.Log(
	//	"msg", "adding new tx to mempool",
	//	"hash", hash,
	//	"mempoolPending", s.mempool.PendingCount(),
	//)

	go func() {
		err := s.broadcastTx(tx)
		if err != nil {
			s.Logger.Log("error", err)
		}
	}()

	s.mempool.Add(tx)

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
// 为每个传输层启动 goroutine 持续接收消息并转发到服务器 RPC 通道
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

	// For now we are going to use all transactions that are in the pending pool
	// Later on when we know the internal structure of our transaction
	// we will implement some kind of complexity function to determine how
	// many transactions can be included in a block.
	txx := s.mempool.Pending()

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

	// TODO: pending pool of tx should only reflect on validator nodes.
	// Right now "normal nodes" does not have their pending pool cleared.
	s.mempool.ClearPending()

	go func() {
		err := s.broadcastBlock(block)
		if err != nil {
			s.Logger.Log("error", err)
		}
	}()

	return nil
}

// genesisBlock 生成创世区块
// 返回一个初始化的区块指针，包含默认头部信息
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
