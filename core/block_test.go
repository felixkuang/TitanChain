// Package core 实现区块链核心功能的测试
package core

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/felixkuang/titanchain/types"
	"github.com/stretchr/testify/assert"
)

// TestHeader_Encode_Decode 测试区块头的编码和解码功能
// 验证区块头可以正确地序列化和反序列化
func TestHeader_Encode_Decode(t *testing.T) {
	// 创建一个测试用的区块头
	h := &Header{
		Version:   1,                     // 设置版本号为1
		PrevBlock: types.RandomHash(),    // 使用随机哈希作为前一个区块的哈希
		Timestamp: time.Now().UnixNano(), // 使用当前时间戳
		Height:    10,                    // 设置区块高度为10
		Nonce:     989394,                // 设置随机数
	}

	// 将区块头编码到缓冲区
	buf := &bytes.Buffer{}
	assert.Nil(t, h.EncodeBinary(buf))

	// 解码并验证结果
	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, h, hDecode)
}

// TestBlock_Encode_Decode 测试完整区块的编码和解码功能
// 验证区块（包括区块头和交易）可以正确地序列化和反序列化
func TestBlock_Encode_Decode(t *testing.T) {
	// 创建一个测试用的区块
	b := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			Timestamp: time.Now().UnixNano(),
			Height:    10,
			Nonce:     989394,
		},
		Transactions: nil, // 空交易列表
	}

	// 将区块编码到缓冲区
	buf := &bytes.Buffer{}
	assert.Nil(t, b.EncodeBinary(buf))

	// 解码并验证结果
	bDecode := &Block{}
	assert.Nil(t, bDecode.DecodeBinary(buf))
	assert.Equal(t, b, bDecode)
}

// TestBlockHash 测试区块哈希计算功能
// 验证区块哈希计算是否正确，以及确保不会生成零哈希
func TestBlockHash(t *testing.T) {
	// 创建一个测试用的区块
	b := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			Timestamp: time.Now().UnixNano(),
			Height:    10,
		},
		Transactions: []Transaction{}, // 空交易列表
	}

	// 计算并验证区块哈希
	h := b.Hash()
	fmt.Println(h)
	assert.False(t, h.IsZero()) // 确保生成的哈希不是零哈希
}
