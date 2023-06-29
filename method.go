package gcache

import (
	"sync"
	"time"
)

// Put save a new k-v data whether it is existing.
// If it existed in singleCache before then update it.
// flush the lru list to make the key in the tail of lru list cause you use it.
// Use NotExpire the key will not be expired
// ex is an option param which set key lifetime and only the first one value is valid
func (c *Cache[K]) Put(k K, v any, ex ...time.Duration) {
	defer c.eliminate(c.fs)
	if _cache, _ := c.getNode(k); _cache != nil {
		_cache.put(k, v, ex...)
	}

}

// Get query data by the key you set.
// If not exist will return false in second return value.
// flush the lru list to make the key in the tail of lru list cause you use it.
func (c *Cache[K]) Get(k K) (any, bool) {
	if _cache, _ := c.getNode(k); _cache != nil {
		return _cache.get(k)
	}
	return nil, false
}

func (c *Cache[K]) ExpireAt(k K) (time.Time, bool) {
	if _cache, _ := c.getNode(k); _cache != nil {
		return _cache.expireAt(k)
	}
	return NoExpire, false
}

// Keys return all keys saved in the singleCache.
// This method not flush lru list.
func (c *Cache[K]) Keys() []K {
	ks := make([]K, 0, c.Size())

	var wait sync.WaitGroup
	for _, _cache := range c.buckets {
		wait.Add(1)
		go func(_cache *cache[K]) {
			_cache.Lock()
			defer wait.Done()
			defer _cache.Unlock()
			for k, v := range _cache.data {
				if v.expireAt != NoExpire && v.expireAt.Before(time.Now()) {
					go func() {
						_cache.ex <- k
					}()
					continue
				}
				ks = append(ks, k)
			}
		}(_cache)
	}
	wait.Wait()

	return ks
}

// Del delete keys from singleCache.
// remove keys from lru list which key is valid.
func (c *Cache[K]) Del(ks ...K) {
	if l := len(ks); l > 0 {
		for i := 0; i < l; i++ {
			if _cache, _ := c.getNode(ks[i]); _cache != nil {
				_cache.del(ks[i])
			}
		}
	}
}
