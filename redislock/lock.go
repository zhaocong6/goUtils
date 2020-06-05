package redislock

import (
	"github.com/gomodule/redigo/redis"
	"math/rand"
)

type Locker struct {
	conn   redis.Conn
	expire int
	key    string
	val    uint64
}

var DefaultExpire = 20

//创建一个lock
func New(redis redis.Conn, key string, expire int) *Locker {
	if expire < 0 {
		expire = DefaultExpire
	}

	return &Locker{
		conn:   redis,
		expire: expire,
		key:    key,
		val:    uint64(rand.Int63()),
	}
}

//使用redis自带的set命令实现原子加锁
func (l *Locker) TryLock() error {
	_, err := redis.String(l.conn.Do("SET", l.key, l.val, "EX", l.expire, "NX"))

	if err != nil {
		return err
	}

	return nil
}

//使用lua保证代码执行原子性
var lockScript = redis.NewScript(1, `
if redis.call("get",KEYS[1]) == ARGV[1] then
	return redis.call("del",KEYS[1])
else
	return 0
end
`)

func (l *Locker) Unlock() error {
	_, err := redis.String(lockScript.Do(l.conn, l.key, l.val))

	if err != nil {
		return err
	}

	return nil
}
