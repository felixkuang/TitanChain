package core

// Storage 定义了区块存储的接口
type Storage interface {
	// Put 将区块保存到存储中
	// 返回存储过程中可能发生的错误
	Put(*Block) error
}

// MemoryStore 实现了基于内存的区块存储
type MemoryStore struct {
	// 目前是一个简单的实现，后续会添加实际的存储字段
}

// NewMemorystore 创建一个新的内存存储实例
func NewMemorystore() *MemoryStore {
	return &MemoryStore{}
}

// Put 将区块保存到内存存储中
// 目前是一个空实现，后续会添加实际的存储逻辑
func (s *MemoryStore) Put(b *Block) error {
	return nil
}
