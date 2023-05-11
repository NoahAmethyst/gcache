package gcache

import (
	"sync"
	"time"
)

const (
	NotExpire time.Duration = -1
)

var NoExpire time.Time

type Cache[K int | int64 | float64 | string] struct {
	cache *cache[K]
	lru   *lruLink[K]
}

type cache[K int | int64 | float64 | string] struct {
	data map[K]*item
	ex   chan K
	size int
	max  int
	fs   float64
	sync.RWMutex
}

type item struct {
	v        interface{}
	expireAt time.Time
}

// NewCache init a Cache object to save your data
// You can specify [int | int64 | float64 | string] type as your key type,but can't specify multi types.
// option[0] is max data size you will save in cache if the value not set or not greater than zero, the cache will not be lru.
// option[1] set the strategy that free up how many data space when data size greater than max data size you set.
// It's a percentage value,for example:if you want free up 10% of max then the value is 10
// Start a channel to receive expired key and delete it when Cache is made
// WARN:refuse to  set option[0] is not a good ideal.
func NewCache[K int | int64 | float64 | string](option ...int) *Cache[K] {
	_max := 0
	_fs := float64(0)
	if len(option) > 0 && option[0] > 0 {
		_max = option[0]
	}
	if len(option) > 1 {
		if option[1] >= 100 {
			_fs = 0.5
		} else {
			_fs = float64(option[1]) / 100
		}

	}

	_cache := &Cache[K]{
		cache: &cache[K]{
			data:    map[K]*item{},
			ex:      make(chan K),
			size:    0,
			max:     _max,
			fs:      _fs,
			RWMutex: sync.RWMutex{},
		},
		lru: &lruLink[K]{},
	}

	go func() {
		for {
			select {
			case k := <-_cache.cache.ex:
				_cache.del(k)
			}
		}
	}()

	return _cache

}

// SetMax set max data size cache can hold
// Data size not limit when  the value is not greater than zero
func (c *Cache[K]) SetMax(max int) {
	c.cache.max = max
}

// Max get cache current max
func (c *Cache[K]) Max() int {
	return c.cache.max
}

// Size get cache current data size
func (c *Cache[K]) Size() int {
	return c.cache.size
}

func init() {
	NoExpire = time.Time{}
}
