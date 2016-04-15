package vcache

import (
	"errors"
	"time"
)

const (
	CLOCK_PREFIX = "clock:"
	SLOCK_PREFIX = "slock:"
)

// CLock is a concurrent lock, it will lock the key in lockSecond secoend and expire in expireSecond
// When CLock the same key in lockSecond, it will extend the lock time to keep it exclusive in concurrent env.
func (this *VCache) CLock(key string, lockSecond, expireSecond int) (ok bool, err error) {
	rc, err := redisPool.Get()
	if err != nil {
		return
	}
	defer redisPool.CarefullyPut(rc, &err)

	if lockSecond > expireSecond {
		err = errors.New("lockedSecond should not greater than expireSecond")
		return
	}

	key = this.getKey(CLOCK_PREFIX + key)
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

// SLock is a sequencial lock, it will lock the key in lockSecond, just like a normal lock.
func (this *VCache) SLock(key string, lockSecond int) (ok bool, err error) {
	rc, err := redisPool.Get()
	if err != nil {
		return
	}
	defer redisPool.CarefullyPut(rc, &err)

	key = this.getKey(SLOCK_PREFIX + key)
	rs, err := rc.Conn.Cmd("SETNX", key, 1).Int()
	if err != nil {
		return
	}

	if rs == 1 {
		err = rc.Conn.Cmd("EXPIRE", key, lockSecond).Err
		ok = true
	}

	return
}

// UnCLock will unlock CLock
func (this *VCache) UnCLock(key string) (err error) {
	rc, err := redisPool.Get()
	if err != nil {
		return
	}
	defer redisPool.CarefullyPut(rc, &err)

	key = this.getKey(CLOCK_PREFIX + key)
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

// UnSLock will unlock SLock
func (this *VCache) UnSLock(key string) (err error) {
	rc, err := redisPool.Get()
	if err != nil {
		return
	}
	defer redisPool.CarefullyPut(rc, &err)

	key = this.getKey(SLOCK_PREFIX + key)
	err = rc.Conn.Cmd("DEL", key).Err
	return
}
