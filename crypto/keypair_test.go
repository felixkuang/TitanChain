// Package crypto 实现了区块链所需的加密功能的测试
package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestKeypairSignVerifySuccess 测试签名和验证的成功场景
// 验证使用正确的密钥对和消息时，签名验证应该通过
func TestKeypairSignVerifySuccess(t *testing.T) {
	privKey := GeneratePrivateKey()  // 生成新的私钥
	publicKey := privKey.PublicKey() // 获取对应的公钥
	msg := []byte("hello world")     // 测试消息

	// 使用私钥签名消息
	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)

	// 验证签名应该成功
	assert.True(t, sig.Verify(publicKey, msg))
}

// TestKeypairSignVerifyFail 测试签名验证的失败场景
// 验证在使用错误的公钥或消息时，签名验证应该失败
func TestKeypairSignVerifyFail(t *testing.T) {
	privKey := GeneratePrivateKey()  // 生成第一个私钥
	PublicKey := privKey.PublicKey() // 获取第一个公钥
	msg := []byte("hello world")     // 测试消息

	// 使用第一个私钥签名
	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)

	// 生成另一对密钥用于测试
	otherPrivKey := GeneratePrivateKey()
	otherPublicKey := otherPrivKey.PublicKey()

	// 使用错误的公钥验证应该失败
	assert.False(t, sig.Verify(otherPublicKey, msg))
	// 使用错误的消息验证应该失败
	assert.False(t, sig.Verify(PublicKey, []byte("xxxxxx")))
}
