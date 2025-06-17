package core

import (
	"fmt"
	"sync"

	"github.com/go-kit/log"
)

// Blockchain 表示区块链的核心数据结构
type Blockchain struct {
	logger    log.Logger
	store     Storage // 区块存储接口
	lock      sync.RWMutex
	headers   []*Header // 所有区块头的有序列表
	validator Validator // 区块验证器
}

// NewBlockchain 创建一个新的区块链实例
// genesis: 创世区块
// 返回新创建的区块链实例和可能发生的错误
func NewBlockchain(l log.Logger, genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		headers: []*Header{},
		store:   NewMemorystore(),
		logger:  l,
	}
	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)

	return bc, err
}

// SetValidator 设置区块链的验证器
// v: 要使用的验证器实例
func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

// AddBlock 添加新的区块到链中
// b: 要添加的区块
// 在添加之前会进行验证，返回可能发生的错误
func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	for _, tx := range b.Transactions {
		bc.logger.Log("msg", "executing code", "len", len(tx.Data), "hash", tx.Hash(&TxHasher{}))

		vm := NewVM(tx.Data)
		if err := vm.Run(); err != nil {
			return err
		}

		bc.logger.Log("vm result", vm.stack.data[vm.stack.sp])
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[height], nil
}

// HasBlock 检查指定高度的区块是否存在
// height: 要检查的区块高度
// 返回是否存在该高度的区块
func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

// Height 获取当前区块链的高度
// 返回最新区块的高度
func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return uint32(len(bc.headers) - 1)
}

// addBlockWithoutValidation 在不进行验证的情况下添加区块
// b: 要添加的区块
// 返回存储过程中可能发生的错误
func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	defer bc.lock.Unlock()

	bc.logger.Log(
		"msg", "new block",
		"hash", b.Hash(BlockHasher{}),
		"height", b.Height,
		"transactions", len(b.Transactions),
	)

	return bc.store.Put(b)
}
