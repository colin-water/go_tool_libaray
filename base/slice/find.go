package slice

// Find 查找元素
// 如果没有找到，第二个返回值返回 false
func Find[T any](src []T, match matchFunc[T]) (T, bool) {
	for _, value := range src {
		if match(value) {
			return value, true
		}
	}
	var zero T
	return zero, false
}

// FindAll 查找所有符合条件的元素
// 永远不会返回 nil
func FindAll[T any](src []T, match matchFunc[T]) []T {
	// 符合条件元素应该是少数
	// 所以会除以 8
	// 也就是触发扩容的情况下，最多三次就会和原本的容量一样
	// +1 是为了保证，至少有一个元素
	result := make([]T, 0, len(src)>>3+1)

	for _, value := range src {
		if match(value) {
			result = append(result, value)
		}
	}
	return result
}
