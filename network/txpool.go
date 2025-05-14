package network

import (
	"sort"

	"github.com/felixkuang/titanchain/core"
	"github.com/felixkuang/titanchain/types"
)

// TxMapSorter 用于对交易池中的交易按 firstSeen 字段排序
// 便于优先级处理和批量打包
// transactions: 交易指针切片
// 实现 sort.Interface
type TxMapSorter struct {
	transactions []*core.Transaction
}

// NewTxMapSorter 根据交易哈希映射创建排序器
func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter {
	txx := make([]*core.Transaction, len(txMap))
	i := 0
	for _, tx := range txMap {
		txx[i] = tx
		i++
	}

	s := &TxMapSorter{txx}

	sort.Sort(s)

	return s
}

// Len 返回交易数量
func (s *TxMapSorter) Len() int {
	return len(s.transactions)
}

// Swap 交换两笔交易的位置
func (s *TxMapSorter) Swap(i, j int) {
	s.transactions[i], s.transactions[j] = s.transactions[j], s.transactions[i]
}

// Less 比较两笔交易的 firstSeen 字段，决定排序顺序
func (s *TxMapSorter) Less(i, j int) bool {
	return s.transactions[i].FirstSeen() < s.transactions[j].FirstSeen()
}

// TxPool 表示内存中的交易池，支持并发添加、查询、排序等操作
type TxPool struct {
	transactions map[types.Hash]*core.Transaction // 交易哈希到交易的映射
}

// NewTxPool 创建一个新的交易池实例
func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

// Transactions 返回排序后的所有交易切片
func (p *TxPool) Transactions() []*core.Transaction {
	s := NewTxMapSorter(p.transactions)
	return s.transactions
}

// Add 向交易池添加一笔交易
// 如果已存在则覆盖
func (p *TxPool) Add(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})
	p.transactions[hash] = tx

	return nil
}

// Has 判断交易池中是否存在指定哈希的交易
func (p *TxPool) Has(hash types.Hash) bool {
	_, ok := p.transactions[hash]

	return ok
}

// Len 返回交易池中的交易数量
func (p *TxPool) Len() int {
	return len(p.transactions)
}

// Flush 清空交易池
func (p *TxPool) Flush() {
	p.transactions = make(map[types.Hash]*core.Transaction)
}
