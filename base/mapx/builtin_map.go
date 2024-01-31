package mapx

// simpleMap 是对 map 的二次封装
type simpleMap[K comparable, V any] struct {
	data map[K]V
}

// newSimpleMap 创建一个新的 simpleMap 实例
func newSimpleMap[K comparable, V any](capacity int) *simpleMap[K, V] {
	return &simpleMap[K, V]{
		data: make(map[K]V, capacity),
	}
}

func newSimpleMapWithData[K comparable, V any](data map[K]V) *simpleMap[K, V] {
	return &simpleMap[K, V]{data: data}
}

// Put 将键值对插入到 map 中
func (b *simpleMap[K, V]) Put(key K, val V) error {
	b.data[key] = val
	return nil
}

// Get 根据键获取对应的值，并返回是否存在
func (b *simpleMap[K, V]) Get(key K) (V, bool) {
	val, ok := b.data[key]
	return val, ok
}

// Delete 根据键删除 map 中的键值对，并返回被删除的值及是否存在
func (b *simpleMap[K, V]) Delete(k K) (V, bool) {
	v, ok := b.data[k]
	delete(b.data, k)
	return v, ok
}

// Keys 返回 map 中所有键的切片，顺序是随机的
// 即便对于同一个实例，调用两次，得到的结果都可能不同
// 调用 map的Keys方法
func (b *simpleMap[K, V]) Keys() []K {
	return Keys[K, V](b.data)
}

// Values 返回 map 中所有值的切片
// 调用 map的Values方法
func (b *simpleMap[K, V]) Values() []V {
	return Values[K, V](b.data)
}

// Len 返回 map 中键值对的数量
func (b *simpleMap[K, V]) Len() int64 {
	return int64(len(b.data))
}
