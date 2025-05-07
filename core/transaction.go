// Package core 实现了TitanChain区块链的核心功能
package core

import (
	"fmt"

	"github.com/felixkuang/titanchain/crypto"
)

// Transaction 表示区块链中的一个交易
// 目前实现了一个简单的数据载荷模型
type Transaction struct {
	Data      []byte            // 交易的原始数据
	From      crypto.PublicKey  // 交易发起者的公钥
	Signature *crypto.Signature // 交易的数字签名
}

// Sign 使用给定的私钥对交易进行签名
// privKey: 用于签名的私钥
// 返回签名过程中可能发生的错误
func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(tx.Data)
	if err != nil {
		return err
	}

	tx.From = privKey.PublicKey()
	tx.Signature = sig

	return nil
}

// Verify 验证交易的签名是否有效
// 检查交易是否包含签名，以及签名是否与公钥和数据匹配
// 返回验证过程中可能发生的错误
func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}
