package slice

// IntersectSet 取交集，已去重
func IntersectSet[T comparable](sliceA, sliceB []T) []T {
	result := make([]T, 0, len(sliceA))
	mapA := toMap(sliceA)
	for _, value := range sliceB {
		if _, ok := mapA[value]; ok {
			result = append(result, value)
		}
	}
	return deduplicate(result)
}

// IntersectSetFunc 支持任意类型
func IntersectSetFunc[T any](sliceA, sliceB []T, equal equalFunc[T]) []T {
	var ret = make([]T, 0, len(sliceA))
	for _, v := range sliceB {
		if ContainsFunc(sliceA, func(t T) bool {
			return equal(t, v)
		}) {
			ret = append(ret, v)
		}
	}
	return deduplicateFunc(ret, equal)
}
