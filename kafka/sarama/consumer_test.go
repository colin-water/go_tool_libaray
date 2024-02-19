package sarama

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

// TestConsumer 函数：测试 Kafka 消费者，创建了一个消费者组，设置了一个超时上下文，然后使用 Consume 方法进行消息的消费。
//
//testConsumerGroupHandler 结构体：实现了 sarama.ConsumerGroupHandler 接口，处理消费者组的各种生命周期事件和消息消费逻辑。
//
//Setup 方法：在消费者组启动时调用，用于设置分区偏移量。
//
//Cleanup 方法：在消费者组停止时调用，用于清理资源。
//
//ConsumeClaim 方法：实际处理从 Kafka 分区中消费的消息的方法。使用了 errgroup 包来支持对消息的批量处理，具备超时控制，同时处理可能的错误。
//
//ConsumeClaimV1 方法：一个简化版本，只处理一个消息。
//
//MyBizMsg 结构体：示例业务消息的结构体。

// 测试 Kafka 消费者
func TestConsumer(t *testing.T) {
	// 创建 Sarama 配置
	cfg := sarama.NewConfig()

	// 创建 Kafka 消费者组
	consumer, err := sarama.NewConsumerGroup(addrs, "test_group", cfg)
	require.NoError(t, err)

	// 创建上下文，设置超时
	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Minute*10, func() {
		cancel()
	})
	//启动 Kafka 消费者组并开始消费指定主题的消息
	// []string{"test_topic"}：是一个包含要消费的主题名称的字符串切片。在这里，消费者组将会从名为 "test_topic" 的主题中消费消息。
	//testConsumerGroupHandler{}：是一个实现了 sarama.ConsumerGroupHandler 接口的结构体，用于定义消费者组的处理逻辑，包括消息的处理方法等
	err = consumer.Consume(ctx, []string{"test_topic"}, testConsumerGroupHandler{})
	// 消费结束后会到这里
	t.Log(err, time.Since(start).String())
}

// Kafka 消费者组处理器
type testConsumerGroupHandler struct {
}

// Setup 方法在消费者组启动时调用
func (t testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	// 获取分配给消费者组的分区偏移量
	partitions := session.Claims()["test_topic"]

	// 不是最优方案，建议走离线渠道重置偏移量
	// 针对每个分区将偏移量重置为最早的位置
	for _, part := range partitions {
		session.ResetOffset("test_topic", part, sarama.OffsetOldest, "")
	}

	return nil
}

// Cleanup 方法在消费者组停止时调用
func (t testConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("CleanUp")
	return nil
}

// session 是和kafka的会话，（从简历连接到连接断掉）
// ConsumeClaim 方法用于实际处理从 Kafka 分区中消费的消息
// 在循环中，它从消息通道中获取一批消息，
//然后启动 goroutine 并行处理这批消息。
//处理完毕后，标记最新的消息为已处理。
//整个处理过程通过上下文的超时控制，以及 errgroup 的协程组方式，
//实现了批量且有超时控制的消息处理。
func (t testConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 获取分配给当前消费者组的消息通道
	msgs := claim.Messages()

	// 设置批量处理的消息数量
	const batchSize = 10

	// 使用无限循环处理消息
	for {
		// 创建带有超时的上下文
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)

		// 使用 errgroup.Group 支持批量处理
		var eg errgroup.Group
		var last *sarama.ConsumerMessage
		done := false

		// 循环处理消息，每次处理 batchSize 条消息，或者超时
		for i := 0; i < batchSize && !done; i++ {
			select {
			// 如果超时，结束当前处理，在这里就是1秒内凑不足10个，就是超时，for循环结束
			case <-ctx.Done():
				done = true
			// 从消息通道中获取消息
			case msg, ok := <-msgs:
				if !ok {
					// 通道关闭，代表消费者被关闭，结束当前处理
					cancel()
					return nil
				}
				// 记录最新的消息
				last = msg
				// 启动 goroutine 处理消息
				eg.Go(func() error {
					// 模拟消息处理过程，这里是睡眠一秒
					time.Sleep(time.Second)
					// 在这里可以进行消息的实际处理逻辑，比如解析、存储等
					log.Println(string(msg.Value))
					return nil
				})
			}
		}
		// 取消上下文，结束超时处理
		cancel()

		// 等待所有消息处理 goroutines 完成
		err := eg.Wait()
		if err != nil {
			// 处理错误，可以记录日志等
			continue
		}

		// 标记已处理的最新消息（最后一个）
		if last != nil {
			session.MarkMessage(last, "")
		}
	}
}

// ConsumeClaimV1 是一个简化版本，只处理一个消息
func (t testConsumerGroupHandler) ConsumeClaimV1(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		// 处理消息
		log.Println(string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}
