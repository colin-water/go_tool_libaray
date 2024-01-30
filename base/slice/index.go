package slice

// Index 返回src中和 val 相等的第一个元素下标
// -1 表示没找到
func Index[T comparable](src []T, val T) int {
	return IndexFunc[T](src, func(param T) bool {
		return param == val
	})
}

// IndexFunc 返回 match 返回 true 的第一个下标
// -1 表示没找到
func IndexFunc[T any](src []T, match matchFunc[T]) int {
	for index, value := range src {
		if match(value) {
			return index
		}
	}
	return -1
}

// LastIndex 返回src中和 val 相等的最后一个元素下标
// -1 表示没找到
func LastIndex[T comparable](src []T, val T) int {
	return LastIndexFunc[T](src, func(param T) bool {
		return param == val
	})
}

// LastIndexFunc 返回 match 返回 true 的最后一个元素下标
// -1 表示没找到
func LastIndexFunc[T any](src []T, match matchFunc[T]) int {
	for i := len(src) - 1; i >= 0; i-- {
		if match(src[i]) {
			return i
		}
	}
	return -1
}

// IndexAll 返回和 val 相等的所有元素的下标
func IndexAll[T comparable](src []T, val T) []int {
	return IndexAllFunc[T](src, func(param T) bool {
		return param == val
	})
}

// IndexAllFunc 返回和 match 返回 true 的所有元素的下标（切片）
// 你应该优先使用 IndexAll
func IndexAllFunc[T any](src []T, match matchFunc[T]) []int {
	var indexes = make([]int, 0, len(src)>>3+1)
	for index, value := range src {
		if match(value) {
			indexes = append(indexes, index)
		}
	}
	return indexes
}
