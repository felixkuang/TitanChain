package network

import (
	"fmt"
	"sync"
	"time"
)

// TransportError 表示传输层的自定义错误类型
type TransportError struct {
	Message string
}

func (e *TransportError) Error() string {
	return e.Message
}

// LocalTransport 实现了用于本地网络通信的Transport接口
type LocalTransport struct {
	addr      NetAddr                     // 当前节点的地址
	consumeCh chan RPC                    // 用于接收消息的通道
	lock      sync.RWMutex                // 用于线程安全的读写锁
	peers     map[NetAddr]*LocalTransport // 已连接的对等节点映射
	maxPeers  int                         // 允许的最大对等节点数
	timeout   time.Duration               // 发送操作的超时时间
	running   bool                        // 表示传输层是否正在运行
}

// LocalTransportOpts 包含创建新的LocalTransport实例的配置选项
type LocalTransportOpts struct {
	Addr     NetAddr       // 节点地址
	MaxPeers int           // 最大对等节点数
	Timeout  time.Duration // 超时时间
}

// NewLocalTransport 使用给定的配置选项创建一个新的LocalTransport实例
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
func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.running {
		return &TransportError{"传输层未运行"}
	}

	if len(t.peers) >= t.maxPeers {
		return &TransportError{"已达到最大对等节点数限制"}
	}

	// 类型断言确保连接的是LocalTransport类型
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
func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if !t.running {
		return &TransportError{"传输层未运行"}
	}

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: 无法发送消息到 %s: 未找到对等节点", t.addr, to)
	}

	// 创建消息
	rpc := RPC{
		From:    t.addr,
		Payload: payload,
	}

	// 带超时的发送
	select {
	case peer.consumeCh <- rpc:
		return nil
	case <-time.After(t.timeout):
		return &TransportError{"发送超时"}
	}
}

// Addr 返回传输层的网络地址
func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}

// Disconnect 从传输层的对等节点列表中移除指定节点
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

// Close 关闭传输层
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
