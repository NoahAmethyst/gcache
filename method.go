package gcache

import "time"

// Put save a new k-v data whether it is existing.
// If it existed in cache before then update it.
// flush the lru list to make the key in the tail of lru list cause you use it.
// Use NotExpire the key will not be expired
// ex is an option param which set key lifetime and only the first one value is valid
func (c *Cache[K]) Put(k K, v interface{}, ex ...time.Duration) {
	c.cache.Lock()
	defer c.eliminate(c.cache.fs)
	defer c.cache.Unlock()
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
	c.cache.data[k] = _item
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
		if v.expireAt != NoExpire && v.expireAt.Before(time.Now()) {
			go func() {
				c.cache.ex <- k
			}()
			return nil, false
		}
		c.lru.flush(k)
		return v.v, ok
	}
	return nil, ok
}

func (c *Cache[K]) ExpireAt(k K) (time.Time, bool) {
	c.cache.RLock()
	defer c.cache.RUnlock()
	v, ok := c.cache.data[k]
	if !ok {
		return NoExpire, ok
	}
	return v.expireAt, ok
}

// Keys return all keys saved in the cache.
// This method not flush lru list.
func (c *Cache[K]) Keys() []K {
	c.cache.RLock()
	defer c.cache.RUnlock()
	ks := make([]K, 0, len(c.cache.data))
	for k, v := range c.cache.data {
		if v.expireAt != NoExpire && v.expireAt.Before(time.Now()) {
			go func() {
				c.cache.ex <- k
			}()
			continue
		}
		ks = append(ks, k)
	}
	return ks
}

// Del delete keys from cache.
// remove keys from lru list which key is valid.
func (c *Cache[K]) Del(ks ...K) {
	if l := len(ks); l > 0 {
		for i := 0; i < l; i++ {
			c.del(ks[i])
		}
	}
}

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
