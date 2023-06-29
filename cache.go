package gcache

import (
	"sync"
	"time"
)

const (
	// NotExpire not expire
	NotExpire time.Duration = -1
	// threshold add bucket node when data size trigger threshold
	threshold = 1024 * 8
	// replicas virtual node magnification
	replicas = 3
)

var NoExpire time.Time

type Cache[K int | int64 | float64 | string] struct {
	//store hash of cache nodes
	nodes []uint32
	//store map relationship of cache nodes and hash
	buckets map[uint32]*cache[K]
	//max data size that will be stored in cache
	max int
	//percentage of max data size that will be free of
	fs float64

	sync.RWMutex
}

type cache[K int | int64 | float64 | string] struct {
	//saved data
	data map[K]*item
	//store all key and sort by recent used,for single cache
	lru *lruLink[K]
	//expire channel which receive expired keys and delete it
	ex chan K
	//data size
	size int
	sync.RWMutex
}

type item struct {
	//data value
	v any
	//expire time
	expireAt time.Time
}

// NewCache init a Cache object to save your data
// You can specify [int | int64 | float64 | string] type as your key type,but can't specify multi types.
// option[0] is max data size you will save in singleCache if the value not set or not greater than zero, the singleCache will not be lru.
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
		RWMutex: sync.RWMutex{},
		max:     _max,
		fs:      _fs,
		buckets: make(map[uint32]*cache[K]),
	}
	_cache.addNode()
	return _cache

}

// SetMax set max data size singleCache can hold
// Data size not limit when  the value is not greater than zero
func (c *Cache[K]) SetMax(max int) {
	c.max = max
	if max > threshold {
		for i := threshold * len(c.nodes); i < max; {
			c.addNode()
			i = threshold * len(c.nodes)
		}
	}
}

// Max get singleCache current max
func (c *Cache[K]) Max() int {
	return c.max
}

// Size get singleCache current data size
func (c *Cache[K]) Size() int {
	size := 0
	for _, _cache := range c.buckets {
		size += _cache.size
	}
	return size
}

func init() {
	NoExpire = time.Time{}
}
