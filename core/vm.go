package core

// Instruction 表示虚拟机支持的指令类型。
type Instruction byte

const (
	// InstrPush 表示将数据压入栈的指令。
	InstrPush Instruction = 0x0a // 10
	// InstrAdd 表示对栈顶两个元素进行加法操作的指令。
	InstrAdd Instruction = 0x0b // 11
	// InstrPushByte 表示将单字节数据压入栈的指令。
	InstrPushByte Instruction = 0x0c
	// InstrPack 表示将多个字节打包为字节切片的指令。
	InstrPack Instruction = 0x0d
	// InstrSub 表示对栈顶两个元素进行减法操作的指令。
	InstrSub Instruction = 0x0e // 14
)

// Stack 表示虚拟机的操作数栈，支持任意类型元素的入栈和出栈操作。
type Stack struct {
	data []any // 存储栈中元素的切片
	sp   int   // 栈顶指针，指向下一个可用位置
}

// NewStack 创建一个指定大小的新栈。
// size: 栈的最大容量。
// 返回 *Stack：新建的栈实例。
func NewStack(size int) *Stack {
	return &Stack{
		data: make([]any, size),
		sp:   0,
	}
}

// Push 将元素 v 压入栈顶。
// v: 任意类型的待入栈元素。
func (s *Stack) Push(v any) {
	s.data[s.sp] = v
	s.sp++
}

// Pop 弹出栈顶元素并返回。
// 返回 any：被弹出的元素。
func (s *Stack) Pop() any {
	value := s.data[0]
	s.data = append(s.data[:0], s.data[1:]...)
	s.sp--

	return value
}

// VM 表示一个简单的字节码虚拟机，用于执行智能合约或脚本。
type VM struct {
	data  []byte // 待执行的字节码数据
	ip    int    // 指令指针，指向当前执行的指令位置
	stack *Stack
}

// NewVM 创建一个新的虚拟机实例，data 为待执行的字节码数据。
func NewVM(data []byte) *VM {
	return &VM{
		data:  data,
		ip:    0,
		stack: NewStack(128),
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
		vm.stack.Push(int(vm.data[vm.ip-1]))

	case InstrPushByte:
		vm.stack.Push(byte(vm.data[vm.ip-1]))

	case InstrPack:
		n := vm.stack.Pop().(int)
		b := make([]byte, n)

		for i := 0; i < n; i++ {
			b[i] = vm.stack.Pop().(byte)
		}

		vm.stack.Push(b)

	case InstrSub:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		c := a - b
		vm.stack.Push(c)

	case InstrAdd:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		c := a + b
		vm.stack.Push(c)
	}

	return nil
}
