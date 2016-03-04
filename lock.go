package vcache

import (
	"errors"
	"time"
)

// Lock will lock the key in lockSecond secoend and expire in expireSecond
func (this *VCache) Lock(key string, lockSecond, expireSecond int) (ok bool, err error) {
	rc, err := redisPool.Get()
	if err != nil {
		return
	}
	defer redisPool.CarefullyPut(rc, &err)

	if lockSecond > expireSecond {
		err = errors.New("lockedSecond should not greater than expireSecond")
		return
	}

	key = this.getKey(key)
	curTime := int(time.Now().Unix())
	expireTime := curTime + lockSecond + 1

	rs, err := rc.Conn.Cmd("SETNX", key, expireTime).Int()
	if err != nil {
		return
	}

	if rs == 1 {
		err = rc.Conn.Cmd("EXPIRE", key, expireSecond).Err
		ok = true
		return
	}

	lockTime, _ := rc.Conn.Cmd("GET", key).Int()
	oldLockTime, _ := rc.Conn.Cmd("GETSET", key, expireTime).Int()
	if curTime > lockTime && curTime > oldLockTime {
		err = rc.Conn.Cmd("EXPIRE", key, expireSecond).Err
		ok = true
		return
	}

	return
}

// unlock will unlock the key
func (this *VCache) UnLock(key string) (err error) {
	rc, err := redisPool.Get()
	if err != nil {
		return
	}
	defer redisPool.CarefullyPut(rc, &err)

	key = this.getKey(key)
	curTime := int(time.Now().Unix())
	lockTime, err := rc.Conn.Cmd("GET", key).Int()
	if err != nil || lockTime == 0 {
		return
	}
	if curTime < lockTime {
		err = rc.Conn.Cmd("DEL", key).Err
	}
	return
}
