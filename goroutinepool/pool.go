package goroutinepool

import (
	"context"
	"github.com/zhaocong6/goUtils/chanlock"
	"log"
	"runtime/debug"
)

type pool struct {
	chanlock.ChanLock
	max     int
	min     int
	running int
	jobs    chan Job
	ctx     context.Context
}

func (p *pool) build() {
	p.Lock()
	defer p.Unlock()

	if p.running < p.max {
		p.running++

		go func(runID int) {

			defer func() {
				p.Lock()
				p.running--
				p.Unlock()
			}()

			for {
				select {
				case <-p.ctx.Done():
					log.Printf("worker %d exit.", runID)
					return
				case co := <-p.jobs:
					func() {
						//异常恢复
						defer p.catchRecover(runID)
						if err := co.Handle(); err != nil {
							panic(err)
						}
					}()
				}
			}

		}(p.running)
	}
}

func (p *pool) send(j Job) {
	p.jobs <- j
}

//异常恢复
func (p *pool) catchRecover(runID int) {
	if err := recover(); err != nil {
		p.Lock()
		log.Printf("exec worker %d error: %s", runID, err)
		debug.PrintStack()
		p.Unlock()
	}
}
