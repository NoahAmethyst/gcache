package gotest

import (
	"github.com/NoahAmethyst/gcache"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test_ConcurrentPut(t *testing.T) {
	max := 20000
	localCache := gcache.NewCache[int](max, 10)

	start := time.Now().UnixMilli()

	var wait sync.WaitGroup

	for i := 0; i < max*2; i++ {
		wait.Add(1)
		go func(k int, v string) {
			defer wait.Done()
			localCache.Put(k, v, time.Duration(max*2-i)*time.Millisecond)
		}(i, strconv.Itoa(i))
	}

	wait.Wait()

	end := time.Now().UnixMilli()

	time.Sleep(10 * time.Second)

	t.Logf("data size:%d,consume:%dms", len(localCache.Keys()), end-start)
}

func Test_BatchPut(t *testing.T) {
	max := 5000
	localCache := gcache.NewCache[int](max, 10)

	for i := 0; i < max*2; i++ {
		localCache.Put(i, strconv.Itoa(i), time.Duration(max*2-i)*time.Millisecond)
	}

	time.Sleep(3 * time.Second)

	keys := localCache.Keys()

	t.Logf("data size:%d", len(keys))
}

func Test_PutWithExpire(t *testing.T) {
	localCache := gcache.NewCache[int](0, 10)

	size := 50
	for i := 0; i < size; i++ {
		if i < size/2 {
			localCache.Put(i, strconv.Itoa(i), time.Duration(i)*time.Second)
		} else {
			localCache.Put(i, strconv.Itoa(i), gcache.NotExpire)
		}
	}

	time.Sleep(5 * time.Second)

	for _, k := range localCache.Keys() {
		v, _ := localCache.Get(k)
		expireAt, _ := localCache.ExpireAt(k)
		var exS string
		if expireAt == gcache.NoExpire {
			exS = "NoExpire"
		} else {
			exS = expireAt.Format("2006-01-02 15:04")
		}

		t.Logf("key:%v value:%v,expireAt:%s", k, v, exS)
	}
}

func Test_ConcurrentGetKey(t *testing.T) {
	max := 20000
	localCache := gcache.NewCache[int](max, 10)

	var wait sync.WaitGroup

	for i := 0; i < max*2; i++ {
		wait.Add(1)
		go func(k int, v string) {
			defer wait.Done()
			localCache.Put(k, v, time.Duration(max*2-i)*time.Millisecond)
		}(i, strconv.Itoa(i))
	}

	wait.Wait()

	start := time.Now().UnixMilli()

	keys := localCache.Keys()
	for _, k := range keys {
		wait.Add(1)
		go func(_k int) {
			defer wait.Done()
			_, _ = localCache.Get(_k)
		}(k)
	}

	wait.Wait()

	end := time.Now().UnixMilli()

	t.Logf("consume: %dms", end-start)
}
