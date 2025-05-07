package core

import (
	"testing"

	"github.com/felixkuang/titanchain/crypto"
	"github.com/stretchr/testify/assert"
)

// TestSignTransaction 测试交易签名功能
func TestSignTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)
}

// TestVerifyTransaction 测试交易签名验证功能
func TestVerifyTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.Nil(t, tx.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	tx.From = otherPrivKey.PublicKey()
	// 伪造公钥，验证应失败
	assert.NotNil(t, tx.Verify())
}

// randomTxWithSignature 辅助函数：生成带签名的交易
func randomTxWithSignature(t *testing.T) *Transaction {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}
	assert.Nil(t, tx.Sign(privKey))
	return tx
}
