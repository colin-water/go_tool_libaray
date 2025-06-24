package list

import (
	"errors"
	"github.com/colin-water/go_tool_libaray/base/common"
	"math/rand"
)

// 跳表 skip list
//  FactorP  = 0.25,  MaxLevel = 32 列表可包含 2^64 个元素
const (
	FactorP  = float32(0.25) // level i 上的结点 有FactorP的比例出现在level i + 1上
	MaxLevel = 32
)

// skipListNode 表示跳表中的结点
type skipListNode[T any] struct {
	Val     T                  // 结点的值
	Forward []*skipListNode[T] // Forward 切片存储了结点在不同层级上的下一个结点
	//具体来说，对于第 i 层，Forward[i] 存储了当前节点在第 i 层的右侧节点。
}

// SkipList 是跳表的实现
type SkipList[T any] struct {
	header *skipListNode[T] // 跳表头结点 header 被用作虚拟头节点，它不包含实际的数据。
	// header 节点本身只有一个，但它的 Forward 数组中有指向每一层的指针
	level   int                  // 跳表的最大层级
	compare common.Comparator[T] // 用于比较元素大小的比较器
	size    int                  // 跳表的元素个数
}

// newSkipListNode 创建一个新的跳表结点
func newSkipListNode[T any](Val T, level int) *skipListNode[T] {
	return &skipListNode[T]{
		Val,
		make([]*skipListNode[T], level+1),
	}
}

// NewSkipList 创建一个新的跳表
func NewSkipList[T any](compare common.Comparator[T]) *SkipList[T] {
	return &SkipList[T]{
		header: &skipListNode[T]{
			Forward: make([]*skipListNode[T], MaxLevel+1),
		},
		level:   1,
		compare: compare,
	}
}

// randomLevel 随机生成结点的层级
// level=1的概率是1-FactorP
// level=2的概率是(1-FactorP)*FactorP
// level越高，概率越低
func (sl *SkipList[T]) randomLevel() int {
	// 初始化层级为 1
	level := 1
	// 获取跳表的因子 p
	p := FactorP

	// 使用位运算生成随机数，直到生成的随机数小于 p 的比例
	//`0xFFFF` 是一个十六进制数，对应的十进制是 65535
	//`p * 0xFFFF` 将因子 `p` 乘以 65535，得到的结果是一个小数。
	//`int32()` 将上述结果转换为 `int32` 类型
	for (rand.Int31() & 0xFFFF) < int32(p*0xFFFF) {
		// 层级加一
		level++
	}

	// 如果生成的层级小于最大层级 MaxLevel，则返回生成的层级，否则返回最大层级 MaxLevel
	if level < MaxLevel {
		return level
	}
	return MaxLevel
}

// Peek 获取跳表的第一个元素，如果跳表为空则返回错误
func (sl *SkipList[T]) Peek() (T, error) {
	// 获取头结点的下一个结点（第一层的第一个结点）
	curr := sl.header.Forward[1]

	// 定义一个零值用于错误返回
	var zero T

	// 如果第一个结点为空，表示跳表为空
	if curr == nil {
		return zero, errors.New("跳表为空")
	}

	// 返回第一个结点的值
	return curr.Val, nil
}

// Len 返回跳表的元素个数
func (sl *SkipList[T]) Len() int {
	return sl.size
}

// Get 获取跳表指定索引处的元素(沿着第1层)
func (sl *SkipList[T]) Get(index int) (T, error) {
	// 检查索引是否在有效范围内
	if index < 0 || index >= sl.size {
		var zero T
		return zero, common.NewErrIndexOutOfRange(sl.size, index)
	}

	// 初始化当前结点为头结点
	curr := sl.header

	// 循环直到达到目标索引
	for i := 0; i <= index; i++ {
		// 获取当前结点的下一个结点（沿着第1层）
		curr = curr.Forward[1]
	}

	// 返回目标索引处的元素值
	return curr.Val, nil
}

