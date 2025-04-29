// Package core 实现了TitanChain区块链的核心功能
package core

import "io"

// Encoder 定义了通用的编码器接口
// T 是要编码的数据类型
type Encoder[T any] interface {
	// Encode 将数据编码并写入到io.Writer中
	// 返回可能发生的错误
	Encode(io.Writer, T) error
}

// Decoder 定义了通用的解码器接口
// T 是要解码的数据类型
type Decoder[T any] interface {
	// Decode 从io.Reader中读取数据并解码
	// 返回可能发生的错误
	Decode(io.Reader, T) error
}
