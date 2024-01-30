package common

// RealNumber 实数
// 绝大多数情况下，你都应该用这个来表达数字的含义
// 通过使用波浪符，可以表示这些类型的集合。
type RealNumber interface{
~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
~int | ~int8 | ~int16 | ~int32 | ~int64 |
~float32 | ~float64
}

// 所有数字，包含复数
// 在 RealNumber 接口的基础上，包含了复数的表示
type Number interface {
	RealNumber| ~complex64 | ~complex128
}
// 定义新的函数类型
type Comparator[T any] func(src T, dst T) int

//具体的比较器
func ComparatorRealNumber[T RealNumber](numA T, numB T) int {
	if numA < numB {
		return -1
	} else if numA == numB {
		return 0
	} else {
		return 1
	}
}
