package mem

import (
	"sync"
	"time"

	"github.com/impact-eintr/ecache/cache"
)

type value struct {
	v       []byte
	created time.Time
}

type memCache struct {
	c     map[string]value
	mutex sync.RWMutex
	cache.Stat
	ttl time.Duration
}

func (mc *memCache) Set(k string, v []byte) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	tmp, exist := mc.c[k]
	if exist {
		mc.Remove(k, tmp.v)
	}

	mc.c[k] = value{v, time.Now()}
	mc.Add(k, v)
	return nil

}

func (mc *memCache) Get(k string) ([]byte, error) {
	if mc.ttl != 0 {
		mc.mutex.Lock()
		val := mc.c[k].v
		mc.c[k] = value{val, time.Now()} // 更新缓存

		defer mc.mutex.Unlock()
		return val, nil

	} else {
		mc.mutex.RLock()
		defer mc.mutex.RUnlock()
		return mc.c[k].v, nil

	}
}

func (mc *memCache) Del(k string) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	v, exist := mc.c[k]
	if exist {
		delete(mc.c, k)
		mc.Remove(k, v.v) //修改状态
	}
	return nil

}

func (mc *memCache) GetStat() cache.Stat {
	return mc.Stat

}

func NewCache(ttl int) *memCache {
	mc := &memCache{
		make(map[string]value),
		sync.RWMutex{},
		cache.Stat{},
		time.Duration(ttl) * time.Second}

	if ttl > 0 {
		go mc.expirer()
	}

	return mc

}

func (mc *memCache) expirer() {
	for {
		time.Sleep(mc.ttl)
		mc.mutex.Lock()
		for k, v := range mc.c {
			mc.mutex.Unlock()
			if v.created.Add(mc.ttl).Before(time.Now()) {
				mc.Del(k)
			}
			mc.mutex.Lock()
		}
		mc.mutex.Unlock()
	}
}
