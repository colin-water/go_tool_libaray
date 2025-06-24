package pool

import (
	"github.com/colin-water/go_tool_libaray/base/common"
	"sync"
)

//消息队列，一个消息被多个goroutine消费
//消费被发给多个channel，每个channel对应一个订阅者

type Msg struct {
	Content string
}

type Broker struct {
	mutex    sync.RWMutex
	channels []chan Msg
}

// Send 方法用于向所有订阅者发送消息
func (b *Broker) Send(m Msg) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for _, ch := range b.channels {
		select {
		case ch <- m:
		default:
			return common.NewErrWithMessage("队列已满")
		}
	}
	return nil
}

// Subscribe 方法用于订阅消息，返回一个消息通道
func (b *Broker) Subscribe(capacity int) (<-chan Msg, error) {
	newChan := make(chan Msg, capacity)
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// 将新的订阅者通道添加到 Broker 的 channels 切片中
	b.channels = append(b.channels, newChan)
	return newChan, nil
}

// Close 方法用于关闭所有订阅者通道
func (b *Broker) Close() error {
	b.mutex.Lock()
	channels := b.channels
	b.channels = nil
	b.mutex.Unlock()

	// 避免重复关闭通道
	for _, ch := range channels {
		close(ch)
	}
	return nil
}

//// TestBroker_Send 函数是对 Broker 的 Send 和 Subscribe 方法进行测试的示例
//func TestBroker_Send(t *testing.T) {
//	b := &Broker{}
//
//	// 模拟发送者
//	go func() {
//		for {
//			err := b.Send(Msg{Content: time.Now().String()})
//			if err != nil {
//				t.Log(err)
//				return
//			}
//			time.Sleep(100 * time.Millisecond)
//		}
//	}()
//
//	var wg sync.WaitGroup
//	wg.Add(3)
//	for i := 0; i < 3; i++ {
//		name := fmt.Sprintf("消费者 %d", i)
//
//		// 模拟消费者
//		go func() {
//			defer wg.Done()
//			msgs, err := b.Subscribe(100)
//			if err != nil {
//				t.Log(err)
//				return
//			}
//
//			// 从消息通道中读取消息并打印
//			for msg := range msgs {
//				fmt.Println(name, msg.Content)
//			}
//		}()
//	}
//	wg.Wait()
//}
