package mapx

// Set 接口定义了集合操作的方法
type Set[T comparable] interface {
	Add(key T)        // 添加元素到集合
	Delete(key T)     // 从集合中删除指定元素
	Exist(key T) bool // 检查集合中是否存在指定元素
	Keys() []T        // 返回集合中所有元素的切片
}

// MapSet 结构体实现了 Set 接口，使用映射来存储集合中的元素
type MapSet[T comparable] struct {
	v map[T]struct{}
}

// NewMapSet 是一个构造函数，用于初始化并返回一个新的 MapSet 实例
// 这里使用了 struct{} 是为了节省内存，因为我们只关心集合中的键而不关心值
// 指定容量 size
func NewMapSet[T comparable](size int) *MapSet[T] {
	return &MapSet[T]{
		v: make(map[T]struct{}, size),
	}
}

// Add 方法将元素添加到映射中，确保元素的唯一性
// struct{} 表示空结构体的类型，而 struct{}{} 表示一个空结构体的实例。
func (s *MapSet[T]) Add(val T) {
	s.v[val] = struct{}{}
}

// Delete 方法从映射中删除指定的元素
func (s *MapSet[T]) Delete(key T) {
	delete(s.v, key)
}

// Exist 方法通过检查映射中是否存在指定的键来判断元素是否存在
func (s *MapSet[T]) Exist(key T) bool {
	_, ok := s.v[key]
	return ok
}

// Keys 方法通过遍历映射中的键，将它们放入切片并返回，但元素顺序不固定
func (s *MapSet[T]) Keys() []T {
	result := make([]T, 0, len(s.v))
	for index := range s.v {
		result = append(result, index)
	}
	return result
}
