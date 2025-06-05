package core

import (
	"testing"

	"github.com/go-kit/log"

	"github.com/felixkuang/titanchain/types"
	"github.com/stretchr/testify/assert"
)

// TestAddBlock 测试区块链批量添加区块的功能和高度统计
func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	lenBlocks := 1000
	for i := 0; i < lenBlocks; i++ {
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks+1)
	// 添加高度已存在的区块应返回错误
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 89, types.Hash{})))
}

// TestNewBlockchain 测试区块链初始化和创世区块
func TestNewBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

// TestHasBlock 测试区块存在性判断
func TestHasBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.True(t, bc.HasBlock(0))
	assert.False(t, bc.HasBlock(1))
	assert.False(t, bc.HasBlock(100))
}

// TestGetHeader 测试区块头的获取和正确性
func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenBlocks := 1000

	for i := 0; i < lenBlocks; i++ {
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(block.Height)
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}
}

// TestAddBlockToHigh 测试添加高度跳跃的区块
func TestAddBlockToHigh(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.Nil(t, bc.AddBlock(randomBlock(t, 1, getPrevBlockHash(t, bc, uint32(1)))))
	// 跳跃高度添加区块应返回错误
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 3, types.Hash{})))
}

// newBlockchainWithGenesis 辅助函数：创建带创世区块的区块链
func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(log.NewNopLogger(), randomBlock(t, 0, types.Hash{}))
	assert.Nil(t, err)

	return bc
}

// getPrevBlockHash 辅助函数：获取前一区块的哈希
func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevHeader)
}
