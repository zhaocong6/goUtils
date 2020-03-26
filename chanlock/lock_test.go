package chanlock

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestChanLock_Lock(t *testing.T) {
	var (
		wg    sync.WaitGroup
		l     ChanLock
		n     int
		total = 1000
	)

	wg.Add(total)

	for i := 0; i < total; i++ {
		go func() {
			l.Lock()
			defer func() {
				wg.Done()
				l.Unlock()
			}()
			n++
		}()
	}

	wg.Wait()

	fmt.Println(total == n, n)
}

func BenchmarkChanLock_Lock(b *testing.B) {
	var (
		wg    sync.WaitGroup
		l     ChanLock
		n     int
		total = 200000
	)

	wg.Add(total)

	for i := 0; i < total; i++ {
		go func() {
			l.Lock()
			defer func() {
				wg.Done()
				l.Unlock()
			}()
			n++
		}()
	}

	wg.Wait()
}

func BenchmarkChanLock_MuLock(b *testing.B) {
	var (
		wg    sync.WaitGroup
		l     sync.Mutex
		n     int
		total = 200000
	)

	wg.Add(total)

	for i := 0; i < total; i++ {
		go func() {
			l.Lock()
			defer func() {
				wg.Done()
				l.Unlock()
			}()
			n++
		}()
	}

	wg.Wait()
}

func TestChanLock_TryLock(t *testing.T) {
	var (
		l     ChanLock
		n     int
		total = 1000
	)

	for i := 0; i < total; i++ {
		go func() {
			//测试时注意timeout时间
			//过短可能导致try失败
			if ok := l.TryLock(time.Millisecond * 100); ok {
				defer func() {
					l.Unlock()
				}()

				n++
			}
		}()
	}

	time.Sleep(time.Millisecond * 100)
	fmt.Println(total == n, n)
}
