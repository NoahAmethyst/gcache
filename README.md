# GCache

---

gcache is a lightweight local cache project that supports concurrent access/reading of k-v data
* Concurrency: Test environment: macos 10 core memory limit of 2g. supports concurrent storage of 40,000- data per second and concurrent reading of 60,000 data per second.
* Support elimination policy: Set the maximum data and space release percentage of a single cache server, and use the lru algorithm to eliminate the maximum data percentage when the trigger condition is met
* Support expiration policy: Set expiration time, similar to Redis *Generic key structure. key data type can be specified. Currently, int int64 float64 string is supported 
* value data uses interface, so data needs to be forcibly transferred


**WARN:Data persistence is not supported**

---
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