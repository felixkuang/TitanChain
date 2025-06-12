// Package core 实现了 TitanChain 区块链的核心功能，包括通用的编码/解码接口和 Gob 编解码实现。
package core

import (
	"crypto/elliptic"
	"encoding/gob"
	"io"
)

//Add commentMore actions
// For now we GOB encoding is used for fast bootstrapping of the project
// in a later phase I'm considering using Protobuffers as default encoding / decoding.

// Encoder 定义了通用的编码器接口
// T 是要编码的数据类型
// Encode 将数据编码并写入底层 Writer，返回可能的错误
type Encoder[T any] interface {
	Encode(T) error
}

// Decoder 定义了通用的解码器接口
// T 是要解码的数据类型
// Decode 从底层 Reader 读取并解码数据，返回可能的错误
type Decoder[T any] interface {
	Decode(T) error
}

// GobTxEncoder 实现了基于 gob 的交易编码器
// 支持将 Transaction 编码为二进制流
type GobTxEncoder struct {
	w io.Writer
}

// NewGobTxEncoder 创建一个 gob 交易编码器
// w: 输出 Writer
func NewGobTxEncoder(w io.Writer) *GobTxEncoder {
	return &GobTxEncoder{w: w}
}

// Encode 将交易编码为 gob 格式并写入 Writer
func (e *GobTxEncoder) Encode(tx *Transaction) error {
	return gob.NewEncoder(e.w).Encode(tx)
}

// GobTxDecoder 实现了基于 gob 的交易解码器
// 支持从二进制流解码 Transaction
type GobTxDecoder struct {
	r io.Reader
}

// NewGobTxDecoder 创建一个 gob 交易解码器
// r: 输入 Reader
func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	return &GobTxDecoder{r: r}
}

// Decode 从 Reader 解码 gob 格式的交易
func (e *GobTxDecoder) Decode(tx *Transaction) error {
	return gob.NewDecoder(e.r).Decode(tx)
}

// GobBlockEncoder 实现了基于 gob 的区块编码器
// 支持将 Block 编码为二进制流
// w: 输出 Writer
// 用于将区块结构序列化为 gob 格式，便于网络传输或持久化存储
type GobBlockEncoder struct {
	w io.Writer // gob 编码输出的 Writer
}

// NewGobBlockEncoder 创建一个 gob 区块编码器
// w: 输出 Writer
// 返回 GobBlockEncoder 指针
func NewGobBlockEncoder(w io.Writer) *GobBlockEncoder {
	return &GobBlockEncoder{
		w: w,
	}
}

// Encode 将区块编码为 gob 格式并写入 Writer
// b: 待编码的区块指针
// 返回可能的错误
func (enc *GobBlockEncoder) Encode(b *Block) error {
	return gob.NewEncoder(enc.w).Encode(b)
}

// GobBlockDecoder 实现了基于 gob 的区块解码器
// 支持从二进制流解码 Block
// r: 输入 Reader
// 用于将 gob 格式的区块数据反序列化为 Block 结构
type GobBlockDecoder struct {
	r io.Reader // gob 解码输入的 Reader
}

// NewGobBlockDecoder 创建一个 gob 区块解码器
// r: 输入 Reader
// 返回 GobBlockDecoder 指针
func NewGobBlockDecoder(r io.Reader) *GobBlockDecoder {
	return &GobBlockDecoder{
		r: r,
	}
}

// Decode 从 Reader 解码 gob 格式的区块
// b: 解码目标区块指针
// 返回可能的错误
func (dec *GobBlockDecoder) Decode(b *Block) error {
	return gob.NewDecoder(dec.r).Decode(b)
}

func init() {
	gob.Register(elliptic.P256())
}
