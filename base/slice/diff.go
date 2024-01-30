package slice

// DiffSet 差集，只支持 comparable 类型
// 已去重
// 并且返回值的顺序是不确定的
// a转成map，使用map的delete删除b中出现的元素
// 返回sliceA - sliceB
func DiffSet[T comparable](sliceA, sliceB []T) []T {
	mapA := toMap(sliceA)
	for _, val := range sliceB {
		delete(mapA, val)
	}

	result := make([]T, 0, len(mapA))
	for key := range mapA {
		result = append(result, key)
	}

	return result
}

// DiffSetFunc 差集，已去重
// 你应该优先使用 DiffSet
func DiffSetFunc[T any](sliceA, sliceB []T, equal equalFunc[T]) []T {
	result := make([]T, 0, len(sliceA))
	for _, val := range sliceA {
		//遍历sliceB，判定是否存在元素与当前val一致
		if !ContainsFunc(sliceB, func(src T) bool {
			return equal(src, val)
		}) {
			result = append(result, val)
		}
	}
	return deduplicateFunc(result, equal)
}
