package gcache

// del delete key from cache.
// remove key from lru list which key is valid.
func (c *Cache[K]) del(k K) {
	c.cache.Lock()
	defer c.cache.Unlock()
	if _, ok := c.cache.data[k]; ok {
		delete(c.cache.data, k)
		c.lru.remove(k)
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
