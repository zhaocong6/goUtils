package glog

import (
	"log"
	"sync"
	"testing"
)

func TestDebug(t *testing.T) {
	SetDir("/home/zc/goApplication/src/utils")
	InitLog()

	num := 10000
	var wg sync.WaitGroup
	wg.Add(num)

	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()

			Debug("测试debug msg", nil)
		}()
	}

	wg.Wait()
	log.Println("write log success")
}

func TestLogger_Write(t *testing.T) {
	SetDir("/home/zc/goApplication/src/utils")
	InitLog()
	log.SetOutput(&Logger{})

	num := 10000
	var wg sync.WaitGroup
	wg.Add(num)

	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			log.Println("测试 log Println")
		}()
	}

	wg.Wait()

	//log.Println("write log success")
}
