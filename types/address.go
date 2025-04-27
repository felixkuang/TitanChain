// Package types 提供了TitanChain区块链使用的基础类型和数据结构
package types

import (
	"encoding/hex"
	"fmt"
)

// Address 表示一个20字节（160位）的区块链地址
// 这种长度与以太坊地址兼容
type Address [20]uint8

// ToSlice 将地址转换为字节切片
// 返回一个新的切片以防止对原始地址的修改
func (a Address) ToSlice() []byte {
	b := make([]byte, 20)
	for i := 0; i < 20; i++ {
		b[i] = a[i]
	}
	return b
}

// String 返回地址的十六进制字符串表示
// 用于地址的可读性展示
func (a Address) String() string {
	return hex.EncodeToString(a.ToSlice())
}

// AddressFromBytes 从字节切片创建一个新的地址
// 参数：
//   - b: 20字节的输入切片
//
// 如果输入切片长度不是20字节，将会触发panic
func AddressFromBytes(b []byte) Address {
	if len(b) != 20 {
		msg := fmt.Sprintf("given bytes with length %d should be 20", len(b))
		panic(msg)
	}

	var value [20]uint8
	for i := 0; i < 20; i++ {
		value[i] = b[i]
	}

	return Address(value)
}
