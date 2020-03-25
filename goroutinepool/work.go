package goroutinepool

import "context"

type Job interface {
	Handle() error
}

type Worker struct {
	pool   *pool
	ctx    context.Context
	cancel context.CancelFunc
}

//配置参数
type Options struct {
	Capacity  int
	JobBuffer int
}

//创建连接池
//实例化连接池结构体
func NewPool(o Options) *Worker {
	ctx, cancel := context.WithCancel(context.Background())

	return &Worker{
		pool: &pool{
			max:     o.Capacity,
			jobs:    make(chan Job, o.JobBuffer),
			running: 0,
			ctx:     ctx,
		},
		ctx:    ctx,
		cancel: cancel,
	}
}

//put任务
func (w *Worker) Put(j Job) {
	w.pool.build()
	w.pool.send(j)
}

//关闭所有worker
func (w *Worker) Cancel() {
	w.cancel()
}
