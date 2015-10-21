package vcache

import (
	"time"

	"github.com/coseyo/radixpool"
)

var (
	redisPool *radixpool.Pool
)

// init redis config
func InitRedis(network, addr string, size int, clientTimeout time.Duration) error {
	var err error
	redisPool, err = radixpool.NewPool(network, addr, size, clientTimeout)
	return err
}

func get(key string) (string, error) {
	return redisPool.Cmd("GET", key).Str()
}

func set(key, value string) error {
	return redisPool.Cmd("SET", key, value).Err
}

func del(key string) error {
	return redisPool.Cmd("DEL", key).Err
}

func expire(key string, seconds int) error {
	return redisPool.Cmd("EXPIRE", key, seconds).Err
}

func incr(key string) (int, error) {
	return redisPool.Cmd("INCR", key).Int()
}
