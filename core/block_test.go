// Package core 实现区块链核心功能的测试
package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/felixkuang/titanchain/crypto"
	"github.com/felixkuang/titanchain/types"
	"github.com/stretchr/testify/assert"
)

// TestSignBlock 测试区块签名功能
func TestSignBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(t, 0, types.Hash{})

	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

// TestVerifyBlock 测试区块签名验证功能
func TestVerifyBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(t, 0, types.Hash{})

	assert.Nil(t, b.Sign(privKey))
	assert.Nil(t, b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey().ToSlice()
	// 验证者公钥被篡改，验证应失败
	assert.NotNil(t, b.Verify())

	b.Height = 100
	// 区块高度被篡改，验证应失败
	assert.NotNil(t, b.Verify())
}

func TestDecodeEncodeBlock(t *testing.T) {
	b := randomBlock(t, 1, types.Hash{})
	buf := &bytes.Buffer{}
	assert.Nil(t, b.Encode(NewGobBlockEncoder(buf)))

	bDecode := new(Block)
	assert.Nil(t, bDecode.Decode(NewGobBlockDecoder(buf)))
	assert.Equal(t, b, bDecode)
}

// randomBlock 辅助函数：生成指定高度和前区块哈希的区块
func randomBlock(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	tx := randomTxWithSignature(t)
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	b, err := NewBlock(header, []*Transaction{tx})
	assert.Nil(t, err)
	dataHash, err := CalculateDataHash(b.Transactions)
	assert.Nil(t, err)
	b.Header.DataHash = dataHash
	assert.Nil(t, b.Sign(privKey))

	return b
}
