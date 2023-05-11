# GCache

simple local cache with lru strategy

### USE

```go
go get github.com/NoahAmethyst/gcache
```

**You can watch code in ./gotest see how to use**
**Or see function below:**

Init Cache
```go
import "github.com/NoahAmethyst/gcache"
//New Cache with max data size and percentage of free up size
//It support Generic Key type
//You should specify concrete type you need for the key type
max := 20000
localCache := gcache.NewCache[string](max, 10)

```

Put data into cache
```go
import "github.com/NoahAmethyst/gcache"

max := 20000
localCache := gcache.NewCache[string](max, 10)


//Put data into cache and set expire time
k:="test"
v:="test_data"
localCache.Put(k, v, 10*time.Millisecond)

//Put data into cache without expire
//Option1
k="test_2"
v="test_data_2"
localCache.Put(k, v)

//Put data into cache without expire
//Option2
k="test_2"
v="test_data_2"
localCache.Put(k, v, gcache.NotExpire)
```
Get data from cache
```go
import "github.com/NoahAmethyst/gcache"

max := 20000
localCache := gcache.NewCache[string](max, 10)


//Get data from cache
//ok is false when data(key) not exist or expired
k:= "test"
v, ok := localCache.Get(k)

```

Get data expire time
```go
import "github.com/NoahAmethyst/gcache"

max := 20000
localCache := gcache.NewCache[string](max, 10)


//Get data expire time
//ok is false when data(key) not exist
k:= "test"
expireAt, ok := localCache.ExpireAt(k)
```

Get all keys
```go
import "github.com/NoahAmethyst/gcache"

max := 20000
localCache := gcache.NewCache[string](max, 10)

//Get all keys
//return a slice with exist and not expired keys
keys := localCache.Keys()

```