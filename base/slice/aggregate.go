package slice

import "github.com/colin-water/go_tool_libaray/base/common"

// Max 返回最大值。
// 该方法假设你至少会传入一个值，确保是数字。
// 在使用 float32 或者 float64 的时候要小心精度问题
func Max[T common.RealNumber](src []T) T {
	res := src[0]
	for _, value := range src {
		if value > res {
			res = value
		}
	}
	return res
}

// Min 返回最小值
// 该方法会假设你至少会传入一个值
// 在使用 float32 或者 float64 的时候要小心精度问题
func Min[T common.RealNumber](src []T) T {
	res := src[0]
	for _, value := range src {
		if res > value {
			res = value
		}
	}
	return res
}

// Sum 求和
// 在使用 float32 或者 float64 的时候要小心精度问题
func Sum[T common.Number](src []T) T {
	var res T
	for _, n := range src {
		res += n
	}
	return res
}
