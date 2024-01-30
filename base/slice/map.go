package slice

// toMap 将切片转换为 map
// 键的类型是T也就是切片元素类型，值是空结构体
func toMap[T comparable](src []T) map[T]struct{} {
	resultMap := make(map[T]struct{}, len(src))
	for _, value := range src {
		resultMap[value] = struct{}{} // 使用空结构体,减少内存消耗
	}
	return resultMap
}

// Map 对输入切片 src 中的每个元素应用提供的映射函数 m，得到一个新的切片。
func Map[Src any, Dst any](src []Src, m func(idx int, src Src) Dst) []Dst {
	result := make([]Dst, 0, len(src))
	for i, s := range src {
		result[i] = m(i, s)
	}
	return result
}

// FilterMap 根据提供的映射函数 m 对输入切片 src 进行过滤和映射，得到一个新的切片。
// 如果 m 的第二个返回值是 false，那么我们会忽略第一个返回值
// 即便第二个返回值是 false，后续的元素依旧会被遍历
func FilterMap[Src any, Dst any](src []Src, m func(idx int, src Src) (Dst, bool)) []Dst {
	result := make([]Dst, 0, len(src))
	for i, s := range src {
		dst, ok := m(i, s)
		if ok {
			result = append(result, dst)
		}
	}
	return result
}

// deduplicate 去重，使用上面的toMap
// map去重后重新转成slice
func deduplicate[T comparable](data []T) []T {
	dataMap := toMap(data)
	var newData = make([]T, 0, len(data))
	for key := range dataMap {
		newData = append(newData, key)
	}
	return newData
}

// dev_test
// 对每个元素使用 ContainsFunc 函数来检查是否在新切片 newData 中已经存在。
// type equalFunc[T any] func(src, dst T) bool 是比较两个元素是否相等
func deduplicateFunc[T any](data []T, equal equalFunc[T]) []T {
	newData := make([]T, 0, len(data))
	for index, value := range data {
		if ContainsFunc(data[index+1:], func(src T) bool {
			return equal(src, value)
		}) {
			continue
		}
		newData = append(newData, value)
	}
	return newData
}
