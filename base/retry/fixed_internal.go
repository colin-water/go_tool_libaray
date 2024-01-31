package retry

import (
	"github.com/yishengzhishui/library/base/common"
	"sync/atomic"
	"time"
)

// FixedIntervalRetryStrategy 等间隔重试策略
type FixedIntervalRetryStrategy struct {
	maxRetries int32         // 最大重试次数，如果是 0 或负数，表示无限重试
	interval   time.Duration // 重试间隔时间
	retries    int32         // 当前重试次数
}

func NewFixedIntervalRetryStrategy(interval time.Duration, maxRetries int32) (*FixedIntervalRetryStrategy, error) {
	if interval <= 0 {
		return nil, common.NewErrInvalidIntervalValue(interval)
	}
	return &FixedIntervalRetryStrategy{
		maxRetries: maxRetries,
		interval:   interval,
		retries:    0,
	}, nil
}

func (s FixedIntervalRetryStrategy) Next() (time.Duration, bool) {
	retries := atomic.AddInt32(&s.retries, 1)

	// maxRetries >0 且 当前重试次数达到最大次数
	if s.maxRetries > 0 && retries > s.maxRetries {
		return 0, false
	}
	return s.interval, true
}
