package pool

import (
	"github.com/colin-water/go_tool_libaray/base/common"
	"time"
)

type ReusableObj struct {
}

// ObjPool 对象池，一个缓存ReusableObj 的通道
type ObjPool struct {
	bufChan chan *ReusableObj
}

// NewObjPool 函数，创建并初始化对象池，接受一个参数 numOfObj 表示对象池中可缓存的可重用对象数量
func NewObjPool(num int) *ObjPool {
	pool := ObjPool{}
	pool.bufChan = make(chan *ReusableObj, num)
	//初始化对象
	for i := 0; i < num; i++ {
		pool.bufChan <- &ReusableObj{}
	}
	return &pool
}

// GetObj 方法，从对象池中获取可重用对象，接受一个超时时间参数 timeout
func (p *ObjPool) GetObj(timeout time.Duration) (*ReusableObj, error) {
	select {
	// 拿到对象
	case res := <-p.bufChan:
		return res, nil
	case <-time.After(timeout):
		return nil, common.NewErrWithMessage("time out")
	}
}

// ReleaseObj 方法，将可重用对象释放回对象池
func (p *ObjPool) ReleaseObj(obj *ReusableObj) error {
	select {
	case p.bufChan <- obj:
		return nil
		// 避免阻塞（正常不用）
	default:
		return common.NewErrWithMessage("放回失败")
	}
}
