package goroutinepool

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"
)

type Score struct {
	Num int
}

func TestWorker_Put(t *testing.T) {
	w := NewPool(Worker{
		Capacity: 3,
	})

	var (
		j1, j2 Job
	)
	j1 = &Score{Num: 1}
	j2 = &Score{Num: 1}
	w.Put(j1)
	w.Put(j2)

	time.Sleep(time.Second)
	fmt.Println(w.pool.running)
}

func getGoroutineID() uint64 {

	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)

	return n
}

func (s *Score) Handle() error {
	time.Sleep(time.Millisecond)
	fmt.Println(s.Num, getGoroutineID())
	return nil
}
