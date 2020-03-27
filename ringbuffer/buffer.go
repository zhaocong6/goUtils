package ringbuffer

import (
	"sync"
)

type Value interface{}

type readBuffer <-chan Value

type writeBuffer chan<- Value

type Buffer struct {
	capacity   int
	ReadBuffer readBuffer
	writeBuffer
	lock sync.Mutex
}

type Options struct {
	Capacity int
}

//实例化一个buffer
func New(o Options) *Buffer {

	//最小容量为1
	if o.Capacity < 1 {
		o.Capacity = 1
	}

	ch := make(chan Value, o.Capacity)

	return &Buffer{
		capacity:    o.Capacity,
		ReadBuffer:  ch,
		writeBuffer: ch,
		lock:        sync.Mutex{},
	}
}

//读取buffer
func (b Buffer) Read() interface{} {
	return <-b.ReadBuffer
}

//写入buffer
func (b Buffer) Write(v Value) {
	b.lock.Lock()
	defer b.lock.Unlock()

	//验证写入buffer是否已经满了
	//写入buffer满了以后, 取出最先进入的buffer
	if len(b.writeBuffer) == b.capacity {
		select {
		case <-b.ReadBuffer:
		default:
		}
	}

	b.writeBuffer <- v
}
