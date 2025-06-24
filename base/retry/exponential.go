package retry

import (
	"github.com/colin-water/go_tool_libaray/base/common"
	"math"
	"sync/atomic"
	"time"
)

// ExponentialBackoffRetryStrategy 指数退避重试
type ExponentialBackoffRetryStrategy struct {
	initialInterval    time.Duration // 初始重试间隔
	maxInterval        time.Duration // 最大重试间隔
	maxRetries         int32         // 最大重试次数
	retries            int32         // 当前重试次数
	maxIntervalReached atomic.Value  // 是否已经达到最大重试间隔
}

// NewExponentialBackoffRetryStrategy 创建一个指数退避的重试策略实例
func NewExponentialBackoffRetryStrategy(initialInterval, maxInterval time.Duration, maxRetries int32) (*ExponentialBackoffRetryStrategy, error) {
	// 检查参数有效性
	if initialInterval <= 0 {
		return nil, common.NewErrInvalidIntervalValue(initialInterval)
	} else if initialInterval > maxInterval {
		return nil, common.NewErrInvalidMaxIntervalValue(maxInterval, initialInterval)
	}
	return &ExponentialBackoffRetryStrategy{
		initialInterval:    initialInterval,
		maxInterval:        maxInterval,
		maxRetries:         maxRetries,
		retries:            0,
		maxIntervalReached: atomic.Value{},
	}, nil
}

// Next 计算下一次重试的间隔时间，并返回是否可以继续重试
func (s *ExponentialBackoffRetryStrategy) Next() (time.Duration, bool) {
	retries := atomic.AddInt32(&s.retries, 1)
	// 检查是否达到最大重试次数
	if s.maxRetries > 0 && retries > s.maxRetries {
		return 0, false
	}

	// 检查是否已经达到最大重试间隔
	if reached, ok := s.maxIntervalReached.Load().(bool); ok && reached {
		return s.maxInterval, true
	}

	// 计算当前重试间隔，使用指数算法
	interval := s.initialInterval * time.Duration(math.Pow(2, float64(retries-1)))

	// 防止溢出或超过最大重试间隔
	if interval < 0 || interval > s.maxInterval {
		s.maxIntervalReached.Store(true)
		return s.maxInterval, true
	}

	return interval, true
}
