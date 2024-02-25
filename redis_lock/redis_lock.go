package redis_lock

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"time"
)

var (
	ErrFailedToPreemptLock = errors.New("redis-lock: 抢锁失败")
	ErrLockNotHold         = errors.New("redis-lock: 你没有持有锁")

	//go:embed lua/unlock.lua
	luaUnlock string

	//go:embed lua/refresh.lua
	luaRefresh string

	//go:embed lua/lock.lua
	luaLock string
)

// Client 就是对 redis.Cmdable 的二次封装
// Client 结构体封装了 Redis 客户端和一些并发控制机制
type Client struct {
	client redis.Cmdable // Redis 客户端对象
	// 用于处理对相同资源的重复请求的并发控制
	// 确保对于相同 key 的请求只有一个会执行实际的获取锁操作
	g singleflight.Group
}

func NewClient(client redis.Cmdable) *Client {
	return &Client{
		client: client,
	}
}

type Lock struct {
	client     redis.Cmdable
	key        string
	value      string
	expiration time.Duration
	unlockChan chan struct{}
}

// SingleflightLock 是对 Lock 方法的包装，使用 singleflight 保证对相同 key 的并发请求只会有一个获得锁
func (c *Client) SingleflightLock(ctx context.Context,
	key string,
	expiration time.Duration,
	timeout time.Duration, retry RetryStrategy) (*Lock, error) {
	for {
		flag := false
		// 使用 singleflight 来包装 Lock 方法，
		// 确保对相同 key 的并发请求只有一个会执行实际的获取锁操作

		// 当调用 DoChan 时，它首先检查缓存中是否已经存在标识符为 key 的函数调用结果。
		//如果存在，则直接返回结果；否则，执行提供的函数。

		// resCh 是一个 DoChan 的结果通道，可以通过它获取函数调用的结果或错误。
		//在这里，主要用于获取 c.Lock 方法的结果。
		resCh := c.g.DoChan(key, func() (interface{}, error) {
			// 在函数调用之前，将 flag 设为 true，以便标记当前正在执行获取锁的操作
			flag = true
			// 调用 Lock 方法
			return c.Lock(ctx, key, expiration, timeout, retry)
		})
		select {
		case res := <-resCh:
			if flag {
				// 如果是当前请求的锁，则 Forget 掉 singleflight，然后返回结果
				c.g.Forget(key)
				if res.Err != nil {
					return nil, res.Err
				}
				return res.Val.(*Lock), nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// Lock 尝试获取锁，如果锁未被其他协程持有，则成功获取锁，返回 Lock 实例，
//否则通过重试策略进行重试。
// expiration: 锁的过期时间，即锁被自动释放的时间。
//timeout: 获取锁的超时时间，即尝试获取锁的最长等待时间。
//retry: 重试策略接口，用于确定下一次重试的间隔和是否继续重试。
func (c *Client) Lock(ctx context.Context, key string, expiration time.Duration, timeout time.Duration, retry RetryStrategy) (*Lock, error) {
	var timer *time.Timer
	val := uuid.New().String() // 生成唯一的锁值

	for {
		// 在这里重试
		lctx, cancel := context.WithTimeout(ctx, timeout)

		// 使用 Lua 脚本尝试获取锁
		// 这边是有三种情况：
		// 1.key 不存在
		// 2.你上次加锁成功了但是返回超时了
		// 3.锁被人家拿着
		res, err := c.client.Eval(lctx, luaLock, []string{key}, val, expiration.Seconds()).Result()
		cancel()

		// 处理获取锁的结果和错误
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return nil, err // 如果出现非超时错误，直接返回错误
		}

		// 如果成功获取到锁，返回 Lock 实例
		if res == "OK" {
			return &Lock{
				client:     c.client,
				key:        key,
				value:      val,
				expiration: expiration,
				unlockChan: make(chan struct{}, 1),
			}, nil
		}

		// 根据重试策略获取下一次重试的间隔和是否继续重试的标志
		interval, ok := retry.Next()

		// 如果不允许继续重试，返回超出重试限制的错误
		if !ok {
			return nil, fmt.Errorf("redis-lock: 超出重试限制, %w", ErrFailedToPreemptLock)
		}

		// 设置定时器，等待下一次重试
		if timer == nil {
			timer = time.NewTimer(interval)
		} else {
			timer.Reset(interval)
		}

		// 在定时器超时或者上下文被取消时，退出循环
		select {
		case <-timer.C:
			// 定时器触发，执行下一次重试
		case <-ctx.Done():
			// 上下文被取消，退出循环并返回错误
			return nil, ctx.Err()
		}
	}
}

// AutoRefresh 自动续约
//自动刷新锁的过期时间，interval 表示刷新间隔，timeout 表示刷新的超时时间
func (l *Lock) AutoRefresh(interval time.Duration, timeout time.Duration) error {
	// 创建一个带缓冲通道，用于在超时时通知刷新
	timeoutChan := make(chan struct{}, 1)
	// 创建定时器，每隔 interval 时间触发一次
	ticker := time.NewTicker(interval)

	// 无限循环，实现自动续约
	for {
		select {
		case <-ticker.C:
			// 定时器触发，执行刷新操作
			// 刷新的超时时间怎么设置
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()
			if errors.Is(err, context.DeadlineExceeded) {
				// 如果刷新超时，向 timeoutChan 发送通知
				timeoutChan <- struct{}{}
				continue
			}
			if err != nil {
				// 如果刷新遇到其他错误，返回错误
				return err
			}
		case <-timeoutChan:
			// 从 timeoutChan 接收通知，执行刷新操作
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			// 出现了 error 了怎么办？
			err := l.Refresh(ctx)
			cancel()
			if errors.Is(err, context.DeadlineExceeded) {
				// 如果刷新超时，再次向 timeoutChan 发送通知
				timeoutChan <- struct{}{}
				continue
			}
			if err != nil {
				// 如果刷新遇到其他错误，返回错误
				return err
			}
		case <-l.unlockChan:
			// 收到解锁通知，结束循环
			return nil
		}
	}
}

// Refresh 续约 刷新锁的过期时间
func (l *Lock) Refresh(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaRefresh, []string{l.key}, l.value, l.expiration.Seconds()).Int64()
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotHold
	}
	return nil
}

//--- 基础加锁和释放锁

// TryLock 尝试获取锁，如果锁未被其他协程持有，则成功获取锁，返回 Lock 实例，否则返回 ErrFailedToPreemptLock 错误。
func (c *Client) TryLock(ctx context.Context, key string, expiration time.Duration) (*Lock, error) {
	// 生成一个唯一的锁值
	val := uuid.New().String()

	// 使用 SetNX 命令尝试设置锁，如果成功返回 true，表示锁未被其他协程持有
	// expiration 是过期时间
	ok, err := c.client.SetNX(ctx, key, val, expiration).Result()
	if err != nil {
		return nil, err // 发生错误时返回错误信息
	}

	// 如果 SetNX 返回 false，说明锁已被其他协程持有，返回 ErrFailedToPreemptLock 错误
	if !ok {
		return nil, ErrFailedToPreemptLock
	}

	// 如果成功获取到锁，创建并返回 Lock 实例
	return &Lock{
		client:     c.client,               // 设置锁的 Redis 客户端
		key:        key,                    // 锁的键
		value:      val,                    // 锁的值，用于释放锁时校验
		expiration: expiration,             // 锁的过期时间
		unlockChan: make(chan struct{}, 1), // 创建用于通知释放锁的通道
	}, nil
}

// Unlock 释放锁
func (l *Lock) Unlock(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaUnlock, []string{l.key}, l.value).Int64()
	defer func() {
		select {
		case l.unlockChan <- struct{}{}:
		default:
			// 说明没有人调用 AutoRefresh
		}
	}()
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotHold
	}
	return nil
}
