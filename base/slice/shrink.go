package slice

// 缩容切片
func Shrink[T any](src []T) []T {
	// 得到当前容量和长度
	capacity, length := cap(src), len(src)
	newCapacity, changed := CalCapacity(capacity, length)
	if !changed {
		return src
	}

	newSrc := make([]T, 0, newCapacity) //新建一个长度为0，容量为newC的切片
	newSrc = append(newSrc, src...)
	return newSrc
}

// 比较容量和大小，判断是否需要缩容
func CalCapacity(capacity, length int) (int, bool) {
	if capacity <= 64 { // 如果容量小于等于64，则不进行缩小
		return capacity, false
	}

	// 如果容量小于等于2048且超过了长度的四倍，则缩小为原容量的一半
	if capacity <= 2048 && (capacity/length >= 4) {
		return capacity / 2, true
	}

	// 如果容量大于2048且超过了长度的两倍，则按照0.625倍缩小
	if capacity > 2048 && (capacity/length >= 2) {
		factor := 0.625
		return int(float32(capacity) * float32(factor)), true
	}

	return capacity, false
}
