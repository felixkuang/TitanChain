// Package network 实现了 TitanChain 区块链的网络消息与RPC处理，包括消息类型、序列化、RPC处理器等。
package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/felixkuang/titanchain/core"
)

// MessageType 表示网络消息的类型
// 支持交易、区块等类型
type MessageType byte

const (
	MessageTypeTx        MessageType = 0x1
	MessageTypeBlock     MessageType = 0x2
	MessageTypeGetBlocks MessageType = 0x3 // 区块消息

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

// DecodedMessage 表示解码后的网络消息
// From: 消息发送者地址
// Data: 解码后的消息内容（如交易、区块等）
type DecodedMessage struct {
	From NetAddr // 发送者地址
	Data any     // 解码后的消息内容
}

// RPCDecodeFunc 定义了RPC消息解码函数类型
// 输入为RPC消息，输出为解码后的消息结构体和错误
type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

// DefaultRPCDecodeFunc 是默认的RPC消息解码函数
// 支持交易消息的解码，后续可扩展更多类型
// rpc: 输入的RPC消息
// 返回解码后的消息结构体和错误
func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
	}

	logrus.WithFields(logrus.Fields{
		"from": rpc.From,
		"type": msg.Header,
	}).Debug("new incoming message")

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodedMessage{
			From: rpc.From,
			Data: tx,
		}, nil
	case MessageTypeBlock:
		block := new(core.Block)
		if err := block.Decode(core.NewGobBlockDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: block,
		}, nil

	default:
		return nil, fmt.Errorf("invalid message header %x", msg.Header)
	}
}

// RPCProcessor 定义了RPC处理器需要实现的接口
// 目前仅包含消息处理方法，参数为解码后的消息
type RPCProcessor interface {
	ProcessMessage(*DecodedMessage) error
}
