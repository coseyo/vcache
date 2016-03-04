package vcache

import (
	"time"
)

// http://blog.csdn.net/lihao21/article/details/49104695
func (this *VCache) Lock(key string, expireSecond int) (ok bool, err error) {
	rc, err := redisPool.Get()
	if err != nil {
		return
	}
	defer redisPool.CarefullyPut(rc, &err)

	key = this.getKey(key)
	expireTime := int(time.Now().Unix()) + expireSecond + 1

	rs, err := rc.Conn.Cmd("SETNX", key, expireTime).Int()
	if err != nil {
		return
	}

	if rs == 1 {
		ok = true
	} else {

		ok = false
	}
	return
}

func (this *VCache) UnLock(key string) error {
	return del(this.getKeyWithVersionNum(key))
}
