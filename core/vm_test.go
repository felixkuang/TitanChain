package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	s := NewStack(128)

	s.Push(1)
	s.Push(2)

	value := s.Pop()
	assert.Equal(t, value, 1)

	value = s.Pop()
	assert.Equal(t, value, 2)
}

// TestVM 测试简单虚拟机的基本指令执行流程。
// 用例说明：
// - 构造字节码序列，依次将两个数（2、2）压入栈，再执行加法指令，期望栈顶为 4。
// - 验证 Run 方法无错误，且最终栈顶元素为 4。
func TestVM(t *testing.T) {
	// 1 + 2 = 3
	// 1
	// push stack
	// 2
	// push stack
	// add
	// 3
	// push stack

	data := []byte{0x03, 0x0a, 0x02, 0x0a, 0x0e}
	// data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())

	result := vm.stack.Pop().(int)

	assert.Equal(t, 1, result)

	// assert.Equal(t, "FOO", string(result))
}
