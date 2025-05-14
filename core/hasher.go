// Package core 实现了 TitanChain 区块链的核心功能，包括通用哈希接口和区块、交易哈希实现。
package core

import (
	"crypto/sha256"

	"github.com/felixkuang/titanchain/types"
)

// Hasher 定义了通用的哈希计算接口
// T 是要计算哈希的数据类型
// Hash 计算输入数据的哈希值，返回 types.Hash
type Hasher[T any] interface {
	Hash(T) types.Hash
}

// BlockHasher 实现了区块头的哈希计算
// 使用 SHA256 计算区块头的哈希，不包含交易数据
type BlockHasher struct{}

// Hash 计算区块头的哈希值
// b: 区块头指针
// 返回区块头的 SHA256 哈希
func (BlockHasher) Hash(b *Header) types.Hash {
	h := sha256.Sum256(b.Bytes())
	return types.Hash(h)
}

// TxHasher 实现了交易的哈希计算
// 使用 SHA256 计算交易数据的哈希
type TxHasher struct{}

// Hash 计算交易的哈希值
// tx: 交易指针
// 返回交易数据的 SHA256 哈希
func (TxHasher) Hash(tx *Transaction) types.Hash {
	return types.Hash(sha256.Sum256(tx.Data))
}
