package goroutinepool

import "github.com/zhaocong6/goUtils/chanlock"

type pool struct {
	capacity int64
	jobs     chan Job
	state    int64
	running  int64
	l        chanlock.ChanLock
}

func (p *pool) build() {
	p.l.Lock()
	defer p.l.Unlock()

	if p.running < p.capacity {
		p.running++

		go func() {
			defer func() {
				p.l.Lock()
				p.running--
				p.l.Unlock()
			}()

			for {
				select {
				case Job := <-p.jobs:
					Job.Handle()
				}
			}
		}()
	}
}

func (p *pool) send(j Job) {
	p.jobs <- j
}
