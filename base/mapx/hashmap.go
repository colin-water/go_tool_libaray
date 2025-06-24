package mapx

import "github.com/colin-water/go_tool_libaray/base/pool"

// Hashable 接口定义了元素需要实现的两个方法：Code 和 Equals。
type Hashable interface {
	// Code 返回该元素的哈希值。
	Code() uint64
	// Equals 比较两个元素是否相等。如果返回 true，那么我们会认为两个键是一样的。
	Equals(key any) bool
}

// node 是 HashMap 中的节点结构，包含 key、value 和指向下一个节点的指针。
// T 是一个类型参数，可以表示任何类型。
// Hashable 是对 T 的约束，表示 T 必须满足 Hashable 接口，即具有 Code 方法和 Equals 方法。
type node[T Hashable, ValType any] struct {
	key   T
	value ValType
	next  *node[T, ValType]
}

// HashMap 是一个哈希映射的实现，使用链表处理冲突。
// 它使用一个底层的哈希表（hashmap）存储键值对，每个键对应一个链表，解决哈希冲突。
// 使用 nodePool 作为节点的对象池，避免频繁创建和销毁节点，提高性能。
type HashMap[T Hashable, ValType any] struct {
	hashmap  map[uint64]*node[T, ValType]  // 用于存储键值对的哈希表，每个键对应一个链表
	nodePool *pool.Pool[*node[T, ValType]] // 用于节点的对象池，避免频繁创建和销毁节点
}

// NewHashMap 创建一个新的哈希映射实例。
func NewHashMap[T Hashable, ValType any](size int) *HashMap[T, ValType] {
	return &HashMap[T, ValType]{
		nodePool: pool.NewPool[*node[T, ValType]](func() *node[T, ValType] {
			return &node[T, ValType]{}
		}),
		hashmap: make(map[uint64]*node[T, ValType], size),
	}
}

// newNode 创建一个新的节点。
func (m *HashMap[T, ValType]) newNode(key T, val ValType) *node[T, ValType] {
	newNode := m.nodePool.Get()
	newNode.value = val
	newNode.key = key
	return newNode
}

// Put 将键值对插入到哈希映射中。
// 第一次出现就自己作为头节点
// 否则使用链表，更新或者在链表后加入新的节点
func (m *HashMap[T, ValType]) Put(key T, val ValType) error {
	// 计算 key 的哈希值
	hash := key.Code()

	// 获取哈希值对应的链表头节点
	root, ok := m.hashmap[hash]
	if !ok {
		newNode := m.newNode(key, val) // 创建新节点
		m.hashmap[hash] = newNode      // 将新节点作为链表头节点
		return nil
	}

	// 遍历链表，查找 key
	pre := root
	for root != nil {
		if root.key.Equals(key) {
			// 如果找到 key，更新值并返回
			root.value = val
			return nil
		}
		pre = root
		root = root.next
	}

	newNode := m.newNode(key, val) // 键不存在，创建新节点
	pre.next = newNode             // 将新节点添加到链表尾部
	return nil
}

// Get 从哈希映射中获取元素。
func (m *HashMap[T, ValType]) Get(key T) (ValType, bool) {
	hash := key.Code()          // 获取键的哈希值
	root, ok := m.hashmap[hash] // 在哈希表中查找对应哈希值的链表头节点
	var val ValType
	if !ok {
		return val, false // 如果链表头节点不存在，说明键不存在
	}
	// 遍历链表
	for root != nil {
		if root.key.Equals(key) {
			return root.value, true // 遍历链表，找到键时返回对应的值
		}
		root = root.next
	}
	return val, false // 遍历完链表未找到键，返回 false
}

// Keys 返回哈希映射里面的所有的 key。
// 注意：key 的顺序是随机的。
// Keys 返回哈希映射中所有键的切片
func (m *HashMap[T, ValType]) Keys() []T {
	res := make([]T, 0)
	// 遍历哈希表中的每个桶（桶是链表的头节点）
	for _, bucketNode := range m.hashmap {
		curNode := bucketNode
		// 遍历链表，将每个节点的键添加到结果切片中
		for curNode != nil {
			res = append(res, curNode.key)
			curNode = curNode.next
		}
	}
	return res
}

// Values 返回哈希映射里面的所有的 value。
// 注意：value 的顺序是随机的。
func (m *HashMap[T, ValType]) Values() []ValType {
	res := make([]ValType, 0)
	for _, bucketNode := range m.hashmap {
		curNode := bucketNode
		for curNode != nil {
			res = append(res, curNode.value)
			curNode = curNode.next
		}
	}
	return res
}

// Delete 从哈希映射中删除元素。
// 第一个返回值为删除 key 的值，第二个是哈希映射是否真的有这个 key。
func (m *HashMap[T, ValType]) Delete(key T) (t ValType, b bool) {
	// 获取哈希值对应的链表头节点
	root, ok := m.hashmap[key.Code()]
	if !ok {
		return t, false // 如果哈希映射中不存在该 key，返回默认值和 false
	}

	pre := root
	num := 0
	for root != nil {
		if root.key.Equals(key) {
			// 找到目标节点，删除操作
			if num == 0 && root.next == nil {
				delete(m.hashmap, key.Code()) // 如果是链表的唯一节点，直接删除链表头
			} else if num == 0 && root.next != nil {
				m.hashmap[key.Code()] = root.next // 如果是链表头，但有下一个节点，更新链表头
			} else {
				pre.next = root.next // 中间或者尾部节点，删除该节点
			}
			val := root.value    // 记录删除的值
			root.formatting()    // 重置节点
			m.nodePool.Put(root) // 将节点放回池中
			return val, true     // 返回删除的值和 true
		}
		num++
		pre = root
		root = root.next
	}

	return t, false // 如果未找到目标节点，返回默认值和 false
}

// formatting 重置节点为空节点
func (n *node[T, ValType]) formatting() {
	var val ValType
	var t T
	n.key = t
	n.value = val
	n.next = nil
}

// Len 返回哈希映射中的元素数量。
func (m *HashMap[T, ValType]) Len() int64 {
	return int64(len(m.hashmap))
}
