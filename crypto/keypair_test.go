// Package crypto 实现了区块链所需的加密功能的测试
// 本文件主要测试密钥对的生成、签名、验证等功能，覆盖正常与异常场景。
package crypto

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestKeypairSignVerifySuccess 测试签名和验证的成功场景
// 用例说明：
//   - 输入：正确的密钥对和消息
//   - 预期：签名验证通过
func TestKeypairSignVerifySuccess(t *testing.T) {
	privKey := GeneratePrivateKey()  // 生成新的私钥
	publicKey := privKey.PublicKey() // 获取对应的公钥
	msg := []byte("hello world")     // 测试消息

	// 使用私钥签名消息
	sig, err := privKey.Sign(msg)
	assert.Nil(t, err, "签名过程中不应出现错误")

	// 验证签名应该成功
	assert.True(t, sig.Verify(publicKey, msg), "签名验证应通过")
}

// TestKeypairSignVerifyFail 测试签名验证的失败场景
// 用例说明：
//   - 输入：错误的公钥或消息
//   - 预期：签名验证失败
func TestKeypairSignVerifyFail(t *testing.T) {
	privKey := GeneratePrivateKey()  // 生成第一个私钥
	publicKey := privKey.PublicKey() // 获取第一个公钥
	msg := []byte("hello world")     // 测试消息

	// 使用第一个私钥签名
	sig, err := privKey.Sign(msg)
	assert.Nil(t, err, "签名过程中不应出现错误")

	// 生成另一对密钥用于测试
	otherPrivKey := GeneratePrivateKey()
	otherPublicKey := otherPrivKey.PublicKey()

	// 使用错误的公钥验证应该失败
	assert.False(t, sig.Verify(otherPublicKey, msg), "使用错误公钥验证应失败")
	// 使用错误的消息验证应该失败
	assert.False(t, sig.Verify(publicKey, []byte("xxxxxx")), "使用错误消息验证应失败")
}

// TestKeypairSignEmptyMessage 测试对空消息的签名与验证
// 用例说明：
//   - 输入：空消息
//   - 预期：签名和验证均应正常工作
func TestKeypairSignEmptyMessage(t *testing.T) {
	privKey := GeneratePrivateKey()
	publicKey := privKey.PublicKey()
	emptyMsg := []byte("")

	sig, err := privKey.Sign(emptyMsg)
	assert.Nil(t, err, "空消息签名不应报错")
	assert.True(t, sig.Verify(publicKey, emptyMsg), "空消息签名验证应通过")
}

// TestKeypairSignTamperedMessage 测试签名内容被篡改的场景
// 用例说明：
//   - 输入：原始消息与被篡改的消息
//   - 预期：原始消息验证通过，篡改后验证失败
func TestKeypairSignTamperedMessage(t *testing.T) {
	privKey := GeneratePrivateKey()
	publicKey := privKey.PublicKey()
	msg := []byte("secure data")
	tamperedMsg := []byte("secure data!")

	sig, err := privKey.Sign(msg)
	assert.Nil(t, err, "签名过程中不应报错")
	assert.True(t, sig.Verify(publicKey, msg), "原始消息验证应通过")
	assert.False(t, sig.Verify(publicKey, tamperedMsg), "篡改消息验证应失败")
}

// TestGeneratePrivateKeyAndPublicKey 测试私钥生成与公钥导出
// 用例说明：
//   - 能否正常生成私钥，并从私钥导出公钥
func TestGeneratePrivateKeyAndPublicKey(t *testing.T) {
	privKey := GeneratePrivateKey()
	assert.NotNil(t, privKey.key, "生成的私钥不应为nil")

	pubKey := privKey.PublicKey()
	assert.NotNil(t, pubKey.Key, "导出的公钥不应为nil")
}

// TestPublicKeyToSliceAndToPublicKey 测试公钥序列化与反序列化（压缩格式）
// 用例说明：
//   - 公钥序列化为压缩字节切片，再反序列化应一致
//   - 错误长度的字节切片反序列化应报错
func TestPublicKeyToSliceAndToPublicKey(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()
	pubBytes := pubKey.ToSlice()
	assert.Equal(t, PublicKeyLength, len(pubBytes), "公钥压缩后长度应为33字节")

	// 反序列化
	newPubKey, err := ToPublicKey(pubBytes)
	assert.Nil(t, err, "公钥反序列化应无错误")
	assert.Equal(t, pubKey.Key.X, newPubKey.Key.X, "反序列化后X坐标应一致")
	assert.Equal(t, pubKey.Key.Y, newPubKey.Key.Y, "反序列化后Y坐标应一致")

	// 错误长度
	_, err = ToPublicKey([]byte{0x01, 0x02, 0x03})
	assert.NotNil(t, err, "错误长度的公钥字节切片应报错")
}

// TestPublicKeyAddress 测试公钥生成地址
// 用例说明：
//   - 公钥生成的地址长度应为20字节，且不同密钥对生成的地址应不同
func TestPublicKeyAddress(t *testing.T) {
	privKey1 := GeneratePrivateKey()
	addr1 := privKey1.PublicKey().Address()
	assert.Equal(t, 20, len(addr1[:]), "地址长度应为20字节")

	privKey2 := GeneratePrivateKey()
	addr2 := privKey2.PublicKey().Address()
	assert.NotEqual(t, addr1, addr2, "不同密钥对生成的地址应不同")
}

// TestSignatureNilFields 测试Signature结构体字段为nil的健壮性
// 用例说明：
//   - S、R为nil时，验证应直接返回false
func TestSignatureNilFields(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()
	msg := []byte("test")

	sig := &Signature{R: &big.Int{}, S: &big.Int{}}
	assert.False(t, sig.Verify(pubKey, msg), "Signature字段为nil时验证应失败")
}
