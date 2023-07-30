package Cache

import "time"

type Cache interface {
	SetMaxMemory(size string) bool
	Set(key string, val interface{}, expire time.Duration) bool
	Get(key string) (*memCacheValue, bool)
	Del(key string) bool
	Exists(key string) bool
	Flush() bool
	Keys() int64
}
