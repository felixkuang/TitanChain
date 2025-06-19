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

// MessageType 表示网络消息的类型枚举
// 用于区分不同类型的网络消息（如交易、区块、状态等）
type MessageType byte

const (
	// MessageTypeTx 交易消息类型
	MessageTypeTx MessageType = 0x1
	// MessageTypeBlock 区块消息类型
	MessageTypeBlock MessageType = 0x2
	// MessageTypeGetBlocks 区块请求消息类型
	MessageTypeGetBlocks MessageType = 0x3
	// MessageTypeStatus 节点状态消息类型
	MessageTypeStatus MessageType = 0x4
	// MessageTypeGetStatus 节点状态请求消息类型
	MessageTypeGetStatus MessageType = 0x5
)

// RPC 表示远程过程调用的消息结构体
// 用于网络层节点间消息传递的统一封装
// 字段说明：
//
//	From: 发送者网络地址
//	Payload: 消息内容（io.Reader，通常为序列化后的消息体）
type RPC struct {
	From    NetAddr   // 发送者的网络地址
	Payload io.Reader // 消息内容
}

// Message 表示网络消息的结构体
// 用于网络层消息的统一封装和序列化
// 字段说明：
//
//	Header: 消息类型（MessageType）
//	Data: 序列化后的消息内容
//
// 支持 gob 序列化/反序列化
// 便于网络传输和类型识别
type Message struct {
	Header MessageType // 消息类型
	Data   []byte      // 消息内容
}

// NewMessage 创建一个新的网络消息实例
// t: 消息类型
// data: 消息内容（已序列化）
// 返回 Message 指针
func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t,
		Data:   data,
	}
}

// Bytes 将消息序列化为字节切片
// 便于网络传输
func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

// DecodedMessage 表示解码后的网络消息结构体
// 用于 RPC 解码后传递给上层处理逻辑
// 字段说明：
//
//	From: 消息发送者地址
//	Data: 解码后的消息内容（如交易、区块、状态等）
type DecodedMessage struct {
	From NetAddr // 发送者地址
	Data any     // 解码后的消息内容
}

// RPCDecodeFunc 定义了 RPC 消息解码函数类型
// 输入为 RPC 消息，输出为解码后的消息结构体和错误
// 便于自定义不同的消息解码逻辑
type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

// DefaultRPCDecodeFunc 是默认的 RPC 消息解码函数
// 支持交易、区块、状态等消息类型的解码
// rpc: 输入的 RPC 消息
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

	case MessageTypeGetStatus:
		return &DecodedMessage{
			From: rpc.From,
			Data: &GetStatusMessage{},
		}, nil

	case MessageTypeStatus:
		statusMessage := new(StatusMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(statusMessage); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: statusMessage,
		}, nil
	default:
		return nil, fmt.Errorf("invalid message header %x", msg.Header)
	}
}

// RPCProcessor 定义了 RPC 处理器需要实现的接口
// 目前仅包含消息处理方法，参数为解码后的消息
// 便于自定义不同的消息处理逻辑
type RPCProcessor interface {
	ProcessMessage(*DecodedMessage) error
}
