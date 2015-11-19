package vcache

import (
	"fmt"

	"github.com/coseyo/vcache/util"
)

var (
	GlobalKeyPrefix string
)

type VCache struct {
	KeyPrefix    string
	ExpireSecond int

	versionKey string
}

// new a VCache instance, thread safe
func New(keyPrefix string, expireSecond int) *VCache {
	instance := new(VCache)
	instance.KeyPrefix = keyPrefix
	instance.ExpireSecond = expireSecond
	return instance
}

// get cache data by key
func (this *VCache) Get(key string) (interface{}, error) {
	key = this.getKeyWithVersionNum(key)
	data, _ := get(key)
	if data == "" {
		return nil, nil
	}
	return util.JsonDecode(data)
}

// set cache data
func (this *VCache) Set(key string, value interface{}) error {
	key = this.getKeyWithVersionNum(key)
	data, err := util.JsonEncode(value)
	if err != nil {
		return err
	}
	expire(key, this.ExpireSecond)
	return set(key, data)
}

// delete cache data
func (this *VCache) Del(key string) error {
	return del(this.getKeyWithVersionNum(key))
}

// incr cache data
func (this *VCache) Incr(key string) (int, error) {
	return incr(this.getKeyWithVersionNum(key))
}

// incr cache data
func (this *VCache) Decr(key string) (int, error) {
	return decr(this.getKeyWithVersionNum(key))
}

// set expireTime cache data
func (this *VCache) Expire(key string, expireSecond int) error {
	return expire(this.getKeyWithVersionNum(key), expireSecond)
}

// increase cache version num
func (this *VCache) IncrVersionNum() error {
	_, err := incr(this.getKey(this.versionKey))
	return err
}

// get the version num by version key
func (this *VCache) getVersionNum() string {
	versionNum, _ := get(this.getKey(this.versionKey))
	if versionNum == "" {
		versionNum = "0"
	}
	return versionNum
}

// set version key according to the params, the params should not including the
// unnecessary params,  like the page, offset
func (this *VCache) SetVersionKey(params map[string]interface{}) *VCache {
	this.versionKey = this.GenerateKey(params)
	return this
}

// generate key by params
func (this *VCache) GenerateKey(params map[string]interface{}, prefix ...string) string {
	var sm util.SortedMaps
	sortedParams := sm.Sort(params)
	jsonString, _ := util.JsonEncode(sortedParams)
	key := jsonString
	for _, v := range prefix {
		key = v + key
	}
	return key
}

func (this *VCache) getKey(key string) string {
	return fmt.Sprintf("%s:%s:%s", GlobalKeyPrefix, this.KeyPrefix, util.MD5(key))
}

func (this *VCache) getKeyWithVersionNum(key string) string {
	key = this.getKey(key)
	if this.versionKey != "" {
		key += ":" + this.getVersionNum()
	}
	return key
}
