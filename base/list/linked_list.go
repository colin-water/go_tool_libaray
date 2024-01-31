package list

import "github.com/yishengzhishui/library/base/common"

// node 双向循环链表结点
type node[T any] struct {
	prev *node[T]
	next *node[T]
	val  T
}

// LinkedList 双向循环链表
type LinkedList[T any] struct {
	head   *node[T] // 头结点
	tail   *node[T] // 尾结点
	length int      // 链表长度
}

var (
	_ List[any] = &LinkedList[any]{}
)

// NewLinkedList 创建一个双向循环链表
func NewLinkedList[T any]() *LinkedList[T] {
	// 创建头结点和尾结点，并使其相互指向，形成循环链表
	head := &node[T]{}
	tail := &node[T]{}
	head.next, head.prev = tail, tail
	tail.next, tail.prev = head, head
	return &LinkedList[T]{
		head:   head,
		tail:   tail,
		length: 0,
	}
}

// NewLinkedListOf 将切片转换为双向循环链表
func NewLinkedListOf[T any](data []T) *LinkedList[T] {
	list := NewLinkedList[T]()
	if err := list.Append(data...); err != nil {
		panic(err)
	}
	return list
}

func (l *LinkedList[T]) findNode(index int) *node[T] {
	var curNode *node[T]
	// 如果索引在链表长度的一半以内，从链表头部开始向后遍历
	if index <= l.length/2 {
		curNode = l.head
		for i := 0; i <= index; i++ {
			curNode = curNode.next
		}
	} else {
		curNode = l.tail
		for i := l.length - 1; i >= index; i-- {
			curNode = curNode.prev
		}
	}
	return curNode
}

func (l *LinkedList[T]) Get(index int) (T, error) {
	// 检查索引是否在有效范围内
	if index < 0 || index >= l.length {
		var zeroValue T
		return zeroValue, common.NewErrIndexOutOfRange(l.length, index)
	}
	// 获取索引位置的节点的值
	n := l.findNode(index)
	return n.val, nil
}

func (l *LinkedList[T]) Append(values ...T) error {
	for _, value := range values {
		// 创建新的节点，设置其值，并插入到链表尾部（在tail的前一个节点处插入新节点）
		node := &node[T]{prev: l.tail.prev, next: l.tail, val: value}
		// 调整新节点的前后节点的引用，使其正确插入到链表中
		node.prev.next, node.next.prev = node, node
		l.length++
	}
	return nil
}

func (l *LinkedList[T]) Add(index int, value T) error {
	if index < 0 || index > l.length {
		return common.NewErrIndexOutOfRange(l.Len(), index)
	}
	if index == l.length {
		return l.Append(value)
	}
	indexNode := l.findNode(index)
	// 创建新的节点，设置其值，并插入到链表尾部（在tail的前一个节点处插入新节点）
	node := &node[T]{prev: indexNode.prev, next: indexNode, val: value}
	// 调整新节点的前后节点的引用，使其正确插入到链表中
	node.prev.next, node.next.prev = node, node
	l.length++
	return nil
}

func (l *LinkedList[T]) Set(index int, value T) error {
	// 检查索引是否在有效范围内
	if index < 0 || index >= l.length {
		return common.NewErrIndexOutOfRange(l.length, index)
	}
	// 获取索引位置的节点，并设置其值
	node := l.findNode(index)
	node.val = value
	return nil
}

func (l *LinkedList[T]) Delete(index int) (T, error) {
	// 检查索引是否在有效范围内
	if index < 0 || index >= l.length {
		var zeroValue T
		return zeroValue, common.NewErrIndexOutOfRange(l.length, index)
	}
	node := l.findNode(index)
	// 删除节点，并调整链表结构
	node.prev.next, node.next.prev = node.next, node.prev
	node.prev, node.next = nil, nil
	l.length--
	return node.val, nil

}

func (l *LinkedList[T]) Len() int {
	// 获取链表长度
	return l.length
}

func (l *LinkedList[T]) Cap() int {
	// 获取链表的容量（等同于长度）
	return l.length
}

func (l *LinkedList[T]) Range(fn func(index int, value T) error) error {
	// 遍历链表，执行指定的操作函数
	for curNode, i := l.head.next, 0; i < l.length; i++ {
		err := fn(i, curNode.val)
		if err != nil {
			return err
		}
		curNode = curNode.next
	}
	return nil
}

// 将链表转换为切片
func (l *LinkedList[T]) AsSlice() []T {
	result := make([]T, l.length)
	for curNode, i := l.head.next, 0; i < l.length; i++ {
		result[i] = curNode.val
		curNode = curNode.next
	}
	return result
}
