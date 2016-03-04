package main

import (
	"fmt"
	"time"

	"github.com/coseyo/vcache"
)

func main() {
	err := vcache.InitRedis("tcp", "10.20.187.251:11311", 30, 900*time.Second)
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

	// lock use
	ok, err := cache.Lock("aa", 30, 20)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("one", ok)

	ok, _ = cache.Lock("aa", 10, 20)

	fmt.Println("two", ok)
}
