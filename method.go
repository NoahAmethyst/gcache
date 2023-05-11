package gcache

func (c *Cache[K]) Put(k K, v interface{}) {
	c.cache.Lock()
	defer c.eliminate(c.cache.fs)
	defer c.cache.Unlock()
	defer c.lru.flush(k)
	c.cache.data[k] = v
	c.cache.size = len(c.cache.data)
}

func (c *Cache[K]) Get(k K) (interface{}, bool) {
	c.cache.RLock()
	defer c.cache.RUnlock()
	v, ok := c.cache.data[k]
	if ok {
		c.lru.flush(k)
	}
	return v, ok
}

func (c *Cache[K]) Keys() []K {
	c.cache.RLock()
	defer c.cache.RUnlock()
	ks := make([]K, 0, len(c.cache.data))
	for k := range c.cache.data {
		ks = append(ks, k)
	}
	return ks
}

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
