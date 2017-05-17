package vcache

import (
	"time"

	"github.com/coseyo/radixpool"
)

var (
	RedisPool *radixpool.Pool
)

// init redis config
func InitRedis(network, addr string, size int, clientTimeout time.Duration) error {
	var err error
	RedisPool, err = radixpool.NewPool(network, addr, size, clientTimeout)
	return err
}

func get(key string) (string, error) {
	return RedisPool.Cmd("GET", key).Str()
}

func set(key, value string) error {
	return RedisPool.Cmd("SET", key, value).Err
}

func setex(key string, seconds int, value string) error {
	return RedisPool.Cmd("SETEX", key, seconds, value).Err
}

func del(key string) error {
	return RedisPool.Cmd("DEL", key).Err
}

func expire(key string, seconds int) error {
	return RedisPool.Cmd("EXPIRE", key, seconds).Err
}

func incr(key string) (int, error) {
	return RedisPool.Cmd("INCR", key).Int()
}

func decr(key string) (int, error) {
	return RedisPool.Cmd("DECR", key).Int()
}
