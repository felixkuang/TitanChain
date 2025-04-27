// Package core 实现了TitanChain区块链的核心功能
package core

import "io"

// Transaction 表示区块链中的一个交易
// 目前实现了一个简单的数据载荷模型
type Transaction struct {
	Data []byte // 交易的原始数据
}

// DecodeBinary 从二进制格式反序列化交易数据
// 当前为占位实现，待完善
func (tx *Transaction) DecodeBinary(r io.Reader) error { return nil }

// EncodeBinary 将交易序列化为二进制格式
// 当前为占位实现，待完善
func (tx *Transaction) EncodeBinary(w io.Writer) error { return nil }
