# vcache
A cache middleware using redis. Support version cotroller.

When we change some content and want to update all pages's cache, it is hard work to do that in normal way. So we can use cache version to control the cache key. So this lib is a implement by go.

Example below:
```go
package main

import (
	"fmt"

	"github.com/coseyo/vcache"
)

func main() {
	cache, err := vcache.Init("tcp", "127.0.0.1:11311", 30)
	if err != nil {
		fmt.Println(err)
		return
	}

	cache.KeyPrefix = "testPrefix"

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
	if err := cache.SetVersionKey(versionParams); err != nil {
		fmt.Println(err)
	}

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
	cache.IncrVersionNum()

	// because the version num was changed, the data is null
	value, _ = cache.Get(keya)
	fmt.Println(value)
	value, _ = cache.Get(keyb)
	fmt.Println(value)
}

```
