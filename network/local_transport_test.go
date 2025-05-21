package network

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试创建新的LocalTransport实例
func TestNewLocalTransport(t *testing.T) {
	opts := LocalTransportOpts{
		Addr:     "node1",
		MaxPeers: 5,
		Timeout:  2 * time.Second,
	}

	tr := NewLocalTransport(opts).(*LocalTransport)
	if tr.addr != opts.Addr {
		t.Errorf("期望地址为 %s，实际得到 %s", opts.Addr, tr.addr)
	}
	if tr.maxPeers != opts.MaxPeers {
		t.Errorf("期望最大对等节点数为 %d，实际得到 %d", opts.MaxPeers, tr.maxPeers)
	}
	if tr.timeout != opts.Timeout {
		t.Errorf("期望超时时间为 %v，实际得到 %v", opts.Timeout, tr.timeout)
	}
}

// 测试LocalTransport的连接功能
func TestLocalTransportConnect(t *testing.T) {
	tr1 := NewLocalTransport(LocalTransportOpts{Addr: "node1"}).(*LocalTransport)
	tr2 := NewLocalTransport(LocalTransportOpts{Addr: "node2"}).(*LocalTransport)

	// 测试成功连接
	err := tr1.Connect(tr2)
	if err != nil {
		t.Errorf("期望连接成功，但得到错误: %v", err)
	}

	err = tr2.Connect(tr1)
	if err != nil {
		t.Errorf("期望连接成功，但得到错误: %v", err)
	}

	assert.Equal(t, tr1.peers[tr2.Addr()], tr2)
	assert.Equal(t, tr2.peers[tr1.Addr()], tr1)

	// 测试自连接
	err = tr1.Connect(tr1)
	if err == nil {
		t.Error("期望自连接时返回错误")
	}

	// 测试重复连接
	err = tr1.Connect(tr2)
	if err == nil {
		t.Error("期望重复连接时返回错误")
	}
}

// 测试LocalTransport的消息发送功能
func TestLocalTransportSendMessage(t *testing.T) {
	tr1 := NewLocalTransport(LocalTransportOpts{Addr: "node1"})
	tr2 := NewLocalTransport(LocalTransportOpts{Addr: "node2"})

	// 连接传输层
	err := tr1.Connect(tr2)
	if err != nil {
		t.Fatalf("连接传输层失败: %v", err)
	}

	err = tr2.Connect(tr1)
	if err != nil {
		t.Fatalf("连接传输层失败: %v", err)
	}

	// 测试发送消息
	msg := []byte("hello world")
	assert.Nil(t, tr1.SendMessage(tr2.Addr(), msg))

	// 测试接收消息
	rpc := <-tr2.Consume()
	b, err := ioutil.ReadAll(rpc.Payload)

	assert.Nil(t, err)
	assert.Equal(t, b, msg)
	assert.Equal(t, rpc.From, tr1.Addr())
}

func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport(LocalTransportOpts{Addr: "A"})
	trb := NewLocalTransport(LocalTransportOpts{Addr: "B"})
	trc := NewLocalTransport(LocalTransportOpts{Addr: "C"})

	tra.Connect(trb)
	tra.Connect(trc)

	msg := []byte("foo")
	assert.Nil(t, tra.(*LocalTransport).Broadcast(msg))

	rpcb := <-trb.Consume()
	b, err := ioutil.ReadAll(rpcb.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)

	rpcc := <-trc.Consume()
	b, err = ioutil.ReadAll(rpcc.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)

}

// 测试LocalTransport的断开连接功能
func TestLocalTransportDisconnect(t *testing.T) {
	tr1 := NewLocalTransport(LocalTransportOpts{Addr: "node1"}).(*LocalTransport)
	tr2 := NewLocalTransport(LocalTransportOpts{Addr: "node2"}).(*LocalTransport)

	// 连接传输层
	err := tr1.Connect(tr2)
	if err != nil {
		t.Fatalf("连接传输层失败: %v", err)
	}

	// 测试断开连接
	err = tr1.Disconnect(tr2.Addr())
	if err != nil {
		t.Errorf("断开连接失败: %v", err)
	}

	// 验证对等节点数量
	if tr1.PeerCount() != 0 {
		t.Errorf("期望对等节点数为 0，实际得到 %d", tr1.PeerCount())
	}

	// 测试向已断开连接的节点发送消息
	err = tr1.SendMessage(tr2.Addr(), []byte("test"))
	if err == nil {
		t.Error("期望向已断开连接的节点发送消息时返回错误")
	}
}

// 测试LocalTransport的关闭功能
func TestLocalTransportClose(t *testing.T) {
	tr := NewLocalTransport(LocalTransportOpts{Addr: "node1"}).(*LocalTransport)

	// 测试关闭传输层
	err := tr.Close()
	if err != nil {
		t.Errorf("关闭传输层失败: %v", err)
	}

	// 验证传输层已停止运行
	if tr.IsRunning() {
		t.Error("期望传输层停止运行")
	}

	// 测试对已关闭的传输层进行操作
	tr2 := NewLocalTransport(LocalTransportOpts{Addr: "node2"})
	err = tr.Connect(tr2)
	if err == nil {
		t.Error("期望对已关闭的传输层进行连接操作时返回错误")
	}

	err = tr.SendMessage("node2", []byte("test"))
	if err == nil {
		t.Error("期望对已关闭的传输层进行发送操作时返回错误")
	}
}

// 测试LocalTransport的对等节点操作
func TestLocalTransportPeerOperations(t *testing.T) {
	tr1 := NewLocalTransport(LocalTransportOpts{Addr: "node1"}).(*LocalTransport)
	tr2 := NewLocalTransport(LocalTransportOpts{Addr: "node2"}).(*LocalTransport)
	tr3 := NewLocalTransport(LocalTransportOpts{Addr: "node3"}).(*LocalTransport)

	// 连接多个对等节点
	tr1.Connect(tr2)
	tr1.Connect(tr3)

	// 测试对等节点数量
	if count := tr1.PeerCount(); count != 2 {
		t.Errorf("期望对等节点数为 2，实际得到 %d", count)
	}

	// 测试获取对等节点列表
	peers := tr1.GetPeers()
	if len(peers) != 2 {
		t.Errorf("期望有 2 个对等节点，实际得到 %d 个", len(peers))
	}

	// 验证对等节点地址
	peerMap := make(map[NetAddr]bool)
	for _, peer := range peers {
		peerMap[peer] = true
	}
	if !peerMap[tr2.Addr()] || !peerMap[tr3.Addr()] {
		t.Error("缺少预期的对等节点地址")
	}
}
