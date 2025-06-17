package core

import "fmt"

// State 表示区块链的简单状态存储，采用内存中的键值对映射实现。
type State struct {
	data map[string][]byte // 存储状态数据的 map，key 为字符串，value 为字节切片
}

// NewState 创建一个新的状态存储实例。
// 返回 *State：新建的状态存储对象。
func NewState() *State {
	return &State{
		data: make(map[string][]byte),
	}
}

// Put 向状态存储中写入一对键值。
// k: 键（字节切片），v: 值（字节切片）。
// 返回 error：写入过程中遇到的错误，正常返回 nil。
func (s *State) Put(k, v []byte) error {
	s.data[string(k)] = v

	return nil
}

// Delete 从状态存储中删除指定键。
// k: 待删除的键（字节切片）。
// 返回 error：删除过程中遇到的错误，正常返回 nil。
func (s *State) Delete(k []byte) error {
	delete(s.data, string(k))

	return nil
}

// Get 根据键从状态存储中获取对应的值。
// k: 待查询的键（字节切片）。
// 返回 value: 查询到的值（字节切片）；error: 若未找到则返回错误。
func (s *State) Get(k []byte) ([]byte, error) {
	key := string(k)

	value, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("given key %s not found", key)
	}

	return value, nil
}
