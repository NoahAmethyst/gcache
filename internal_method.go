package gcache

// del delete key from cache.
// remove key from lru list which key is valid.
func (c *Cache[K]) del(k K) {
	c.singleCache.Lock()
	defer c.singleCache.Unlock()
	if _, ok := c.singleCache.data[k]; ok {
		delete(c.singleCache.data, k)
		c.lru.remove(k)
		c.singleCache.size = len(c.singleCache.data)
	}
}

// eliminate Trigger elimination to clear the keys that haven't used in longest time duration.
// When cached data size is greater than set max you set
// Clear data count(free up space) decided by percentage of max count
// Clear until cached data size less than max even if cache.fs is not a valid value(0-1)
// WARN: if max you set is not greater than 0 this strategy won't be triggered
func (c *Cache[K]) eliminate(fs float64) {
	if c.singleCache.size <= c.max || c.max == 0 {
		return
	}
	c.singleCache.Lock()
	defer c.singleCache.Unlock()

	if fs > 0 {
		for c.singleCache.size > c.max {
			for i := float64(0); i < float64(c.max)*fs; i++ {
				if eldK, ok := c.lru.popHead(); ok {
					delete(c.singleCache.data, eldK)
				}
			}
			c.singleCache.size = len(c.singleCache.data)
		}
	} else {
		for c.singleCache.size > c.max {
			if eldK, ok := c.lru.popHead(); ok {
				delete(c.singleCache.data, eldK)
				c.singleCache.size = len(c.singleCache.data)
			}
		}
	}
}
