package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	data := []byte{0x02, 0x0a, 0x02, 0x0a, 0x0b}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())

	assert.Equal(t, byte(4), vm.stack[vm.sp])
}
