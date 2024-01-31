package mapx

import "fmt"

// Keys 返回 map 里面的所有的 key。
// 需要注意：这些 key 的顺序是随机。
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// Values 返回 map 里面的所有的 value。
// 需要注意：这些 value 的顺序是随机。
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

// KeysValues 返回 map 里面的所有的 key,value。
// 需要注意：这些 (key,value) 的顺序是随机,相对顺序是一致的。
func KeysValues[K comparable, V any](m map[K]V) ([]K, []V) {
	keys := make([]K, 0, len(m))
	values := make([]V, 0, len(m))
	for key, value := range m {
		keys = append(keys, key)
		values = append(values, value)
	}
	return keys, values
}

// ToMapWithKeyValues 将会返回一个map[K]V
// keys 与 values 的长度必须相同，且不为 nil。
// 如果 keys 或 values 为 nil，则函数返回错误。
// 如果 keys 中存在相同的元素，则在返回的 map 中，相同的 key 对应的 value 为最后出现的对应 value。
func ToMapWithKeyValues[K comparable, V any](keys []K, values []V) (map[K]V, error) {
	if keys == nil || values == nil {
		return nil, fmt.Errorf("keys与values均不可为nil")
	}
	n := len(keys)
	if n != len(values) {
		return nil, fmt.Errorf("keys与values的长度不同, len(keys)=%d, len(values)=%d", n, len(values))
	}
	result := make(map[K]V, n)
	for index, key := range keys {
		result[key] = values[index]
	}
	return result, nil

}
