// Package core 实现了 TitanChain 区块链的核心功能，包括交易结构、签名、哈希、序列化等。
package core

import (
	"fmt"

	"github.com/felixkuang/titanchain/types"

	"github.com/felixkuang/titanchain/crypto"
)

// Transaction 表示区块链中的一个交易
// 包含原始数据、发起者公钥、签名、哈希
// 支持签名、哈希、验证、序列化等操作
type Transaction struct {
	Data      []byte            // 交易的原始数据
	From      []byte            // 交易发起者的公钥
	Signature *crypto.Signature // 交易的数字签名
	hash      types.Hash        // 交易哈希缓存
}

// NewTransaction 创建一个新的交易实例
// data: 交易数据
func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

// Sign 使用给定的私钥对交易进行签名
// privKey: 用于签名的私钥
// 返回签名过程中可能发生的错误
func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(tx.Data)
	if err != nil {
		return err
	}

	tx.From = privKey.PublicKey().ToSlice()
	tx.Signature = sig

	return nil
}

// Hash 计算并返回交易的哈希值（带缓存）
// hasher: 哈希器实例
func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}
	return tx.hash
}

// Verify 验证交易的签名是否有效
// 检查交易是否包含签名，以及签名是否与公钥和数据匹配
// 返回验证过程中可能发生的错误
func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	publicKey, err := crypto.ToPublicKey(tx.From)
	if err != nil {
		return err
	}

	if !tx.Signature.Verify(publicKey, tx.Data) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}

// Decode 使用指定解码器解码交易
// dec: 解码器实例
func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}

// Encode 使用指定编码器编码交易
// enc: 编码器实例
func (tx *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(tx)
}
