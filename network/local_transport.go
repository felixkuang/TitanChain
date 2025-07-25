// Package network 实现了 TitanChain 区块链的网络传输层，包括本地传输、节点连接、消息发送等功能。
package network

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

// TransportError 表示传输层的自定义错误类型
// 用于描述网络传输过程中的异常和错误
// Message: 错误信息
type TransportError struct {
	Message string
}

// Error 实现 error 接口，返回错误信息
func (e *TransportError) Error() string {
	return e.Message
}

// LocalTransport 实现了用于本地网络通信的 Transport 接口
// 支持点对点连接、消息发送、超时、最大连接数、线程安全等
type LocalTransport struct {
	addr      NetAddr                     // 当前节点的地址
	consumeCh chan RPC                    // 用于接收消息的通道
	lock      sync.RWMutex                // 用于线程安全的读写锁
	peers     map[NetAddr]*LocalTransport // 已连接的对等节点映射
	maxPeers  int                         // 允许的最大对等节点数
	timeout   time.Duration               // 发送操作的超时时间
	running   bool                        // 表示传输层是否正在运行
}

// LocalTransportOpts 包含创建新的 LocalTransport 实例的配置选项
// Addr: 节点地址
// MaxPeers: 最大对等节点数
// Timeout: 超时时间
type LocalTransportOpts struct {
	Addr     NetAddr       // 节点地址
	MaxPeers int           // 最大对等节点数
	Timeout  time.Duration // 超时时间
}

// NewLocalTransport 使用给定的配置选项创建一个新的 LocalTransport 实例
func NewLocalTransport(opts LocalTransportOpts) Transport {
	if opts.MaxPeers == 0 {
		opts.MaxPeers = 100 // 默认最大对等节点数
	}
	if opts.Timeout == 0 {
		opts.Timeout = 5 * time.Second // 默认超时时间
	}

	return &LocalTransport{
		addr:      opts.Addr,
		consumeCh: make(chan RPC, 1024), // 使用带缓冲的通道防止阻塞
		peers:     make(map[NetAddr]*LocalTransport),
		maxPeers:  opts.MaxPeers,
		timeout:   opts.Timeout,
		running:   true,
	}
}

// Consume 返回用于接收消息的通道
func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}

// Connect 与另一个传输节点建立双向连接
// tr: 目标传输层实例
// 返回连接建立过程中的错误
func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.running {
		return &TransportError{"传输层未运行"}
	}

	if len(t.peers) >= t.maxPeers {
		return &TransportError{"已达到最大对等节点数限制"}
	}

	// 类型断言确保连接的是 LocalTransport 类型
	localTr, ok := tr.(*LocalTransport)
	if !ok {
		return &TransportError{"无效的传输层类型"}
	}

	// 防止自连接
	if t.addr == localTr.Addr() {
		return &TransportError{"不能与自身建立连接"}
	}

	// 检查是否已经连接
	if _, exists := t.peers[tr.Addr()]; exists {
		return &TransportError{"已经与该节点建立连接"}
	}

	// 建立连接
	t.peers[tr.Addr()] = localTr
	return nil
}

// SendMessage 向指定的对等节点发送消息，支持超时机制
// to: 目标节点地址
// payload: 消息内容
// 返回发送过程中的错误
func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if !t.running {
		return &TransportError{"传输层未运行"}
	}

	if t.addr == to {
		return &TransportError{"不能向自身发送消息"}
	}

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not send message to unknown peer %s", t.addr, to)
	}

	// 创建消息
	rpc := RPC{
		From:    t.addr,
		Payload: bytes.NewReader(payload),
	}

	// 带超时的发送
	select {
	case peer.consumeCh <- rpc:
		return nil
	case <-time.After(t.timeout):
		return &TransportError{"发送超时"}
	}
}

// Broadcast 向所有已连接的对等节点广播消息
// payload: 要广播的消息内容
// 返回广播过程中遇到的第一个错误（如有），否则返回 nil
func (t *LocalTransport) Broadcast(payload []byte) error {
	for _, peer := range t.peers {
		if err := t.SendMessage(peer.Addr(), payload); err != nil {
			return err
		}
	}
	return nil
}

// Addr 返回传输层的网络地址
func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}

// Disconnect 从传输层的对等节点列表中移除指定节点
// addr: 要断开的对等节点地址
// 返回断开过程中的错误
func (t *LocalTransport) Disconnect(addr NetAddr) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.running {
		return &TransportError{"传输层未运行"}
	}

	if _, exists := t.peers[addr]; !exists {
		return &TransportError{"未找到对等节点"}
	}

	delete(t.peers, addr)
	return nil
}

// Close 关闭传输层，释放所有资源
// 返回关闭过程中的错误
func (t *LocalTransport) Close() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.running {
		return &TransportError{"传输层已经关闭"}
	}

	t.running = false
	close(t.consumeCh)
	t.peers = make(map[NetAddr]*LocalTransport)
	return nil
}

// PeerCount 返回当前连接的对等节点数量
func (t *LocalTransport) PeerCount() int {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return len(t.peers)
}

// IsRunning 返回传输层的运行状态
func (t *LocalTransport) IsRunning() bool {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.running
}

// GetPeers 返回已连接的对等节点地址列表
func (t *LocalTransport) GetPeers() []NetAddr {
	t.lock.RLock()
	defer t.lock.RUnlock()

	peers := make([]NetAddr, 0, len(t.peers))
	for addr := range t.peers {
		peers = append(peers, addr)
	}
	return peers
}
