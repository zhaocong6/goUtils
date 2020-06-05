package redislock

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"sync"
	"testing"
	"time"
)

func newRedis() (redis.Conn, error) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func TestLocker_TryLock(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 20; i++ {
		ii := i
		go func() {
			defer wg.Done()

			redisConn, err := newRedis()
			if err != nil {
				log.Panicln("redis连接失败")
			}
			defer redisConn.Close()

			lock := New(redisConn, "order:1", 100)

			if lock.TryLock() != nil {
				log.Println("fail:", ii)
			} else {
				defer lock.Unlock()
				//模拟代码运行
				time.Sleep(time.Second)
				log.Println("success:", ii)
			}
		}()
	}

	wg.Wait()
}
