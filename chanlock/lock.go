package chanlock

import (
	"sync"
	"time"
)

const (
	waited int32 = iota
	filled
)

//lock empty类型
type empty interface{}

//lock基础结构体
type ChanLock struct {
	ch    chan empty
	state int32
	sync.Mutex
}

//实例化一个lock结构体
//延迟绑定给c结构体实例
//使用mutex控制并行
func (c *ChanLock) newChan() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if c.state == waited {
		c.ch = make(chan empty, 1)
		c.state = filled
	}
}

//加锁
func (c *ChanLock) Lock() {
	if c.state == waited {
		c.newChan()
	}

	c.ch <- nil
}

//解锁
func (c *ChanLock) Unlock() {
	if c.state == waited {
		return
	}
	<-c.ch
}

//尝试加锁
//失败后, 返回false并停止定时器
func (c *ChanLock) TryLock(timeout time.Duration) bool {
	if c.state == waited {
		c.newChan()
	}

	t := time.NewTimer(timeout)
	defer t.Stop()

	select {
	case c.ch <- nil:
		return true
	case <-t.C:
		return false
	}
}
