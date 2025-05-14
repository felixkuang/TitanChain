// Package core 实现了 TitanChain 区块链的核心功能，包括通用的编码/解码接口和 Gob 编解码实现。
package core

import (
	"crypto/elliptic"
	"encoding/gob"
	"io"
)

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
	gob.Register(elliptic.P256())
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
	gob.Register(elliptic.P256())
	return &GobTxDecoder{r: r}
}

// Decode 从 Reader 解码 gob 格式的交易
func (e *GobTxDecoder) Decode(tx *Transaction) error {
	return gob.NewDecoder(e.r).Decode(tx)
}
