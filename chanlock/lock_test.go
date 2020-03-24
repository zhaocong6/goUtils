package chanlock

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestChanLock_Lock(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(100)

	l := NewLock()

	for i := 0; i < 100; i++ {
		i := i
		go func() {
			l.Lock()
			defer func() {
				wg.Done()
				l.Unlock()
			}()

			time.Sleep(time.Millisecond)
			fmt.Println(i)
		}()
	}

	wg.Wait()
}

func TestChanLock_TryLock(t *testing.T) {
	l := NewLock()

	for i := 0; i < 100; i++ {
		i := i
		go func() {
			if ok := l.TryLock(time.Nanosecond); ok {
				defer func() {
					l.Unlock()
				}()

				time.Sleep(time.Millisecond)
				fmt.Println(i)
			} else {
				fmt.Println("lock fail")
			}
		}()
	}

	time.Sleep(time.Second)
}
