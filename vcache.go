package vcache

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

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

// init with redis config and return the struct instance, and set default expire time
func Init(network, addr string, size int) (*VCache, error) {
	if err := initPool(network, addr, size); err != nil {
		return nil, err
	}
	return &VCache{}, nil
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
	data, err := get(key)
	if err != nil {
		return nil, err
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

// increase cache version num
func (this *VCache) IncrVersionNum() error {
	_, err := incr(this.getKey(this.versionKey))
	return err
}

// get the version num by version key
func (this *VCache) getVersionNum() string {
	versionNum, _ := get(this.getKey(this.versionKey))
	if versionNum == "" {
		versionNum = "1"
	}
	return versionNum
}

// set version key according to the params, the params should not including the
// unnecessary params,  like the page, offset
func (this *VCache) SetVersionKey(params map[string]interface{}) (err error) {
	for k, v := range params {
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Int {
			this.versionKey += fmt.Sprintf("_%s-%s", k, strconv.Itoa(int(rv.Int())))
		} else if rv.Kind() == reflect.String {
			this.versionKey += fmt.Sprintf("_%s-%s", k, rv.String())
		} else {
			err = errors.New("Invalid type value")
			break
		}
	}
	return
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
