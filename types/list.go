// Package types 提供了 TitanChain 区块链使用的基础类型和数据结构
package types

import (
	"fmt"
	"reflect"
)

// List 是一个通用的泛型列表结构，支持任意类型元素的动态数组
// 提供插入、删除、查找、清空等常用操作
// T 为元素类型
// 适用于区块链中需要动态管理的对象集合
type List[T any] struct {
	Data []T // 存储元素的切片
}

// NewList 创建一个空的泛型列表
// 返回新建的 List 指针
func NewList[T any]() *List[T] {
	return &List[T]{
		Data: []T{},
	}
}

// Get 获取指定索引的元素
// index: 元素索引
// 返回对应的元素，若越界会 panic
func (l *List[T]) Get(index int) T {
	if index > len(l.Data)-1 {
		err := fmt.Sprintf("the given index (%d) is higher than the length (%d)", index, len(l.Data))
		panic(err)
	}
	return l.Data[index]
}

// Insert 向列表末尾插入一个元素
// v: 要插入的元素
func (l *List[T]) Insert(v T) {
	l.Data = append(l.Data, v)
}

// Clear 清空列表中的所有元素
func (l *List[T]) Clear() {
	l.Data = []T{}
}

// GetIndex 获取指定元素在列表中的索引
// v: 要查找的元素
// 返回元素索引，未找到返回 -1
func (l *List[T]) GetIndex(v T) int {
	for i := 0; i < l.Len(); i++ {
		if reflect.DeepEqual(l.Get(i), v) {
			return i
		}
	}
	return -1
}

// Remove 移除列表中的指定元素（只移除第一个匹配项）
// v: 要移除的元素
func (l *List[T]) Remove(v T) {
	index := l.GetIndex(v)
	if index == -1 {
		return
	}
	l.Pop(index)
}

// Pop 移除指定索引的元素
// index: 要移除的元素索引
func (l *List[T]) Pop(index int) {
	l.Data = append(l.Data[:index], l.Data[index+1:]...)
}

// Contains 判断列表中是否包含指定元素
// v: 要查找的元素
// 返回是否存在
func (l *List[T]) Contains(v T) bool {
	for i := 0; i < len(l.Data); i++ {
		if reflect.DeepEqual(l.Data[i], v) {
			return true
		}
	}
	return false
}

// Last 返回列表中的最后一个元素
// 若列表为空会 panic
func (l List[T]) Last() T {
	return l.Data[l.Len()-1]
}

// Len 返回列表的长度
func (l *List[T]) Len() int {
	return len(l.Data)
}
