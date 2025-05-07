// Package core 实现了TitanChain区块链的核心功能
package core

import (
	"crypto/sha256"

	"github.com/felixkuang/titanchain/types"
)

// Hasher 定义了通用的哈希计算接口
// T 是要计算哈希的数据类型
type Hasher[T any] interface {
	// Hash 计算输入数据的哈希值
	// 返回计算得到的哈希值
	Hash(T) types.Hash
}

// BlockHasher 实现了区块的哈希计算
type BlockHasher struct{}

// Hash 使用SHA256算法计算区块的哈希值
// 只计算区块头数据的哈希，不包含交易数据
func (BlockHasher) Hash(b *Header) types.Hash {
	h := sha256.Sum256(b.Bytes())
	return types.Hash(h)
}
