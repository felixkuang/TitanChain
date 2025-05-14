package network

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/felixkuang/titanchain/core"
	"github.com/stretchr/testify/assert"
)

// TestTxPool 测试交易池初始化和长度
func TestTxPool(t *testing.T) {
	p := NewTxPool()
	assert.Equal(t, 0, p.Len())
}

// TestTxPoolAddTx 测试交易池添加交易、去重和清空
func TestTxPoolAddTx(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("fooo"))
	assert.Nil(t, p.Add(tx))
	assert.Equal(t, 1, p.Len())

	core.NewTransaction([]byte("fooo"))
	assert.Equal(t, 1, p.Len())

	p.Flush()
	assert.Equal(t, 0, p.Len())
}

// TestSortTransactions 测试交易池按 firstSeen 字段排序
func TestSortTransactions(t *testing.T) {
	p := NewTxPool()
	txLen := 1000

	for i := 0; i < txLen; i++ {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		tx.SetFirstSeen(int64(i * rand.Intn(1000000)))
		assert.Nil(t, p.Add(tx))
	}

	assert.Equal(t, txLen, p.Len())

	txx := p.Transactions()
	for i := 0; i < len(txx)-1; i++ {
		assert.True(t, txx[i].FirstSeen() < txx[i+1].FirstSeen())
	}
}
