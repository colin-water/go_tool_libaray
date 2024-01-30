package slice

// 反转切片
// Reverse 将会完全创建一个新的切片，而不是直接在 src 上进行翻转。
// 传入可比较的类型
func Reverse[T any](src []T) []T {
	result := make([]T, 0, len(src))
	for i := len(src) - 1; i >= 0; i-- {
		result = append(result, src[i])
	}
	return result
}

// ReverseSelf 會直接在 src 上进行翻转。
// 前后指针法
func ReverseSelf[T any](src []T) {
	for left, right := 0, len(src)-1; left < right; {
		src[left], src[right] = src[right], src[left]
		left++
		right--
	}
}
