package chanlock

import "time"

//lock empty类型
type empty interface{}

//lock基础结构体
type ChanLock struct {
	ch chan empty
}

//实例化一个lock结构体
func NewLock() *ChanLock {
	return &ChanLock{
		ch: make(chan empty, 1), //创建一个带缓冲/延迟的chan
	}
}

//加锁
func (c *ChanLock) Lock() {
	c.ch <- nil
}

//解锁
func (c *ChanLock) Unlock() {
	<-c.ch
}

//尝试加锁
//失败后, 返回false并停止定时器
func (c *ChanLock) TryLock(timeout time.Duration) bool {
	t := time.NewTimer(timeout)
	defer t.Stop()

	select {
	case c.ch <- nil:
		return true
	case <-t.C:
		return false
	}
}
