package queue

import (
	"github.com/colin-water/go_tool_libaray/base/common"
	"sync/atomic"
	"unsafe"
)

// 节点结构
type node[T any] struct {
	val T
	// 下一个节点指针
	next unsafe.Pointer
}

// ConcurrentLinkedQueue 无界并发安全队列
// 使用 atomic.LoadPointer
type ConcurrentLinkedQueue[T any] struct {
	// 头节点指针
	head unsafe.Pointer
	// 尾节点指针
	tail unsafe.Pointer
}

// NewConcurrentLinkedQueue 创建一个新的 ConcurrentLinkedQueue 实例
func NewConcurrentLinkedQueue[T any]() *ConcurrentLinkedQueue[T] {
	head := &node[T]{}
	ptr := unsafe.Pointer(head) // 将这个head的指针转换为通用指针类型
	return &ConcurrentLinkedQueue[T]{
		head: ptr,
		tail: ptr,
	}
}

// Enqueue 将元素入队
func (c *ConcurrentLinkedQueue[T]) Enqueue(t T) error {
	// 创建新节点
	newNode := &node[T]{val: t}
	// 将新节点的指针转换为通用指针类型
	newPtr := unsafe.Pointer(newNode)

	for {
		// 获取尾节点指针（原子性的）
		tailPtr := atomic.LoadPointer(&c.tail)
		// 类型断言，将尾节点指针转换为节点类型
		// tail 是指向 tail_node 的指针
		tail := (*node[T])(tailPtr)
		// 获取尾节点的下一个节点指针（tail.next）
		// tail.next 是下一节点的指针
		// &tail.next 是下一节点的指针 的地址
		tailNext := atomic.LoadPointer(&tail.next)

		// 如果尾节点的下一个节点不为空，表示已经被其他线程修改
		if tailNext != nil {
			// 已经被其他线程修改，不需要修复，因为预期中修改的那个线程会更新 c.tail 指针
			continue
		}

		// 使用 CAS 操作设置尾节点的下一个节点为新节点
		// 如果*(&tail.next) 等于 tailNext ，就将*(&tail.next)替换为newPtr

		//tail 是一个指针，指向一个 tail_node，它代表了 LinkedQueue 的尾节点
		// &tail.next 是这个节点的下一节点指针 的地址
		// 如果是false，就是这个tail下一节点的指针不是nil来
		if atomic.CompareAndSwapPointer(&tail.next, tailNext, newPtr) {
			// 如果失败也不用担心，说明有其他线程抢先一步了

			// 最后再尝试更新尾节点指针，确保它指向最新的节点
			// 就是这个新节点成为 尾节点
			// 就是c.tail 变成 newPtr
			atomic.CompareAndSwapPointer(&c.tail, tailPtr, newPtr)

			// 入队成功，退出循环
			return nil
		}
	}
}

// Dequeue 从队列中取出元素
// head默认是指向一个空的node
// 出队返回的是下一个node的val，并且head指向这个node
// 所以说这个node实际上还在这个链表中
func (c *ConcurrentLinkedQueue[T]) Dequeue() (T, error) {
	for {
		// 获取头节点指针
		headPtr := atomic.LoadPointer(&c.head)
		head := (*node[T])(headPtr)

		// 获取尾节点指针
		tailPtr := atomic.LoadPointer(&c.tail)
		tail := (*node[T])(tailPtr)

		// 如果头节点和尾节点相同，表示队列为空
		if head == tail {
			// 在当下这一刻，我们就认为没有元素，即便这时候正好有其他线程入队
			// 但是并不妨碍我们在它彻底入队完成之前认为其实还是没有元素
			var t T
			return t, common.NewErrWithMessage("empty queue")
		}

		// 获取头节点的下一个节点指针
		headNextPtr := atomic.LoadPointer(&head.next)

		// 使用 CAS 操作设置头节点为其下一个节点
		if atomic.CompareAndSwapPointer(&c.head, headPtr, headNextPtr) {
			// 如果 CAS 操作成功，表示头节点出队成功

			// 将头节点的下一个节点的指针转换为节点类型
			headNext := (*node[T])(headNextPtr)

			// 返回出队的元素值
			return headNext.val, nil
		}
		// 如果 CAS 操作失败，说明有其他线程在抢先一步，需要重新尝试出队操作
	}
}
