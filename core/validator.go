package core

import "fmt"

// Validator 定义了区块验证器的接口
type Validator interface {
	// ValidateBlock 验证区块的有效性
	// 返回验证过程中可能发生的错误
	ValidateBlock(*Block) error
}

// BlockValidator 实现了基本的区块验证器
type BlockValidator struct {
	bc *Blockchain // 关联的区块链实例
}

// NewBlockValidator 创建一个新的区块验证器实例
// bc: 要关联的区块链实例
func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

// ValidateBlock 实现了区块的验证逻辑
// b: 要验证的区块
// 验证内容包括：
// 1. 检查区块高度是否已存在
// 2. 验证区块的签名
// 返回验证过程中可能发生的错误
func (v *BlockValidator) ValidateBlock(b *Block) error {
	if v.bc.HasBlock(b.Height) {
		return fmt.Errorf("chain already contains block (%d) with hash (%s)", b.Height, b.Hash(BlockHasher{}))
	}

	if b.Height != v.bc.Height()+1 {
		return fmt.Errorf("block (%s) with height (%d) is too high => current height (%d)", b.Hash(BlockHasher{}), b.Height, v.bc.Height())
	}

	prevHeader, err := v.bc.GetHeader(b.Height - 1)
	if err != nil {
		return err
	}

	hash := BlockHasher{}.Hash(prevHeader)
	if hash != b.PrevBlockHash {
		return fmt.Errorf("the hash of the previous block (%s) is invalid", b.PrevBlockHash)
	}

	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
