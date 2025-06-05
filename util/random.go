package util

import (
	"math/rand"
	"testing"

	"github.com/felixkuang/titanchain/core"
	"github.com/felixkuang/titanchain/crypto"
	"github.com/felixkuang/titanchain/types"
	"github.com/stretchr/testify/assert"
)

// RandomBytes 生成指定长度的随机字节切片
// size: 字节长度
// 返回生成的随机字节切片
func RandomBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

// RandomHash 生成一个随机的 32 字节哈希值
// 返回 types.Hash 类型的随机哈希
func RandomHash() types.Hash {
	return types.HashFromBytes(RandomBytes(32))
}

// NewRandomTransaction 创建一个未签名的随机交易
// size: 交易数据长度
// 返回新的 Transaction 实例
func NewRandomTransaction(size int) *core.Transaction {
	return core.NewTransaction(RandomBytes(size))
}

// NewRandomTransactionWithSignature 创建一个带签名的随机交易
// t: 测试对象
// privKey: 用于签名的私钥
// size: 交易数据长度
// 返回已签名的 Transaction 实例
func NewRandomTransactionWithSignature(t *testing.T, privKey crypto.PrivateKey, size int) *core.Transaction {
	tx := NewRandomTransaction(size)
	assert.Nil(t, tx.Sign(privKey))
	return tx
}
