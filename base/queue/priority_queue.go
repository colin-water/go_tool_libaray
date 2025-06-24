package queue

import (
	"github.com/colin-water/go_tool_libaray/base/common"
	"github.com/colin-water/go_tool_libaray/base/slice"
)

// PriorityQueue 是一个基于小顶堆的优先队列
// 当 capacity <= 0 时，为无界队列，切片容量会动态扩缩容
// 当 capacity > 0 时，为有界队列，初始化后就固定容量，不会扩缩容
type PriorityQueue[T any] struct {
	// 用于比较前一个元素是否小于后一个元素
	compare common.Comparator[T]
	// 队列容量
	capacity int
	// 队列中的元素，为便于计算父子节点的 index，0 位置留空，根节点从 1 开始
	// i, left i*2, right i*2+1
	data []T
}

// NewPriorityQueue 创建优先队列
// 当 capacity <= 0 时，为无界队列，否则有界队列
func NewPriorityQueue[T any](capacity int, compare common.Comparator[T]) *PriorityQueue[T] {
	// 根据传入的容量设置切片的容量
	sliceSize := capacity + 1
	if capacity < 1 {
		capacity = 0
		sliceSize = 2 << 5
	}
	// 创建并初始化优先队列实例
	return &PriorityQueue[T]{
		capacity: capacity,
		data:     make([]T, 1, sliceSize),
		compare:  compare,
	}
}

// Len 返回队列长度
func (p *PriorityQueue[T]) Len() int {
	return len(p.data) - 1
}

// Cap 无界队列返回 0，有界队列返回创建队列时设置的值
func (p *PriorityQueue[T]) Cap() int {
	return p.capacity
}

// IsBoundless 判断是否为无界队列
func (p *PriorityQueue[T]) IsBoundless() bool {
	return p.capacity <= 0
}

// isFull 判断队列是否已满
func (p *PriorityQueue[T]) isFull() bool {
	return p.capacity > 0 && len(p.data)-1 == p.capacity
}

// isEmpty 判断队列是否为空
func (p *PriorityQueue[T]) isEmpty() bool {
	return len(p.data) < 2
}

// Peek 获取队列首个元素
func (p *PriorityQueue[T]) Peek() (t T, e error) {
	if p.isEmpty() {
		return t, common.NewErrWithMessage("empty queue")
	}
	return p.data[1], nil
}

// Enqueue 将元素入队
func (p *PriorityQueue[T]) Enqueue(t T) error {
	// 判断队列是否已满
	if p.isFull() {
		return common.NewErrWithMessage("full queue")
	}

	// 将元素添加到队列末尾
	p.data = append(p.data, t)

	// 调整队列，使其保持小顶堆的性质
	// 计算新加入元素和其父节点的索引
	node, parent := len(p.data)-1, (len(p.data)-1)/2

	// 循环直到满足小顶堆的性质
	for parent > 0 && p.compare(p.data[node], p.data[parent]) < 0 {
		// 如果新加入元素比其父节点小，交换它们的位置
		p.data[parent], p.data[node] = p.data[node], p.data[parent]

		// 更新当前节点和父节点的索引，继续上浮操作
		node = parent
		parent = node / 2
	}

	return nil
}

// Dequeue 将元素出队
func (p *PriorityQueue[T]) Dequeue() (t T, e error) {
	// 判断队列是否为空
	if p.isEmpty() {
		return t, common.NewErrWithMessage("empty queue")
	}

	// 弹出队列的首个元素
	pop := p.data[1]

	// 将队列末尾的元素移到首个位置
	p.data[1] = p.data[len(p.data)-1]

	// 缩短队列切片，去除末尾元素
	p.data = p.data[:len(p.data)-1]

	// 根据需要缩短队列切片容量（仅对无界队列有效）
	p.shrinkIfNecessary()

	// 调整队列，使其保持小顶堆的性质
	p.heapify(p.data, len(p.data)-1, 1)

	// 返回出队的元素
	return pop, nil
}

// shrinkIfNecessary 如果是无界队列，根据需要收缩队列切片容量
func (p *PriorityQueue[T]) shrinkIfNecessary() {
	if p.IsBoundless() {
		p.data = slice.Shrink[T](p.data)
	}
}

// heapify 用于维护小顶堆的性质
func (p *PriorityQueue[T]) heapify(data []T, n, i int) {
	// 初始化最小位置
	minPos := i
	for {
		// 计算左子节点的位置, 判定 data[left] 是否比data[minPos]小
		if left := i * 2; left <= n && p.compare(data[left], data[minPos]) < 0 {
			minPos = left
		}
		// 计算右子节点的位置
		if right := i*2 + 1; right <= n && p.compare(data[right], data[minPos]) < 0 {
			minPos = right
		}
		// 如果最小位置等于当前位置，表示堆的性质已满足，跳出循环
		if minPos == i {
			break
		}
		// 交换当前位置和最小位置的元素
		data[i], data[minPos] = data[minPos], data[i]
		// 更新当前位置为最小位置，继续循环调整
		i = minPos
	}
}
