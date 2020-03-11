package vcache

import (
	"fmt"
	"github.com/json-iterator/go"

	"errors"

	"github.com/coseyo/goutil/sortmap"
	"github.com/coseyo/vcache/util"
)

var (
	GlobalKeyPrefix         string
	DefaultVersionKeyExpire int = 7200
)

const (
	ErrCacheEmpty = "CACHE_EMPTY"
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
	key = this.GetKeyWithVersionNum(key)
	data, err := get(key)
	if data == "" {
		return nil, err
	}
	return util.JsonDecode(data)
}

// get cache data string by key
func (this *VCache) GetString(key string) (string, error) {
	key = this.GetKeyWithVersionNum(key)
	return get(key)
}

// GetByType empty cache will return error
func (this *VCache) GetByType(key string, v interface{}) (err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	key = this.GetKeyWithVersionNum(key)
	str, err := get(key)
	if err != nil {
		return
	}
	if str == "" {
		err = errors.New(ErrCacheEmpty)
		return
	}
	err = json.Unmarshal([]byte(str), v)
	return
}

// GetByTypeWithExist return exist and error
func (this *VCache) GetByTypeWithExist(key string, v interface{}) (exist bool, err error) {
	err = this.GetByType(key, v)
	if err == nil {
		exist = true
		return
	}

	if err.Error() == ErrCacheEmpty {
		err = nil
	}
	return
}

// set cache data
func (this *VCache) Set(key string, value interface{}) error {
	key = this.GetKeyWithVersionNum(key)
	data, err := util.JsonEncode(value)
	if err != nil {
		return err
	}
	return setex(key, this.ExpireSecond, data)
}

// set cache data with expire second
func (this *VCache) SetWithExpire(key string, value interface{}, expireSecond int) error {
	key = this.GetKeyWithVersionNum(key)
	data, err := util.JsonEncode(value)
	if err != nil {
		return err
	}
	return setex(key, expireSecond, data)
}

// delete cache data
func (this *VCache) Del(key string) error {
	return del(this.GetKeyWithVersionNum(key))
}

// incr cache data
func (this *VCache) Incr(key string) (int, error) {
	return incr(this.GetKeyWithVersionNum(key))
}

// incr cache data
func (this *VCache) Decr(key string) (int, error) {
	return decr(this.GetKeyWithVersionNum(key))
}

// set expireTime cache data
func (this *VCache) Expire(key string, expireSecond int) error {
	return expire(this.GetKeyWithVersionNum(key), expireSecond)
}

// increase cache version num
func (this *VCache) IncrVersionNum() error {
	key := this.getKey(this.versionKey)
	_, err := incr(key)
	if err == nil {
		expire(key, DefaultVersionKeyExpire)
	}
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
	key := sortmap.MapToMD5String(params)
	for _, v := range prefix {
		key = v + ":" + key
	}
	return key
}

func (this *VCache) GetKey(key string) string {
	return fmt.Sprintf("%s:%s:%s", GlobalKeyPrefix, this.KeyPrefix, util.MD5(key))
}

func (this *VCache) GetKeyWithVersionNum(key string) string {
	key = this.GetKey(key)
	if this.versionKey != "" {
		key += ":" + this.getVersionNum()
	}
	return key
}
