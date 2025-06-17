package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestState_Basic 测试 State 的基本存取、删除功能。
func TestState_Basic(t *testing.T) {
	state := NewState()

	// 测试 Put 和 Get 正常流程
	key := []byte("foo")
	value := []byte("bar")
	assert.Nil(t, state.Put(key, value))

	got, err := state.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, value, got)

	// 测试 Delete
	assert.Nil(t, state.Delete(key))
	_, err = state.Get(key)
	assert.NotNil(t, err)
}

// TestState_GetNotFound 测试获取不存在的 key。
func TestState_GetNotFound(t *testing.T) {
	state := NewState()
	_, err := state.Get([]byte("not_exist"))
	assert.Error(t, err)
}

// TestState_Overwrite 测试同一 key 多次 Put 是否覆盖。
func TestState_Overwrite(t *testing.T) {
	state := NewState()
	key := []byte("k1")
	v1 := []byte("v1")
	v2 := []byte("v2")
	assert.Nil(t, state.Put(key, v1))
	assert.Nil(t, state.Put(key, v2))
	got, err := state.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, v2, got)
}
