package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Kafka 服务器的地址
var addrs = []string{"localhost:9094"}

// 测试同步生产者
func TestSyncProducer(t *testing.T) {
	// 创建 Sarama 配置
	cfg := sarama.NewConfig()

	// 配置生产者返回成功信息
	cfg.Producer.Return.Successes = true

	// 配置分区策略为哈希分区
	cfg.Producer.Partitioner = sarama.NewHashPartitioner

	// 创建同步生产者
	producer, err := sarama.NewSyncProducer(addrs, cfg)
	assert.NoError(t, err)

	// 发送消息
	//_, _, err = producer.SendMessage(&sarama.ProducerMessage{
	//	Topic: "test_topic",
	//	Key:   sarama.StringEncoder("oid-123"),
	//	Value: sarama.StringEncoder("Hello, 这是一条消息 A"),
	//	Headers: []sarama.RecordHeader{
	//		{
	//			Key:   []byte("trace_id"),
	//			Value: []byte("123456"),
	//		},
	//	},
	//	Metadata: "这是metadata",
	//})
	//assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: "read_article",
			Value: sarama.StringEncoder(`{"aid": 1, "uid": 123}`),
		})
		assert.NoError(t, err)
	}

}

// 测试异步生产者
func TestAsyncProducer(t *testing.T) {
	// 创建 Sarama 配置
	cfg := sarama.NewConfig()

	// 配置生产者返回错误和成功信息
	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true

	// 创建异步生产者
	producer, err := sarama.NewAsyncProducer(addrs, cfg)
	require.NoError(t, err)

	// 获取消息通道
	msgCh := producer.Input()

	// 启动一个协程用于向消息通道发送消息
	go func() {
		for {
			msg := &sarama.ProducerMessage{
				Topic: "test_topic",
				Key:   sarama.StringEncoder("oid-123"),
				Value: sarama.StringEncoder("Hello, 这是一条消息 A"),
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte("trace_id"),
						Value: []byte("123456"),
					},
				},
				Metadata: "这是metadata",
			}
			//在通道阻塞的时候，不再发了
			select {
			case msgCh <- msg:
				//default:

				// 向消息通道发送消息
			}
		}
	}()

	// 获取错误和成功的通道
	errCh := producer.Errors()
	succCh := producer.Successes()

	// 循环监听错误和成功的通道
	for {
		// 如果两个情况都没发生，就会阻塞
		select {
		case err := <-errCh:
			t.Log("发送出了问题", err.Err)
		case <-succCh:
			t.Log("发送成功")
		}
	}
}

// 自定义 JSON 编码器
type JSONEncoder struct {
	Data any
}
