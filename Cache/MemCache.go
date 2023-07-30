package Cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type MemCache struct {
	// 最大内存
	maxMemorySize int64
	//当前内存大小
	currentMemorySize int64
	//最大内存字符串表示
	maxMemorySizeStr string
	//缓存键值对
	values map[string]*memCacheValue
	//锁 (读写锁。读可以并发)
	locker sync.RWMutex
	//清除过期缓存时间间隔
	clearExpireTimeInterval time.Duration
}
type memCacheValue struct {
	//value值
	val interface{}
	//过期时间
	expireTime time.Time
	// value大小
	valueSize int64
}

func NewMemCache() Cache {
	mc := &MemCache{
		values:                  make(map[string]*memCacheValue),
		clearExpireTimeInterval: 10 * time.Second,
	}
	//go mc.clearExpired()
	return mc
}
func (mc *MemCache) SetMaxMemory(size string) bool {
	mc.maxMemorySize, mc.maxMemorySizeStr = ParseSize(size)
	log.Println("setMaxMemory：str is " + mc.maxMemorySizeStr)
	log.Println(mc.maxMemorySize)
	return false
}
func (mc *MemCache) Set(key string, val interface{}, expire time.Duration) bool {
	mc.locker.Lock()
	defer mc.locker.Unlock()
	v := &memCacheValue{
		val:        val,
		expireTime: time.Now().Add(expire),
		valueSize:  GetValueSize(val),
	}
	mc.del(key)
	mc.add(key, v)
	if mc.currentMemorySize > mc.maxMemorySize {
		mc.del(key)
		log.Println(fmt.Sprintf("max memory size%s", mc.maxMemorySize))
		panic(fmt.Sprintf("max memory size%s", mc.maxMemorySize))
	}
	return true
}
func (mc *MemCache) Get(key string) (*memCacheValue, bool) {
	mc.locker.RLock()
	defer mc.locker.RUnlock()
	val, ok := mc.get(key)
	return val, ok
}
func (mc *MemCache) Del(key string) bool {
	mc.locker.Lock()
	defer mc.locker.Unlock()
	mc.del(key)
	return true
}
func (mc *MemCache) del(key string) bool {
	tmp, ok := mc.get(key)
	if ok && tmp != nil {
		mc.currentMemorySize -= tmp.valueSize
		delete(mc.values, key)
	}
	return false
}
func (mc *MemCache) Exists(key string) bool {
	_, ok := mc.Get(key)
	if ok {
		return true
	}
	return false
}
func (mc *MemCache) Flush() bool {
	mc.locker.Lock()
	defer mc.locker.Unlock()
	mc.values = make(map[string]*memCacheValue)
	mc.currentMemorySize = 0
	return true
}
func (mc *MemCache) Keys() int64 {
	mc.locker.RLock()
	mc.locker.RUnlock()
	return int64(len(mc.values))
}

func (mc *MemCache) add(key string, val *memCacheValue) {
	mc.values[key] = val
	mc.currentMemorySize += val.valueSize
}

func (mc *MemCache) get(key string) (*memCacheValue, bool) {
	val, ok := mc.values[key]
	// 判断超时
	if ok {
		if val.expireTime.Before(time.Now()) {
			mc.Del(key)
			log.Println(fmt.Sprintf("the value of the key%s is expired", key))
			return nil, false
		}
		return val, ok
	}

	return nil, false
}

func (mc *MemCache) clearExpired() {
	// 定时触发器
	timeTicker := time.NewTicker(mc.clearExpireTimeInterval)
	defer timeTicker.Stop()
	for {
		select {
		case <-timeTicker.C:
			for key, item := range mc.values {
				if time.Now().After(item.expireTime) {
					mc.locker.Lock()
					mc.del(key)
					mc.locker.Unlock()
				}
			}
		}
	}
}
