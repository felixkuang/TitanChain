// Package network 实现了 TitanChain 区块链的网络消息与RPC处理，包括消息类型、序列化、RPC处理器等。
package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/felixkuang/titanchain/core"
)

// MessageType 表示网络消息的类型
// 支持交易、区块等类型
type MessageType byte

const (
	MessageTypeTx    MessageType = iota + 1 // 交易消息
	MessageTypeBlock                        // 区块消息
)

// RPC 表示远程过程调用的消息结构
// From: 发送者网络地址
// Payload: 消息内容（io.Reader）
type RPC struct {
	From    NetAddr   // 发送者的网络地址
	Payload io.Reader // 消息内容
}

// Message 表示网络消息的结构体
// Header: 消息类型
// Data: 序列化后的消息数据
type Message struct {
	Header MessageType // 消息类型
	Data   []byte      // 消息内容
}

// NewMessage 创建一个新的网络消息
// t: 消息类型
// data: 消息内容
func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t,
		Data:   data,
	}
}

// Bytes 将消息序列化为字节切片
func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

// RPCHandler 定义了RPC消息处理器接口
// 负责处理收到的RPC消息
type RPCHandler interface {
	HandleRPC(rpc RPC) error
}

// DefaultRPCHandler 是默认的RPC消息处理器
// 负责解码消息并分发到具体处理逻辑
type DefaultRPCHandler struct {
	p RPCProcessor
}

// NewDefaultRPCHandler 创建默认的RPC处理器
// p: RPCProcessor 实例
func NewDefaultRPCHandler(p RPCProcessor) *DefaultRPCHandler {
	return &DefaultRPCHandler{p}
}

// HandleRPC 处理收到的RPC消息
// rpc: RPC消息
// 返回处理过程中的错误
func (p *DefaultRPCHandler) HandleRPC(rpc RPC) error {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
	}

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return err
		}
		return p.p.ProcessTransaction(rpc.From, tx)
	default:
		return fmt.Errorf("invalid message header %x", msg.Header)
	}
}

// RPCProcessor 定义了RPC处理器需要实现的接口
// 目前仅包含交易处理
type RPCProcessor interface {
	ProcessTransaction(NetAddr, *core.Transaction) error
}
