package gcache

// Put save a new k-v data whether it is existing.
// If it existed in cache before then update it.
// flush the lru list to make the key in the tail of lru list cause you use it.
func (c *Cache[K]) Put(k K, v interface{}) {
	c.cache.Lock()
	defer c.eliminate(c.cache.fs)
	defer c.cache.Unlock()
	defer c.lru.flush(k)
	c.cache.data[k] = v
	c.cache.size = len(c.cache.data)
}

// Get query data by the key you set.
// If not exist will return false in second return value.
// flush the lru list to make the key in the tail of lru list cause you use it.
func (c *Cache[K]) Get(k K) (interface{}, bool) {
	c.cache.RLock()
	defer c.cache.RUnlock()
	v, ok := c.cache.data[k]
	if ok {
		c.lru.flush(k)
	}
	return v, ok
}

// Keys return all keys saved in the cache.
// This method not flush lru list.
func (c *Cache[K]) Keys() []K {
	c.cache.RLock()
	defer c.cache.RUnlock()
	ks := make([]K, 0, len(c.cache.data))
	for k := range c.cache.data {
		ks = append(ks, k)
	}
	return ks
}

// Del delete keys from cache.
// remove keys from lru list which key is valid.
func (c *Cache[K]) Del(ks ...K) {
	c.cache.Lock()
	defer c.cache.Unlock()
	if l := len(ks); l > 0 {
		for i := 0; i < l; i++ {
			if _, ok := c.cache.data[ks[i]]; ok {
				delete(c.cache.data, ks[i])
				c.lru.remove(ks[i])
			}
		}
		c.cache.size = len(c.cache.data)
	}
}

// eliminate Trigger elimination to clear the keys that haven't used in longest time duration.
// When cached data size is greater than set max you set
// Clear data count(free up space) decided by percentage of max count
// Clear until cached data size less than max even if cache.fs is not a valid value(0-1)
// WARN: if max you set is not greater than 0 this strategy won't be triggered
func (c *Cache[K]) eliminate(fs float64) {
	if c.cache.size <= c.cache.max || c.cache.max == 0 {
		return
	}
	c.cache.Lock()
	defer c.cache.Unlock()

	if fs > 0 {
		for c.cache.size > c.cache.max {
			for i := float64(0); i < float64(c.cache.max)*fs; i++ {
				if eldK, ok := c.lru.popHead(); ok {
					delete(c.cache.data, eldK)
				}
			}
			c.cache.size = len(c.cache.data)
		}
	} else {
		for c.cache.size > c.cache.max {
			if eldK, ok := c.lru.popHead(); ok {
				delete(c.cache.data, eldK)
				c.cache.size = len(c.cache.data)
			}
		}
	}
}
