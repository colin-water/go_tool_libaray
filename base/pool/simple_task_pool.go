package pool

import "context"

// 任务池的实现

type TaskPool struct {
	ch chan struct{}
}

// NewTaskPool 初始化
//通过limit 限制并发数量
func NewTaskPool(limit int) *TaskPool {
	pool := TaskPool{
		ch: make(chan struct{}, limit),
	}
	for i := 0; i < limit; i++ {
		pool.ch <- struct{}{}
	}
	return &pool
}

// TaskDo 异步执行task
func (t *TaskPool) TaskDo(taskFunc func()) {
	token := <-t.ch
	go func() {
		taskFunc()
		t.ch <- token
	}()

}

// 方案二
type Task func()

// TaskPoolWithClose 使用缓存并可以close
type TaskPoolWithClose struct {
	tasks chan Task
	close chan struct{}
}

func NewTaskPoolWithClose(num int, capacity int) *TaskPoolWithClose {
	pool := &TaskPoolWithClose{
		tasks: make(chan Task, capacity),
		close: make(chan struct{}),
	}

	// 需要close机制，防止goroutine 泄露
	for i := 0; i < num; i++ {
		go func() {
			for {
				select {
				case <-pool.close:
					return
				case task := <-pool.tasks:
					task()
				}
			}
		}()
	}
	return pool
}

// Submit 提交任务
func (p *TaskPoolWithClose) Submit(ctx context.Context, task Task) error {
	select {
	case p.tasks <- task:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// Close 方法会释放资源
func (p *TaskPoolWithClose) Close() error {
	// 重复调用 Close 方法，会 panic
	close(p.close) // 关闭通道
	return nil
}
