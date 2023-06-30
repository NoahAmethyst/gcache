package gcache

import (
	"encoding/binary"
	"hash/crc32"
	"math"
	"sort"
	"sync"
	"time"
)

// del delete key from cache.
// remove key from lru list which key is valid
func (c *cache[K]) del(k K) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.data[k]; ok {
		delete(c.data, k)
		c.lru.remove(k)
		c.size = len(c.data)
	}
}

// eliminate Trigger elimination to clear the keys that haven't used in longest time duration.
// When cached data size is greater than set max you set
// Clear data count(free up space) decided by percentage of max count
// Clear until cached data size less than max even if cache.fs is not a valid value(0-1)
// WARN: if max you set is not greater than 0 this strategy won't be triggered
func (c *Cache[K]) eliminate(fs float64) {
	if c.Size() <= c.max || c.max == 0 {
		return
	}
	var wait sync.WaitGroup

	for _, _cache := range c.buckets {
		wait.Add(1)
		go func(_cache *cache[K]) {
			_cache.Lock()
			defer wait.Done()
			defer _cache.Unlock()
			if fs > 0 {
				for _cache.size > c.max {
					for i := float64(0); i < float64(c.max)*fs; i++ {
						if eldK, ok := _cache.lru.popHead(); ok {
							delete(_cache.data, eldK)
						}
					}
					_cache.size = len(_cache.data)
				}
			} else {
				for _cache.size > c.max {
					if eldK, ok := _cache.lru.popHead(); ok {
						delete(_cache.data, eldK)
						_cache.size = len(_cache.data)
					}
				}
			}
		}(_cache)
	}
	wait.Wait()

}

// addNode add bucket node
func (c *Cache[K]) addNode() {
	c.Lock()
	defer c.Unlock()
	b := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(b, time.Now().UnixMilli())
	hash := crc32.ChecksumIEEE(b)
	c.nodes = append(c.nodes, hash)
	c.buckets[hash] = &cache[K]{
		data:    map[K]*item{},
		lru:     &lruLink[K]{},
		ex:      make(chan K),
		size:    0,
		RWMutex: sync.RWMutex{},
	}
	sort.Slice(c.nodes, func(i, j int) bool {
		return c.nodes[i] < c.nodes[j]
	})

	go func(cache *cache[K]) {
		for {
			if cache == nil {
				continue
			}
			select {
			case k := <-cache.ex:
				cache.del(k)
			}
		}
	}(c.buckets[hash])
}

// reloadData reload bucket data cause new node added
func (c *Cache[K]) reloadData(hash uint32) {
	if len(c.nodes) == 1 {
		return
	}
	index := 0
	for i, _hash := range c.nodes {
		if _hash == _hash {
			index = i
			break
		}
	}

	var _cache *cache[K]
	currCache := c.buckets[hash]
	if index == 0 {
		_cache = c.buckets[c.nodes[index+1]]
	} else {
		_cache = c.buckets[c.nodes[index-1]]
	}

	for k, v := range _cache.data {
		if _, _hash := c.getNode(k); _hash == hash {
			currCache.data[k] = v
			currCache.size = len(currCache.data)
			currCache.lru.flush(k)
			_cache.del(k)
		}
	}
}

// getNode get bucket node by key
func (c *Cache[K]) getNode(key K) (*cache[K], uint32) {
	if len(c.nodes) == 0 {
		return nil, 0
	}
	var hash uint32
	switch any(key).(type) {
	case string:
		hash = crc32.ChecksumIEEE([]byte(any(key).(string)))
	case int:
		ibytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(ibytes, uint32(any(key).(int)))
		hash = crc32.ChecksumIEEE(ibytes)
	case int64:
		i64bytes := make([]byte, 8)
		binary.PutVarint(i64bytes, any(key).(int64))
		hash = crc32.ChecksumIEEE(i64bytes)
	case float64:
		fbytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(fbytes, math.Float64bits(any(key).(float64)))
		hash = crc32.ChecksumIEEE(fbytes)
	}

	idx := sort.Search(len(c.nodes), func(i int) bool {
		return c.nodes[i] >= hash
	})
	if idx == len(c.nodes) {
		idx = 0
	}

	return c.buckets[c.nodes[idx]], c.nodes[idx]
}

func (c *cache[K]) put(k K, v any, ex ...time.Duration) {
	c.Lock()
	defer c.Unlock()
	defer c.lru.flush(k)
	_item := &item{
		v: v,
	}
	if len(ex) > 0 {
		if ex[0] == NotExpire {
			_item.expireAt = NoExpire
		} else {
			expireAt := time.Now().Add(ex[0])
			_item.expireAt = expireAt
		}
	} else {
		_item.expireAt = NoExpire
	}
	c.data[k] = _item
	c.size = len(c.data)
}

func (c *cache[K]) get(k K) (any, bool) {
	c.RLock()
	defer c.RUnlock()
	v, ok := c.data[k]
	if ok {
		if v.expireAt != NoExpire && v.expireAt.Before(time.Now()) {
			go func() {
				c.ex <- k
			}()
			return nil, false
		}
		c.lru.flush(k)
		return v.v, ok
	}
	return nil, ok
}

func (c *cache[K]) expireAt(k K) (time.Time, bool) {
	c.RLock()
	defer c.RUnlock()
	v, ok := c.data[k]
	if !ok {
		return NoExpire, ok
	}
	return v.expireAt, ok

}
