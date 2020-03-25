package goroutinepool

import (
	"errors"
	"testing"
	"time"
)

type Score struct {
	Num int
}

type ScorePanic struct {
	Num int
}

func TestWorker_Put(t *testing.T) {
	w := NewPool(Options{
		Capacity:  3,
		JobBuffer: 100,
	})

	var (
		j Job
	)
	j = &Score{Num: 1}
	j = &Score{Num: 1}
	w.Put(j)
	w.Put(j)
	w.Put(j)
	w.Put(j)

	time.Sleep(time.Millisecond * 100)
}

func TestWorker_Cancel(t *testing.T) {
	w := NewPool(Options{
		Capacity:  3,
		JobBuffer: 100,
	})

	var (
		j Job
	)
	j = &Score{Num: 1}
	j = &Score{Num: 1}
	w.Put(j)
	w.Put(j)
	w.Put(j)

	time.Sleep(time.Millisecond * 100)
	w.Cancel()
	time.Sleep(time.Millisecond)
}

func TestWorker_Panic(t *testing.T) {
	w := NewPool(Options{
		Capacity:  3,
		JobBuffer: 100,
	})

	var (
		j Job
	)
	j = &ScorePanic{Num: 1}
	j = &ScorePanic{Num: 1}
	w.Put(j)
	w.Put(j)
	w.Put(j)

	time.Sleep(time.Millisecond * 100)
}

func (s *Score) Handle() error {
	time.Sleep(time.Millisecond)
	return nil
}

func (s *ScorePanic) Handle() error {
	time.Sleep(time.Millisecond)
	return errors.New("错误")
}
