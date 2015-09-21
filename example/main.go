package main

import (
	"fmt"

	"github.com/coseyo/vcache"
)

func main() {
	_, err := vcache.Init("tcp", "127.0.0.1:11311", 30)
	if err != nil {
		fmt.Println(err)
		return
	}

	vcache.GlobalKeyPrefix = "globalVcache8"

	cache := vcache.New("test", 900)
	cache2 := vcache.New("test", 900)

	// versionParams is  use to generate the version key, not including the page param
	versionParams := map[string]interface{}{
		"username": "test1",
		"state":    1,
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
	cache2.SetVersionKey(versionParams)

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
	fmt.Println(value)
	value, _ = cache.Get(keyb)
	fmt.Println(value)

	// The editor change some content, want to refresh the page 1
	// and page 2 immediately, may be much more pages.
	// And just execute the IncrVersionNum() method, the cache will be deprecated
	//	cache.IncrVersionNum()
	cache2.IncrVersionNum()

	// because the version num was changed, the data is null
	value, _ = cache.Get(keya)
	fmt.Println(value)
	value, _ = cache2.Get(keya)
	fmt.Println(value)
	value, _ = cache.Get(keyb)
	fmt.Println(value)
}
