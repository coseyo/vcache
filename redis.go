package vcache

import (
	"github.com/fzzy/radix/extra/pool"
	"github.com/fzzy/radix/redis"
)

var (
	redisPool *pool.Pool
)

func getClient() *redis.Client {
	client, _ := redisPool.Get()
	return client
}

func initPool(network, addr string, size int) error {
	var err error
	redisPool, err = pool.NewPool(network, addr, size)
	return err
}

func get(key string) (string, error) {
	return getClient().Cmd("GET", key).Str()
}

func set(key, value string) error {
	return getClient().Cmd("SET", key, value).Err
}

func del(key string) error {
	return getClient().Cmd("DEL", key).Err
}

func expire(key string, seconds int) error {
	return getClient().Cmd("EXPIRE", key, seconds).Err
}

func incr(key string) (int, error) {
	return getClient().Cmd("INCR", key).Int()
}
