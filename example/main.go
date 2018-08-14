package main

import (
	"fmt"
	"time"

	"log"
	"sync"

	"github.com/coseyo/vcache"
)

func main() {
	err := vcache.InitRedis("tcp", "10.13.88.102:11311", 30, 900*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	vcache.GlobalKeyPrefix = "globalVcache8"

	cache := vcache.New("test", 900)

	// versionParams is  use to generate the version key, not including the page param
	versionParams := map[string]interface{}{
		"username": "test1",
	}

	// get page 1 on specific condition
	a := map[string]interface{}{
		"username": "test1",
		"state":    1,
		"page":     1,
	}

	// get page 2 on specific condition
	b := map[string]interface{}{
		"username": "test1",
		"state":    1,
		"page":     2,
	}

	// set version key
	cache.SetVersionKey(versionParams)

	// generate key by params
	keya := cache.GenerateKey(a, "prefix_aa", "prefix_aa_2")
	if err := cache.Set(keya, a); err != nil {
		fmt.Println(err)
	}

	keyb := cache.GenerateKey(b, "prefix_bb")
	if err := cache.Set(keyb, b); err != nil {
		fmt.Println(err)
	}

	//	cache.Del(keya)

	// test the cache data
	value, _ := cache.Get(keya)
	fmt.Println("keya vaule is", value)
	value, _ = cache.Get(keyb)
	fmt.Println("keyb vaule is", value)

	// The editor change some content, want to refresh the page 1
	// and page 2 immediately, may be much more pages.
	// And just execute the IncrVersionNum() method, the cache will be deprecated
	cache.IncrVersionNum()

	// because the version num was changed, the data is null
	value, _ = cache.Get(keya)
	fmt.Println("keya vaule is", value)
	value, _ = cache.Get(keyb)
	fmt.Println("keyb value is", value)

	//bb, _ := vcache.RedisPool.Cmd("SETNX", "BB", 1).Int()
	//log.Println("bb", bb)
	//
	//bb2, _ := vcache.RedisPool.Cmd("GET", "BB").Int()
	//log.Println("bb2", bb2)
	//
	//bb, _ = vcache.RedisPool.Cmd("SETNX", "BB", 2).Int()
	//log.Println("bb", bb)
	//
	//bb2, _ = vcache.RedisPool.Cmd("GET", "BB").Int()
	//log.Println("bb2", bb2)

	//lock use

	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			ok, err := cache.SLock("dasddsda", 100)
			if ok {
				fmt.Println(i, ok, err, "========================", time.Now().Unix())
			}
			wg.Done()
		}()
	}

	wg.Wait()

	log.Println("finish")
}
