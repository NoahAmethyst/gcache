package gcache

import "sync"

type Cache[K int | int64 | float64 | string] struct {
	cache *cache[K]
	lru   *lruLink[K]
}

type cache[K int | int64 | float64 | string] struct {
	data map[K]interface{}
	size int
	max  int
	fs   float64
	sync.RWMutex
}

// NewCache init a Cache object to save your data
// You can specify [int | int64 | float64 | string] type as your key type,but can't specify multi types.
// option[0] is max data size you will save in cache if the value not set or not greater than zero, the cache will not be lru.
// option[1] set the strategy that free up how many data space when data size greater than max data size you set.
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

	return &Cache[K]{
		cache: &cache[K]{
			data:    map[K]interface{}{},
			size:    0,
			max:     _max,
			fs:      _fs,
			RWMutex: sync.RWMutex{},
		},
		lru: &lruLink[K]{},
	}
}
