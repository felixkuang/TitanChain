package util

import (
	"log"
	"reflect"
)

// AssertEqual 判断两个对象是否深度相等
// a, b: 任意类型的待比较对象
// 若不相等则直接终止程序并输出详细信息
func AssertEqual(a, b any) {
	if !reflect.DeepEqual(a, b) {
		log.Fatalf("ASSERTION: %+v != %+v", a, b)
	}
}