// traverse 寻找元素插入位置，并返回当前节点以及每层需要更新的节点
func (sl *SkipList[T]) traverse(Val T, level int) (*skipListNode[T], []*skipListNode[T]) {
	// 初始化一个数组，用于记录每一层需要更新的节点
	//update[i] 包含位于level i 的插入/删除位置左侧的指针
	update := make([]*skipListNode[T], MaxLevel+1)

	// 从跳表的头节点开始，逐层向下查找插入位置
	curr := sl.header
	for i := level; i > 0; i-- {
		// 在当前层中，向右遍历节点，找到插入位置的前一个节点
		for curr.Forward[i] != nil && sl.compare(curr.Forward[i].Val, Val) < 0 {
			curr = curr.Forward[i]
		}
		// 记录当前层的插入位置的前一个节点
		update[i] = curr
	}

	// 返回当前节点（最底层的）和每层需要更新的节点数组
	return curr, update
}

// Insert 向跳表中插入一个元素
func (sl *SkipList[T]) Insert(Val T) {
	// 获取插入位置的信息
	_, update := sl.traverse(Val, sl.level)

	// 随机生成新节点的层级
	level := sl.randomLevel()

	// 如果新节点的层级比跳表的当前层级高，需要更新跳表的层级
	if level > sl.level {
		for i := sl.level + 1; i <= level; i++ {
			// 因为update[i] 包含位于level i 的插入/删除位置左侧的指针
			//将 update 数组中对应层级的元素更新为跳表的头节点 sl.header。
			update[i] = sl.header
		}
		sl.level = level
	}
	// 创建新节点
	newNode := newSkipListNode[T](Val, level)

	// 将新节点插入到跳表中
	for i := 1; i <= level; i++ {
		//第i层链表插入新的节点
		newNode.Forward[i] = update[i].Forward[i]
		update[i].Forward[i] = newNode
	}
	// 更新跳表的元素数量
	sl.size += 1
}

// Search 在跳表中查找指定元素是否存在
func (sl *SkipList[T]) Search(target T) bool {
	//curr 目标节点的前一个节点
	curr, _ := sl.traverse(target, sl.level)
	curr = curr.Forward[1]
	return curr != nil && sl.compare(curr.Val, target) == 0
}

// DeleteElement 从跳表中删除指定元素
func (sl *SkipList[T]) DeleteElement(target T) bool {
	// 找到目标节点的前一个节点
	curr, update := sl.traverse(target, sl.level)
	//目标节点
	node := curr.Forward[1]
	if node == nil || sl.compare(node.Val, target) != 0 {
		return true
	}
	// 删除target结点
	//遍历每个层级，检查目标节点在每个层级的前一个节点是否正确。
	//如果正确，将前一个节点的下一个节点更新为目标节点的下一个节点，实现删除操作。
	for i := 1; i <= sl.level && update[i].Forward[i] == node; i++ {
		update[i].Forward[i] = node.Forward[i]
	}
	sl.size -= 1
	// 更新层级，如果层级元素空了，删除层级
	for sl.level > 1 && sl.header.Forward[sl.level] == nil {
		sl.level--
	}

	return true
}

// AsSlice 将跳表转换为切片并返回
func (sl *SkipList[T]) AsSlice() []T {
	// 初始化当前节点为跳表的头节点
	curr := sl.header
	// 使用 make 创建一个切片，初始长度为 0，容量为 sl.size（跳表的元素个数）
	slice := make([]T, 0, sl.size)

	// 遍历跳表，直到当前节点的第一层的后继节点为空（跳表的结束条件）
	for curr.Forward[1] != nil {
		// 将当前节点的第一层后继节点的值追加到切片中
		slice = append(slice, curr.Forward[1].Val)
		// 移动到下一个节点，继续遍历
		curr = curr.Forward[1]
	}

	// 返回存储了跳表所有元素的切片
	return slice
}

// NewSkipListFromSlice 从切片创建一个新的跳表
func NewSkipListFromSlice[T any](slice []T, compare common.Comparator[T]) *SkipList[T] {
	sl := NewSkipList(compare)
	for _, n := range slice {
		sl.Insert(n)
	}
	return sl
}
