package slice

// UnionSet 并集，只支持 comparable
// 已去重
// 返回值的元素顺序是不定的
func UnionSet[T comparable](sliceA, sliceB []T) []T {
	mapA, mapB := toMap(sliceA), toMap(sliceB)
	// a,b map 合并
	for key := range mapA {
		mapB[key] = struct{}{}
	}
	// 转成slice
	result := make([]T, 0, len(mapB))
	for key := range mapB {
		result = append(result, key)
	}

	return result
}

// UnionSetFunc 并集，支持任意类型
// 直接合并然后去重
func UnionSetFunc[T any](sliceA, sliceB []T, equal equalFunc[T]) []T {
	var ret = make([]T, 0, len(sliceA)+len(sliceB))
	ret = append(ret, sliceA...)
	ret = append(ret, sliceB...)

	return deduplicateFunc(ret, equal)
}
