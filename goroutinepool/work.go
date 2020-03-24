package goroutinepool

type Job interface {
	Handle() error
}

type Worker struct {
	pool     *pool
	Capacity uint64
}

func NewPool(w Worker) *Worker {
	return &Worker{
		pool: &pool{
			capacity: int64(w.Capacity),
			jobs:     make(chan Job, 100),
			state:    1,
			running:  0,
		},
	}
}

func (w *Worker) Put(j Job) {
	w.pool.build()
	w.pool.send(j)
}
