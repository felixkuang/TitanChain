package network

import (
	"sync"

	"github.com/felixkuang/titanchain/core"
	"github.com/felixkuang/titanchain/types"
)

// TxPool 表示交易池，负责管理所有待处理和待打包的交易
// 支持最大容量限制、去重、优先级等功能
// all: 所有交易的有序集合
// pending: 待打包的交易集合
// maxLength: 交易池最大容量，超出时自动移除最早交易
type TxPool struct {
	all       *TxSortedMap // 所有交易的有序集合
	pending   *TxSortedMap // 待打包的交易集合
	maxLength int          // 交易池最大容量
}

// NewTxPool 创建一个新的交易池实例
// maxLength: 交易池最大容量
// 返回新建的 TxPool 指针
func NewTxPool(maxLength int) *TxPool {
	return &TxPool{
		all:       NewTxSortedMap(),
		pending:   NewTxSortedMap(),
		maxLength: maxLength,
	}
}

// Add 向交易池添加一笔交易，自动去重并按容量限制裁剪
// tx: 要添加的交易指针
func (p *TxPool) Add(tx *core.Transaction) {
	// 如果池已满，移除最早的交易
	if p.all.Count() == p.maxLength {
		oldest := p.all.First()
		p.all.Remove(oldest.Hash(core.TxHasher{}))
	}

	if !p.all.Contains(tx.Hash(core.TxHasher{})) {
		p.all.Add(tx)
		p.pending.Add(tx)
	}
}

// Contains 检查指定哈希的交易是否在池中
// hash: 交易哈希
// 返回是否存在
func (p *TxPool) Contains(hash types.Hash) bool {
	return p.all.Contains(hash)
}

// Pending 返回所有待打包的交易切片
func (p *TxPool) Pending() []*core.Transaction {
	return p.pending.txx.Data
}

// ClearPending 清空所有待打包的交易
func (p *TxPool) ClearPending() {
	p.pending.Clear()
}

// PendingCount 返回待打包交易的数量
func (p *TxPool) PendingCount() int {
	return p.pending.Count()
}

// TxSortedMap 是用于管理交易的有序映射结构，支持高效查找、插入、删除
// 内部使用哈希表和泛型列表实现
// lock: 读写锁，保证并发安全
// lookup: 哈希到交易的映射
// txx: 交易列表，保持插入顺序
type TxSortedMap struct {
	lock   sync.RWMutex
	lookup map[types.Hash]*core.Transaction
	txx    *types.List[*core.Transaction]
}

// NewTxSortedMap 创建一个新的 TxSortedMap 实例
func NewTxSortedMap() *TxSortedMap {
	return &TxSortedMap{
		lookup: make(map[types.Hash]*core.Transaction),
		txx:    types.NewList[*core.Transaction](),
	}
}

// First 返回最早插入的交易指针
func (t *TxSortedMap) First() *core.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	first := t.txx.Get(0)
	return t.lookup[first.Hash(core.TxHasher{})]
}

// Get 根据哈希获取交易指针
// h: 交易哈希
// 返回对应的交易指针
func (t *TxSortedMap) Get(h types.Hash) *core.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.lookup[h]
}

// Add 添加一笔交易到映射中，自动去重
// tx: 要添加的交易指针
func (t *TxSortedMap) Add(tx *core.Transaction) {
	hash := tx.Hash(core.TxHasher{})

	t.lock.Lock()
	defer t.lock.Unlock()

	if _, ok := t.lookup[hash]; !ok {
		t.lookup[hash] = tx
		t.txx.Insert(tx)
	}
}

// Remove 移除指定哈希的交易
// h: 交易哈希
func (t *TxSortedMap) Remove(h types.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.txx.Remove(t.lookup[h])
	delete(t.lookup, h)
}

// Count 返回当前映射中的交易数量
func (t *TxSortedMap) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return len(t.lookup)
}

// Contains 检查指定哈希的交易是否存在
// h: 交易哈希
// 返回是否存在
func (t *TxSortedMap) Contains(h types.Hash) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()

	_, ok := t.lookup[h]
	return ok
}

// Clear 清空所有交易
func (t *TxSortedMap) Clear() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.lookup = make(map[types.Hash]*core.Transaction)
	t.txx.Clear()
}
