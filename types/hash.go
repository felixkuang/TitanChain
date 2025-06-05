// Package types 提供了TitanChain区块链使用的基础类型和数据结构
package types

import (
	"encoding/hex"
	"fmt"
)

// Hash 表示一个32字节（256位）的哈希值
type Hash [32]uint8

// IsZero 检查哈希值是否为零值（所有字节都是0）
func (h Hash) IsZero() bool {
	for i := 0; i < 32; i++ {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

// ToSlice 将哈希值转换为字节切片
// 返回一个新的切片以防止对原始哈希的修改
func (h Hash) ToSlice() []byte {
	b := make([]byte, 32)
	for i := 0; i < len(h); i++ {
		b[i] = h[i]
	}
	return b
}

// String 返回哈希值的十六进制字符串表示
func (h Hash) String() string {
	return hex.EncodeToString(h.ToSlice())
}

// HashFromBytes 从字节切片创建一个新的哈希值
// 如果输入切片长度不是32字节，将会触发panic
func HashFromBytes(b []byte) Hash {
	if len(b) != 32 {
		msg := fmt.Sprintf("invalid hash length: %d", len(b))
		panic(msg)
	}
	var value [32]uint8
	for i := 0; i < len(b); i++ {
		value[i] = b[i]
	}
	return Hash(value)
}
