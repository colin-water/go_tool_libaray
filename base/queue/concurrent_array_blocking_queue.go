package queue

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

// ConcurrentArrayBlockingQueue 是一个有界并发阻塞队列的实现。
// 使用互斥锁（c.mutex）来保护队列的并发访问。
// 使用两个信号量（c.enqueueCap 和 c.dequeueCap）来控制入队和出队的容量、超时和阻塞问题。
// head 和 tail 在切片中可能发生循环： 在循环队列中，当 tail 超过数组的末尾时，它会循环回到数组的开头。
// 因此，如果 tail 的值小于 head，说明队列在数组中发生了循环。
type ConcurrentArrayBlockingQueue[T any] struct {
	data  []T
	mutex *sync.RWMutex // 使用指针类型，更小的开销，但是需要显式的初始化

	// 队头元素下标
	head int
	// 队尾元素下标
	tail int
	// 包含多少个元素
	count int

	// 入队信号量
	enqueueCap *semaphore.Weighted
	// 出队信号量
	dequeueCap *semaphore.Weighted

	// zero 不能作为返回值返回，防止用户篡改
	zero T
}

// NewConcurrentArrayBlockingQueue 创建一个有界阻塞队列。
// 容量在最开始就初始化好。
// capacity 必须为正数。
func NewConcurrentArrayBlockingQueue[T any](capacity int) *ConcurrentArrayBlockingQueue[T] {
	mutex := &sync.RWMutex{}

	// 创建入队信号量
	semaForEnqueue := semaphore.NewWeighted(int64(capacity))
	// 创建出队信号量，一开始就设置为满的，即无法进行出队操作
	semaForDequeue := semaphore.NewWeighted(int64(capacity))
	// 它表明代码中应该没有处理 Context 的具体逻辑。它通常在还没有明确的 Context 时使用。
	_ = semaForDequeue.Acquire(context.TODO(), int64(capacity)) // 确保一开始不能进行出队操作

	res := &ConcurrentArrayBlockingQueue[T]{
		data:       make([]T, capacity),
		mutex:      mutex,
		enqueueCap: semaForEnqueue,
		dequeueCap: semaForDequeue,
	}
	return res
}

// Enqueue 入队操作。
// 使用入队信号量控制容量、超时和阻塞。
func (c *ConcurrentArrayBlockingQueue[T]) Enqueue(ctx context.Context, t T) error {
	// 尝试获取入队信号量，能拿到则说明队列还有空位可以入队，否则阻塞等待
	err := c.enqueueCap.Acquire(ctx, 1)
	if err != nil {
		return err
	}

	// 获取互斥锁，保护对队列的并发访问
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 检查上下文是否已取消，避免在获取锁期间发生取消操作
	if ctx.Err() != nil {
		// 超时或取消，释放入队信号量，防止容量泄露
		c.enqueueCap.Release(1)
		return ctx.Err()
	}

	// 执行入队操作
	c.data[c.tail] = t
	c.tail++
	c.count++

	// 队尾已经到达数组末尾，重置下标
	if c.tail == cap(c.data) {
		c.tail = 0
	}

	// 发送出队信号，通知等待的出队操作可以进行了
	c.dequeueCap.Release(1)

	return nil
}

// Dequeue 出队操作。
// 使用出队信号量控制容量、超时和阻塞。
func (c *ConcurrentArrayBlockingQueue[T]) Dequeue(ctx context.Context) (T, error) {
	// 尝试获取出队信号量，能拿到则说明队列有元素可以出队，否则阻塞等待
	err := c.dequeueCap.Acquire(ctx, 1)
	var res T

	if err != nil {
		return res, err
	}

	// 获取互斥锁，保护对队列的并发访问
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 检查上下文是否已取消，避免在获取锁期间发生取消操作
	if ctx.Err() != nil {
		// 超时或取消，释放出队信号量，有元素消费不到
		c.dequeueCap.Release(1)
		return res, ctx.Err()
	}

	// 执行出队操作
	res = c.data[c.head]
	// 释放元素所占用的内存，帮助 GC
	c.data[c.head] = c.zero
	c.head++
	c.count--

	// 队头已经到达数组末尾，重置下标
	if c.head == cap(c.data) {
		c.head = 0
	}

	// 发送入队信号，通知等待的入队操作可以进行了
	c.enqueueCap.Release(1)

	return res, nil
}

// Len 返回队列中元素的数量。
func (c *ConcurrentArrayBlockingQueue[T]) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.count
}

// AsSlice 返回队列中元素的切片表示。
// AsSlice 方法将队列中的元素转换成切片并返回
func (c *ConcurrentArrayBlockingQueue[T]) AsSlice() []T {
	// 使用读取锁保护并发读取
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// 创建一个初始长度为0、容量为队列中元素数量的切片
	res := make([]T, 0, c.count)

	// 初始化迭代计数器
	cnt := 0

	// 获取队列底层数组的容量
	capacity := cap(c.data)

	// 遍历队列中的元素
	for cnt < c.count {
		// 计算当前元素的索引，考虑到队列可能在数组的末尾循环
		//确保索引在底层数组容量范围内，并且当 c.head + cnt 的和超过数组末尾时，能够循环回到数组的开头。
		index := (c.head + cnt) % capacity

		// 将当前元素添加到切片中
		res = append(res, c.data[index])

		// 增加计数器，以处理下一个元素
		cnt++
	}

	// 返回包含队列中所有元素的切片
	return res
}
