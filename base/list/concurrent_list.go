package list

import "sync"

// ConcurrentList 用读写锁封装了对 List 的操作，实现线程安全的列表操作
type ConcurrentList[T any] struct {
	List[T]              // 嵌套匿名字段，表示 ConcurrentList 继承 List 的所有方法
	lock    sync.RWMutex // 读写锁，用于保护并发访问
}

var (
	_ List[any] = &ConcurrentList[any]{}
)

func (c *ConcurrentList[T]) Get(index int) (T, error) {
	c.lock.RLock()         // 加读锁
	defer c.lock.RUnlock() // 释放读锁
	return c.List.Get(index)
}

func (c *ConcurrentList[T]) Append(ts ...T) error {
	c.lock.Lock()         // 加写锁
	defer c.lock.Unlock() // 释放写锁
	return c.List.Append(ts...)
}

func (c *ConcurrentList[T]) Add(index int, t T) error {
	c.lock.Lock()         // 加写锁
	defer c.lock.Unlock() // 释放写锁
	return c.List.Add(index, t)
}

func (c *ConcurrentList[T]) Set(index int, t T) error {
	c.lock.Lock()         // 加写锁
	defer c.lock.Unlock() // 释放写锁
	return c.List.Set(index, t)
}

func (c *ConcurrentList[T]) Delete(index int) (T, error) {
	c.lock.Lock()         // 加写锁
	defer c.lock.Unlock() // 释放写锁
	return c.List.Delete(index)
}

func (c *ConcurrentList[T]) Len() int {
	c.lock.RLock()         // 加读锁
	defer c.lock.RUnlock() // 释放读锁
	return c.List.Len()
}

func (c *ConcurrentList[T]) Cap() int {
	c.lock.RLock()         // 加读锁
	defer c.lock.RUnlock() // 释放读锁
	return c.List.Cap()
}

func (c *ConcurrentList[T]) Range(fn func(index int, t T) error) error {
	c.lock.RLock()         // 加读锁
	defer c.lock.RUnlock() // 释放读锁
	return c.List.Range(fn)
}

func (c *ConcurrentList[T]) AsSlice() []T {
	c.lock.RLock()         // 加读锁
	defer c.lock.RUnlock() // 释放读锁
	return c.List.AsSlice()
}
