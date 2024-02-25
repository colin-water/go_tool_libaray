package redis_lock

import "time"

// RetryStrategy 定义了重试策略的接口
type RetryStrategy interface {
	// Next 返回下一次重试的间隔和是否继续重试的标志
	Next() (time.Duration, bool)
}

// FixedIntervalRetryStrategy 实现了 RetryStrategy 接口，表示固定间隔的重试策略
type FixedIntervalRetryStrategy struct {
	Interval time.Duration // 重试的间隔
	MaxCnt   int           // 最大重试次数
	cnt      int           // 当前重试次数
}

// Next 实现了 RetryStrategy 接口的 Next 方法
func (f *FixedIntervalRetryStrategy) Next() (time.Duration, bool) {
	// 如果当前重试次数已经达到最大重试次数，返回 0 和 false 表示不再继续重试
	if f.cnt >= f.MaxCnt {
		return 0, false
	}
	// 返回重试的间隔和 true 表示继续重试
	return f.Interval, true
}
