package core

// Instruction 表示虚拟机支持的指令类型。
type Instruction byte

const (
	// InstrPush 表示将数据压入栈的指令。
	InstrPush Instruction = 0x0a // 10
	// InstrAdd 表示对栈顶两个元素进行加法操作的指令。
	InstrAdd Instruction = 0x0b // 11
)

// VM 表示一个简单的字节码虚拟机，用于执行智能合约或脚本。
type VM struct {
	data  []byte // 待执行的字节码数据
	ip    int    // 指令指针，指向当前执行的指令位置
	stack []byte // 操作数栈
	sp    int    // 栈指针，指向当前栈顶
}

// NewVM 创建一个新的虚拟机实例，data 为待执行的字节码数据。
func NewVM(data []byte) *VM {
	return &VM{
		data:  data,
		ip:    0,
		stack: make([]byte, 1024),
		sp:    -1,
	}
}

// Run 启动虚拟机，顺序执行字节码指令，直到结束或遇到错误。
// 返回 error 表示执行过程中遇到的错误。
func (vm *VM) Run() error {
	for {
		instr := Instruction(vm.data[vm.ip])

		if err := vm.Exec(instr); err != nil {
			return err
		}

		vm.ip++

		if vm.ip > len(vm.data)-1 {
			break
		}
	}

	return nil
}

// Exec 执行单条指令，根据指令类型进行相应操作。
// instr: 当前要执行的指令。
// 返回 error 表示执行过程中遇到的错误。
func (vm *VM) Exec(instr Instruction) error {
	switch instr {
	case InstrPush:
		// InstrPush: 将前一个字节压入栈
		vm.pushStack(vm.data[vm.ip-1])
	case InstrAdd:
		// InstrAdd: 取出栈顶两个元素相加，并将结果压入栈
		a := vm.stack[0]
		b := vm.stack[1]
		c := a + b
		vm.pushStack(c)
	}

	return nil
}

// pushStack 将一个字节压入操作数栈。
// b: 需要压入栈的字节。
func (vm *VM) pushStack(b byte) {
	vm.sp++
	vm.stack[vm.sp] = b
}
