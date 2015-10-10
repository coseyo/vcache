package vcache

import "github.com/fzzy/radix/extra/pool"

var (
	redisPool *pool.Pool
)

func initPool(network, addr string, size int) error {
	var err error
	redisPool, err = pool.NewPool(network, addr, size)
	return err
}

func get(key string) (string, error) {
	client, redisErr := redisPool.Get()
	defer redisPool.CarefullyPut(client, &redisErr)
	return client.Cmd("GET", key).Str()
}

func set(key, value string) error {
	client, redisErr := redisPool.Get()
	defer redisPool.CarefullyPut(client, &redisErr)
	return client.Cmd("SET", key, value).Err
}

func del(key string) error {
	client, redisErr := redisPool.Get()
	defer redisPool.CarefullyPut(client, &redisErr)
	return client.Cmd("DEL", key).Err
}

func expire(key string, seconds int) error {
	client, redisErr := redisPool.Get()
	defer redisPool.CarefullyPut(client, &redisErr)
	return client.Cmd("EXPIRE", key, seconds).Err
}

func incr(key string) (int, error) {
	client, redisErr := redisPool.Get()
	defer redisPool.CarefullyPut(client, &redisErr)
	return client.Cmd("INCR", key).Int()
}
